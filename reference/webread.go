package main


import (
	"fmt"
	//"net/http"
	//"html/template"
	//"encoding/xml"
	//"io/ioutil"
	//"sync"
    "log"
    "github.com/PuerkitoBio/goquery"
    //"go/doc"
    //"time"
    //"go/doc"
    "os"
    "bufio"
    "net/http"
    "io/ioutil"
    "strings"
    "sync"
    "sync/atomic"
)

type WebResult struct {
    links []string
    num int
}


func webWorker(id string, jobs <-chan string, results chan<- WebResult) {

    //just keeps waiting for job here until channel is closed
    for j := range jobs {

        url := "http://kamadenu.blogspot.ca/"+j

        fmt.Println("worker", id, "started  job : ", url)
        webResult := ScrapeImageLinks(url)
        fmt.Println("worker", id, "finished job : ", url)


        results <- *webResult
    }
}

func ScrapeImageLinks(url string) *WebResult {
    doc, err := goquery.NewDocument(url)
    if err != nil {
        log.Fatal(err)
    }

    MAX_IMG_PER_PAGE := 40
    links := make([]string, MAX_IMG_PER_PAGE)

    //TODO replace with xpath selector
    idx := 0
    doc.Find(".separator  a").Each(func(i int, s *goquery.Selection) {
      s.Find("img").Each(func(i int, s1 *goquery.Selection) {
          link, ok := s1.Attr("src")
          if ok{
              links[idx] = link
              idx += 1
          }
      })
    })

    webResult := new(WebResult)
    webResult.links = links
    webResult.num = idx

    return webResult
}

//https://gobyexample.com/worker-pools
func ScraperWorkers() {

    webjobs := make(chan string, 100)
    webresults := make(chan WebResult, 100) //TODO make *WebResult

    //create web workers to consume jobs
    for w :=1; w<=3; w++{
        go webWorker(fmt.Sprintf("%d/webworker",w), webjobs, webresults)
    }
    //create web jobs
    for year:=2013; year<2014; year++ {
        for month := 1; month < 13; month++ {
            webjobs <- fmt.Sprintf("%d/%02d/", year,month)
        }
    }
    close(webjobs) //close so range wont block (awaiting next sender)


    f, err := os.Create("scrape.txt")
    check(err)
    defer f.Close()
    w := bufio.NewWriter(f)

    //await web workers to finish first N jobs
    for a:=1; a <= 12; a++{
        result := <-webresults
        fmt.Println("worker result: ", result.num)
        for i:=0; i<result.num;i++{
            //fmt.Println(result.links[i])
            w.WriteString(result.links[i])
            w.WriteString("\n")
        }
        //fmt.Println("Got webresult :\n", result)
    }
    w.Flush()

}


func check(e error) {
    if e != nil {
        panic(e)
    }
}


var wg sync.WaitGroup
var saves uint64 = 0
var errors uint64 = 0

func ImageFetchWorker(id string, jobs <-chan string){
    defer wg.Done()
    fmt.Printf("worker %s started\n", id)


    for image_url := range jobs {

        n2 := strings.TrimSpace(image_url)
        resp, err := http.Get(n2)
        if resp.StatusCode == 404 || err != nil{
            //panic(404)
            atomic.AddUint64(&errors, 1)
            continue
        }
        //
        bytes, err := ioutil.ReadAll(resp.Body)
        if err != nil{
            atomic.AddUint64(&errors, 1)
            continue
            //panic(err)
        }


        stringSlice := strings.Split(image_url, "/")
        name := stringSlice[len(stringSlice) - 1]
        //name = name[:len(name)-4] //.jpg removal

        img, _ := os.Create(fmt.Sprintf("images/img%s", name))
        img.Write(bytes)
        img.Close()


        //https://gobyexample.com/mutexes
        //https://gobyexample.com/atomic-counters
        atomic.AddUint64(&saves, 1)
        //current := atomic.LoadInt64(&saves)
        if saves % 20 == 0{
            //fmt.Println("worker %s @ %d saves", id, current)
            fmt.Printf("worker %s @ %d saves", id, saves)
        }

    }
}


//https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/
func FetcherWorkers(){
    f, err := os.Open("scrape.txt")
    check(err)
    defer f.Close()


    //channels can become full, give 100 size for now
    jobs := make(chan(string), 100)

    //queue workers
    for w:=1; w<=3; w++{
        wg.Add(1)
        go ImageFetchWorker(fmt.Sprintf("worker%d",w), jobs)
    }
    //close(jobs) //close so range wont block (awaiting next sender)

    //queue jobs
    scanner := bufio.NewScanner(f)
    i := 0
    for scanner.Scan() {

        if i % 100 == 0{fmt.Println("queued %d jobs", i)}

        jobs<- scanner.Text()
        i++
        if i > 100{
            break
        }
    }

    fmt.Println("waiting")
    close(jobs)
    wg.Wait()


    fmt.Printf("done with %d image saves\n", saves)
}

func main() {
    usage := "Usage: \n\n go run webread.go scrape \n go run webread.go fetch"
    if len(os.Args) < 2 {
        fmt.Println(usage)
        os.Exit(0)
    }

    arg := os.Args[1]
    if (arg == "scrape") {
        fmt.Println("scraping")
        ScraperWorkers()
    }else if (arg == "fetch") {
        fmt.Println("fetching from scraped locations in scrape.txt file")
        FetcherWorkers()
    }else{
        fmt.Println(usage)
    }


}
