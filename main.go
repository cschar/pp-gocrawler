package main

import (
	"net/http"
	"github.com/cschar/pp-gocrawler/controllers"
	"github.com/cschar/pp-gocrawler/imageprocessing"
	"log"
	"time"
	"math/rand"
	"fmt"
	"os"
)



func server() {
	rand.Seed(time.Now().UnixNano())

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/agg/", controllers.NewsAggHandler)

    http.HandleFunc("/upload", controllers.UploadFile)
	http.HandleFunc("/mixed", controllers.MixedImages)
	http.HandleFunc("/input", controllers.InputImages)

	log.Println("Running")
	http.ListenAndServe(":8123", nil)
}

func main(){
	usage := "Usage: \n\n" +
		" go run main.go server \n" +
		" go run main.go init"

	if len(os.Args) < 2 {
        fmt.Println(usage)
        os.Exit(0)
    }

    arg := os.Args[1]
    if (arg == "server") {
        server()
    }else if (arg == "init") {
		fmt.Println("Creating database and slicing/analyzing image rgb sections")
		imageprocessing.InitData()
    }else{
        fmt.Println(usage)
    }
}