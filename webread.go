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
func QueueWorkers() {

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


    //await web workers to finish first N jobs
    for a:=1; a <= 12; a++{
        result := <-webresults
        fmt.Println("worker result: ")
        for i:=0; i<result.num;i++{
            fmt.Println(result.links[i])
        }
        //fmt.Println("Got webresult :\n", result)
    }

}





func ExampleScrape() {
    doc, err := goquery.NewDocument("http://kamadenu.blogspot.ca/2016/02/")
    //doc, err := goquery.NewDocument("http://kamadenu.blogspot.ca")
  //doc, err := goquery.NewDocument("http://metalsucks.net")
  if err != nil {
    log.Fatal(err)
  }


    MAX_IMG_PER_PAGE := 40
    links := make([]string, MAX_IMG_PER_PAGE)
  // Find the review items
    idx := 0
    doc.Find(".separator  a").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    //band := s.Find("a").Text()
    //title := s.Find("i").Text()
    //fmt.Printf("Review %d: %s - %s\n", i, band, title)
      fmt.Println(s.Html())

      fmt.Println()
      s.Find("img").Each(func(i int, s1 *goquery.Selection) {
          link, ok := s1.Attr("src")
          fmt.Println(link, ok)
          if ok{
              links[idx] = link
              idx += 1
          }
      })

  })

    fmt.Println("collected links:")
    for i:=0; i<idx; i++{
        val := links[i]
        fmt.Println(i, val)
    }
}


func main() {
    QueueWorkers()
    //ExampleScrape()

}
