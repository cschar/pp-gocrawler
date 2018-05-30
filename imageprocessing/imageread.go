//https://www.devdungeon.com/content/working-images-go
// This example demonstrates decoding a JPEG image and examining its pixels.
package imageprocessing

import (
	//"encoding/base64"
	"fmt"
	"image"
	"log"
    "math/rand"
	//"strings"
	 _ "image/gif"
	 _ "image/png"
	_ "image/jpeg"
    "os"
    //"math"
    "image/png"
    //"math"
    //"image/color"
    "image/draw"
    //"bytes"
    "strconv"
    //"github.com/boltdb/bolt"
    _ "github.com/mattn/go-sqlite3"
    "github.com/nfnt/resize"

    "database/sql"
    "path/filepath"

    "math"
    "time"
)



func clonePix(b []uint8) []byte {
	c := make([]uint8, len(b))
	copy(c, b)
	return c
}

func CloneToRGBA(src image.Image) draw.Image {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}

func CloneRectToRGBA(src image.Image, rect image.Rectangle) draw.Image{
    dst := image.NewRGBA(rect)
    draw.Draw(dst, rect, src, rect.Min, draw.Src)
    return dst
}


func getImage(imageName string) image.Image{
    choice_ier, err := os.Open(imageName)
    if err != nil {
        log.Fatal(err)
    }
    defer choice_ier.Close()
    m, _, err := image.Decode(choice_ier)
    if err != nil {
        log.Fatal(err)
    }
    //fmt.Println("bounds of opened image", m.Bounds())
    return m
}


func initDB(){
    db, err := sql.Open("sqlite3", "./imagedata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
    //id integer not null primary key,
    sqlStmt := `
	create table rgbavg (
	  imageName text, dim text, x int, y int, r int, g int, b int);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		//log.Printf("%q: %s\n", err, sqlStmt)
	}
}

//const xsplits = 10
//const ysplits = 10
const xsplits = 20
const ysplits = 20
//const xsplits = 30
//const ysplits = 30
// const xsplits = 60
// const ysplits = 60

func InitData(){

    //have to do this so math.rand works...
    rand.Seed(time.Now().UnixNano())

    initDB()
    sliceAnalyzeSave("images/figs.jpg")
    sliceAnalyzeSave("images/oranges.jpg")
    sliceAnalyzeSave("images/apples.jpg")
    sliceAnalyzeSave("images/kiwi.jpg")
    sliceAnalyzeSave("images/blueberries.jpg")
    sliceAnalyzeSave("images/chameleon_blueyellowgreen.jpg")

    //MakeImageFromSlices("public/input/eyemazestyle.jpg")
    //MakeImageFromSlices("input/snowmandala.jpg")


}

func MakeImageFromSlices(imageName string) string{
    // MakeImageFromSlicesCustomThreshold(imageName, 30.0, 70.0, 40.0)
    return MakeImageFromSlicesCustomThreshold(imageName, 20.0, 30.0, 20.0)
}

// func Shuffle(slice interface{}) {
//     rv := reflect.ValueOf(slice)
//     swap := reflect.Swapper(slice)
//     length := rv.Len()
//     for i := length - 1; i > 0; i-- {
//             j := rand.Intn(i + 1)
//             swap(i, j)
//     }
// }

func MakeImageFromSlicesCustomThreshold(imageName string,
     rThreshold float64, gThreshold float64, bThreshold float64) string{
    m := getImage(imageName)
    /////////////
    /////////////
    //slice and analyze
    m = resize.Resize(960, 540, m, resize.Lanczos3)
    bounds := m.Bounds()
    yregionsize := bounds.Max.Y / ysplits
    xregionsize := bounds.Max.X / xsplits

    //make a 3D slice
    var inputPixelSums = make([][][]uint32, xsplits)
    for i:=0;i<xsplits;i++{
        inputPixelSums[i] = make([][]uint32, ysplits)
        for j:=0;j<ysplits;j++{
            inputPixelSums[i][j] = make([]uint32, 3)
        }
    }

    //var inputPixelSums [ysplits][xsplits][3]uint32
    inputRGBBucketAvg := make([][][]float32, xsplits)
    for i:=0;i<xsplits;i++{
        inputRGBBucketAvg[i] = make([][]float32, ysplits)
        for j:=0;j<ysplits;j++{
            inputRGBBucketAvg[i][j] = make([]float32, 3)
        }
    }

    fmt.Printf("Region Size: %dx%d split into: %dx%d  tiles\n",
        xregionsize, yregionsize, xsplits, ysplits)

    //Sum up each pixel's rgb into its xy buckets
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
            r = r>>8
            g = g>>8
            b = b>>8
            ybucket := y / yregionsize
            xbucket := x / xregionsize
            inputPixelSums[xbucket][ybucket][0] += r
            inputPixelSums[xbucket][ybucket][1] += g
            inputPixelSums[xbucket][ybucket][2] += b
		}
    }

    //Compute rgb average of each x,y bucket
    for i:=0; i< xsplits; i++{
        for j:=0; j<ysplits; j++ {
            for k := 0; k < 3; k++ {
                // divide by how many pixels reside in one bucket
                d := float32(xregionsize * yregionsize)
                inputRGBBucketAvg[i][j][k] = float32(inputPixelSums[i][j][k]) / d
            }
        }
    }

    /////////////
    /////////////
    // scan Database for matching tile choices

    MAX_CHOICES := 10000
    var tileNameBank = make([]string, MAX_CHOICES)
    var rgbtile_choices = make([][]byte, MAX_CHOICES)
    for i:=0;i<MAX_CHOICES;i++{
        rgbtile_choices[i] = make([]byte, 3)
    }

    db, err := sql.Open("sqlite3", "./imagedata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    //all rows for now
    //rows, err := db.Query("select imageName, dim, r, g, b from rgbavg ")
    rows, err := db.Query("select imageName, dim, r,g,b from rgbavg ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
    choices_read := 0
	for rows.Next() {
        var imageName string
        var dim string
        var red int
        var green int
        var blue int
        err = rows.Scan(&imageName, &dim, &red, &green, &blue)
        if err != nil {
            log.Fatal(err)
        }
        //fmt.Println(imageName, dim, red)

        tileNameBank[choices_read] = imageName
        rgbtile_choices[choices_read][0] = byte(red)
        rgbtile_choices[choices_read][1] = byte(green)
        rgbtile_choices[choices_read][2] = byte(blue)
        //fmt.Printf("%d choices_read \n",rval[choices_read])
        choices_read += 1
        if choices_read >= MAX_CHOICES {
            break
        }
    }


    //for each bucket, look for a match
    outputImg := image.NewRGBA(bounds)
    for y :=0; y < ysplits; y++{
        for x := 0; x <xsplits; x++{

            redAvg := inputRGBBucketAvg[x][y][0]
            greenAvg := inputRGBBucketAvg[x][y][1]
            blueAvg := inputRGBBucketAvg[x][y][2]


            var rCHOICES = 20 // max chioces for winning tile
            var hitCounter = 0;
            var hitNames = make([]string, rCHOICES)
            for i:=0; i < len(rgbtile_choices); i++{
                if  math.Abs(float64(float32(rgbtile_choices[i][0]) - redAvg)) < rThreshold &&
                    math.Abs(float64(float32(rgbtile_choices[i][1]) - greenAvg)) < gThreshold &&
                    math.Abs(float64(float32(rgbtile_choices[i][2]) - blueAvg)) < bThreshold {

                //fmt.Printf("Found hit for %dx%d spot with value %d near avg %f\n",
                //    x,y,rval[i], regionAvg)

                    hitNames[hitCounter] = tileNameBank[i]
                    hitCounter += 1;
                    if hitCounter >= rCHOICES {
                        break
                    }
                }
            }
            if hitCounter > 0 {
                hitChoice := rand.Intn(hitCounter)
                winnerImageName := hitNames[hitChoice] //0->4
                
                //open picture and draw it to current image
                srcImg := getImage(winnerImageName)

                rect := srcImg.Bounds()
                rectPoint := image.Pt(x*xregionsize,y*yregionsize)
                rectInOutputImg := rect.Add(rectPoint)

                draw.Draw(outputImg, rectInOutputImg, srcImg, image.Pt(0,0), draw.Src)

                fmt.Printf("Hit choice %d out of %d total choices \n",
                    hitChoice, hitCounter)
                fmt.Printf("Drawing @ %dx%d with winner %s\n", x,y,winnerImageName)

            }else{
              
                //no hits... use random tile
                winnerImageName := tileNameBank[rand.Intn(choices_read)]
                srcImg := getImage(winnerImageName)

                rect := srcImg.Bounds()
                rectPoint := image.Pt(x*xregionsize,y*yregionsize)
                rectInOutputImg := rect.Add(rectPoint)

                draw.Draw(outputImg, rectInOutputImg, srcImg, image.Pt(0,0), draw.Src)  
            }


            hitCounter = 0;
        }

    }

    //save as png instead of input jpg file format
    _, noPathName := filepath.Split(imageName)
    var extension = filepath.Ext(noPathName)
    var noExtensionName = noPathName[0:len(noPathName)-len(extension)]

    outputFile, err := os.Create(fmt.Sprintf("public/output/%s.png",noExtensionName))
            if err != nil {
                fmt.Printf("cant save file")
                log.Fatal("cant save file")
            }
            png.Encode(outputFile, outputImg)
            outputFile.Close()
    
    return fmt.Sprintf("%s.png",noExtensionName)
}

func sliceAnalyzeSave(imageName string){

    db, err := sql.Open("sqlite3", "./imagedata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


    m := getImage(imageName)
    m = resize.Resize(960, 540, m, resize.Lanczos3)
    bounds := m.Bounds()


    yregionsize := bounds.Max.Y / ysplits
    xregionsize := bounds.Max.X / xsplits

    //make a 3D slice
    var inputPixelSums = make([][][]uint32, xsplits)
    for i:=0;i<xsplits;i++{
        inputPixelSums[i] = make([][]uint32, ysplits)
        for j:=0;j<ysplits;j++{
            inputPixelSums[i][j] = make([]uint32, 3)
        }
    }

    //var inputPixelSums [ysplits][xsplits][3]uint32
    inputRGBBucketAvg := make([][][]float32, xsplits)
    for i:=0;i<xsplits;i++{
        inputRGBBucketAvg[i] = make([][]float32, ysplits)
        for j:=0;j<ysplits;j++{
            inputRGBBucketAvg[i][j] = make([]float32, 3)
        }
    }

    fmt.Printf("Region Splits: %dx%d %dx%d\n",
        xregionsize, yregionsize, xsplits, ysplits)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
            r = r>>8
            g = g>>8
            b = b>>8
            ybucket := y / yregionsize
            xbucket := x / xregionsize
            inputPixelSums[ybucket][xbucket][0] += r
            inputPixelSums[ybucket][xbucket][1] += g
            inputPixelSums[ybucket][xbucket][2] += b
		}
    }

    //Compute rgb average of each x,y bucket
    for i:=0; i< xsplits; i++{
        for j:=0; j<ysplits; j++ {
            for k := 0; k < 3; k++ {
                //d := float32(255.0 * yregionsize * xregionsize)
                d := float32(xregionsize * yregionsize)
                inputRGBBucketAvg[i][j][k] = float32(inputPixelSums[i][j][k]) / d
            }
        }
    }

    //save rgb avg for each x,y bucket
    fmt.Println("saving avgs to db")
    tx, err := db.Begin();
    if err != nil {
        log.Fatal(err)
    }
    //
    stmt, err := tx.Prepare("insert into" +
        " rgbavg(imageName, dim, x,y,r,g,b)" +
        " values(?,          ?,   ?,?,?,?,?)")
    if err != nil {log.Fatal(err)}
    defer stmt.Close()



    for y :=0; y < ysplits; y++{
        for x := 0; x <xsplits; x++{

            b := image.Rect(0,0, xregionsize, yregionsize)
            rectPoint := image.Pt(x*xregionsize,y*yregionsize)
            b = b.Add(rectPoint)


            myImage := CloneRectToRGBA(m, b)

            dimensions := strconv.Itoa(xsplits) + "x" + strconv.Itoa(ysplits)
            region := strconv.Itoa(x) + "-"+ strconv.Itoa(y)
            var extension = filepath.Ext(imageName)
            var name = imageName[0:len(imageName)-len(extension)]

            fullPathFileName := name+"/"+dimensions + "-" + region +".png"

            _, endname := filepath.Split(fullPathFileName)
            fmt.Println(endname)

            _, err = stmt.Exec(fullPathFileName,
                    fmt.Sprintf("%dx%d-%s", xsplits,ysplits,region),
                    x,y,
                    int(inputRGBBucketAvg[x][y][0]),
                    int(inputRGBBucketAvg[x][y][1]),
                    int(inputRGBBucketAvg[x][y][2]))

            if err != nil {
                log.Fatal(err)
            }

            os.MkdirAll(name, os.ModePerm);
            fmt.Println(fullPathFileName)


            outputFile, err := os.Create(fullPathFileName)
            if err != nil {
                log.Fatal("cant save file")
            }
            png.Encode(outputFile, myImage)
            outputFile.Close()
        }

    }
    tx.Commit()

}







