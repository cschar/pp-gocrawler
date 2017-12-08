

# README

go tut code:
https://pythonprogramming.net/go/introduction-go-language-programming-tutorial/

go image code:
https://www.devdungeon.com/content/working-images-go
coloring:
use image/draw and the structs .Set method
https://stackoverflow.com/questions/28992396/draw-a-rectangle-in-golang

TODO:

#image stuff
- [DONE] split image into big regions of avg rgb
- [DONE] save regions, PRE-avg along with their rgb average data
- load in that data and fill in another image's regions with similar rgb
  - load in data, sliceAnalyze
  - scan other avgs in db for threshold similarity on R axis
  - get matching bits and replace original image
  - svae in output folder
  
  
- cleaning pipeline
   - resize -> split -> save

#web stuff
- worker pool consuming images from a crawler
  sort them into rgb categories (overall red, green or blue images)
https://gobyexample.com/worker-pools
