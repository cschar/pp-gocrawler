package main

import (
	"net/http"
	"ppgocrawler/controllers"
	"log"
	"time"
	"math/rand"
)



func main() {
	rand.Seed(time.Now().UnixNano())

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/agg/", controllers.NewsAggHandler)

    http.HandleFunc("/upload", controllers.UploadFile)
	http.HandleFunc("/mixed", controllers.MixedImages)
	http.HandleFunc("/input", controllers.InputImages)

	log.Println("Running")
	http.ListenAndServe(":8000", nil)
}