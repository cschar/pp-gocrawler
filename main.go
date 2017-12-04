package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/xml"
	"io/ioutil"
)

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Whoa, Go is neat!</h1>")
}



func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	var s Sitemapindex
	var n News
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	news_map := make(map[string]NewsMap)

	dups := 0

	for _, Location := range s.Locations {
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes, &n)

		for idx, _ := range n.Keywords {
			if _, ok := news_map[n.Titles[idx]]; ok {
				dups += 1
			}
			news_map[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}
		fmt.Printf("SiteIndex Location %s\n", Location)
		fmt.Println("keywords", len(n.Keywords))
		fmt.Println("titles", len(n.Titles))
		fmt.Println("locations", len(n.Locations))
	}
	fmt.Println("After crawling, map had duplicates:", dups)
	fmt.Println("Map size is", len(news_map))

	p := NewsAggPage{Title: "Amazing News Aggregator", News: news_map}
	t, _ := template.ParseFiles("newsaggtemplate.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg/", newsAggHandler)
	http.ListenAndServe(":8000", nil)
}