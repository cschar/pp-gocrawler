package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/xml"
	"io/ioutil"
	"sync"
	"io"
	"os"
	"strconv"
	"crypto/md5"
	"time"
	"ppgocrawler/controllers"
	"log"
)

var wg sync.WaitGroup

type NewsMap struct {
	Keyword string
	Location string
}

type NewsAggPage struct {
	Title string
	News map[string]NewsMap
}

type Sitemapindex struct {
	Locations []string `xml:"sitemap>loc"`
}

//http://www.washingtonpost.com/news-politics-sitemap.xml
//http://www.washingtonpost.com/news-technology-sitemap.xml
type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}


func newsRoutine(c chan News, Location string){
    defer wg.Done()

    var n News
    resp, _ := http.Get(Location)
    bytes, _ := ioutil.ReadAll(resp.Body)
    xml.Unmarshal(bytes, &n)
    resp.Body.Close()

    c <- n
}


func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	var s Sitemapindex

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	news_map := make(map[string]NewsMap)
    resp.Body.Close()
    queue := make(chan News, 30)
	dups := 0

	for _, Location := range s.Locations {
        wg.Add(1)
        go newsRoutine(queue, Location)
    }

    wg.Wait()
    close(queue)

    //elem is News type
    for elem := range queue {
		for idx, _ := range elem.Keywords {
			if _, ok := news_map[elem.Titles[idx]]; ok {
				dups += 1
			}
			news_map[elem.Titles[idx]] = NewsMap{elem.Keywords[idx], elem.Locations[idx]}
		}

		fmt.Println("keywords", len(elem.Keywords))
		fmt.Println("titles", len(elem.Titles))
		fmt.Println("locations", len(elem.Locations))
	}

	fmt.Println("After crawling, map had duplicates:", dups)
	fmt.Println("Map size is", len(news_map))

	p := NewsAggPage{Title: "Amazing News Aggregator", News: news_map}
	t, _ := template.ParseFiles("templates/newsaggtemplate.html")
	t.Execute(w, p)
}

func upload(w http.ResponseWriter, r *http.Request) {
       fmt.Println("method:", r.Method)
       if r.Method == "GET" {
           crutime := time.Now().Unix()
           h := md5.New()
           io.WriteString(h, strconv.FormatInt(crutime, 10))
           token := fmt.Sprintf("%x", h.Sum(nil))

           t, _ := template.ParseFiles("templates/upload.gtpl")
           t.Execute(w, token)
       } else {
           r.ParseMultipartForm(32 << 20)
           file, handler, err := r.FormFile("uploadfile")
           if err != nil {
               fmt.Println(err)
               return
           }
           defer file.Close()
           fmt.Fprintf(w, "%v", handler.Header)
           f, err := os.OpenFile("./uploadtest/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(err)
               return
           }
           defer f.Close()
           io.Copy(f, file)
       }
}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.gtpl")
	t.Execute(w, struct {
    Title string
    ImageCount int
	}{"ppgc", 0})

}





func main() {
	//http.HandleFunc("/", indexHandler)
	//http.HandleFunc("/agg/", newsAggHandler)
	//http.HandleFunc("/upload", upload)

	http.Handle("/", http.FileServer(http.Dir("./public")))
    http.HandleFunc("/upload", controllers.UploadFile)
	http.HandleFunc("/mixed", controllers.MixedImages)
	http.HandleFunc("/input", controllers.InputImages)
    log.Println("Running")

	http.ListenAndServe(":8000", nil)
}