

# README




## run server
```
go run main.go init   
go run main.go server
# in browser open localhost:8000
```


### run crawler with workers
```
go run webread.go scrape      # generates scrape.txt with image urls
go run webread.go fetch       # workers fetch the images
```

## building server for linux deployment



```
Make sure go is installed then:
go get -v github.com/cschar/pp-gocrawler   # this will also download the dependencies

arch   #x86_64 usually

#amd64 for x86_64
env GOARCH=amd64 GOOS=linux go build -o mainlin main.go

# or just ssh into linux box and build locally
go build -o mainlin main.go
```


# developing

when developing and pushing changes, to update the dependency on an external box:

```
go get -u all

# single package
go get -u github.com/cschar/pp-gocrawler

```

go tut code:
https://pythonprogramming.net/go/introduction-go-language-programming-tutorial/
https://dev.to/gcdcoder/how-to-upload-files-with-golang-and-ajax

go image code:
https://www.devdungeon.com/content/working-images-go
coloring:
use image/draw and the structs .Set method
https://stackoverflow.com/questions/28992396/draw-a-rectangle-in-golang

TODO:

#image stuff
- [DONE] split image into big regions of avg rgb
- [DONE] save regions, PRE-avg along with their rgb average data
- [DONE] load in that data and fill in another image's regions with similar rgb
        - load in data, sliceAnalyze
        - scan other avgs in db for threshold similarity on R axis
        - get matching bits and replace original image
        - svae in output folder
  
  
- cleaning pipeline
   - resize -> split -> save

#web stuff
- [DONE] worker pool fetching images from blog
- [DONE] worker pool consuming images from a crawler
- sort them into rgb categories (overall red, green or blue images)
https://gobyexample.com/worker-pools
