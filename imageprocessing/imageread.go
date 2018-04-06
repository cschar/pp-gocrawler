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
    "image/color"
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
    reader, err := os.Open(imageName)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()
    m, _, err := image.Decode(reader)
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


func main(){

    //have to do this so math.rand works...
    rand.Seed(time.Now().UnixNano())

    //initDB()
    //sliceAnalyzeSave("images/figs.jpg")
    //sliceAnalyzeSave("images/oranges.jpg")
    //sliceAnalyzeSave("images/apples.jpg")


    MakeImageFromSlices("public/input/eyemazestyle.jpg")
    //MakeImageFromSlices("input/snowmandala.jpg")


}

func MakeImageFromSlices(imageName string){
    m := getImage(imageName)
    /////////////
    /////////////
    //slice and analyze
    m = resize.Resize(960, 540, m, resize.Lanczos3)
    bounds := m.Bounds()
    yregionsize := bounds.Max.Y / ysplits
    xregionsize := bounds.Max.X / xsplits

    //make a 3D slice
    var rgbsumregions = make([][][]uint32, xsplits)
    for i:=0;i<xsplits;i++{
        rgbsumregions[i] = make([][]uint32, ysplits)
        for j:=0;j<ysplits;j++{
            rgbsumregions[i][j] = make([]uint32, 3)
        }
    }

    //var rgbsumregions [ysplits][xsplits][3]uint32
    avgregions := make([][][]float32, xsplits)
    for i:=0;i<xsplits;i++{
        avgregions[i] = make([][]float32, ysplits)
        for j:=0;j<ysplits;j++{
            avgregions[i][j] = make([]float32, 3)
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
            yregionbucket := y / yregionsize
            xregionbucket := x / xregionsize
            rgbsumregions[xregionbucket][yregionbucket][0] += r
            rgbsumregions[xregionbucket][yregionbucket][1] += g
            rgbsumregions[xregionbucket][yregionbucket][2] += b
		}
    }

    //Compute rgb average of each x,y bucket
    for i:=0; i< xsplits; i++{
        for j:=0; j<ysplits; j++ {
            for k := 0; k < 3; k++ {
                //d := float32(255.0 * yregionsize * xregionsize)
                d := float32(xregionsize * yregionsize)
                avgregions[i][j][k] = float32(rgbsumregions[i][j][k]) / d
            }
        }
    }

    /////////////
    /////////////
    //scan other images for R axis similarity

    CHOICES := 1000
    var rname = make([]string, CHOICES) //1000 images to choose from
    var rgbval = make([][]byte, CHOICES)
    for i:=0;i<CHOICES;i++{
        rgbval[i] = make([]byte, 3)
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
    read := 0
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

        rname[read] = imageName
        rgbval[read][0] = byte(red)
        rgbval[read][1] = byte(green)
        rgbval[read][2] = byte(blue)
        //fmt.Printf("%d read \n",rval[read])
        read += 1
        if read >= CHOICES {
            break
        }
    }


    //for each region, look for a match
    outputImg := image.NewRGBA(bounds)
    for y :=0; y < ysplits; y++{
        for x := 0; x <xsplits; x++{

            redAvg := avgregions[x][y][0]
            greenAvg := avgregions[x][y][0]
            blueAvg := avgregions[x][y][2]


            var rCHOICES = 10 // max chioces for winning tile
            var hitCounter = 0;
            var hitNames = make([]string, rCHOICES)
            for i:=0; i < len(rgbval); i++{
                if  math.Abs(float64(float32(rgbval[i][0]) - redAvg)) < 10.0 &&
                    math.Abs(float64(float32(rgbval[i][1]) - greenAvg)) < 70.0 &&
                    math.Abs(float64(float32(rgbval[i][2]) - blueAvg)) < 40.0 {

                //fmt.Printf("Found hit for %dx%d spot with value %d near avg %f\n",
                //    x,y,rval[i], regionAvg)

                    hitNames[hitCounter] = rname[i]
                    hitCounter += 1;
                    if hitCounter >= rCHOICES {
                        break
                    }
                }
            }
            if hitCounter > 0 {
                winnerImageName := hitNames[rand.Intn(hitCounter)] //0->4
                //open picture and draw it to current image
                srcImg := getImage(winnerImageName)

                rect := srcImg.Bounds()
                rectPoint := image.Pt(x*xregionsize,y*yregionsize)
                rectInOutputImg := rect.Add(rectPoint)

                draw.Draw(outputImg, rectInOutputImg, srcImg, image.Pt(0,0), draw.Src)

                //fmt.Printf("Drawing @ %dx%d with winner %s\n", x,y,winnerImageName)
            }

            hitCounter = 0;
        }

    }

    _, noPathName := filepath.Split(imageName)
    var extension = filepath.Ext(noPathName)
    var noExtensionName = noPathName[0:len(noPathName)-len(extension)]

    outputFile, err := os.Create(fmt.Sprintf("public/output/%s.png",noExtensionName))
            if err != nil {
                log.Fatal("cant save file")
            }
            png.Encode(outputFile, outputImg)
            outputFile.Close()
    //tx.Commit()

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
    var rgbsumregions = make([][][]uint32, xsplits)
    for i:=0;i<xsplits;i++{
        rgbsumregions[i] = make([][]uint32, ysplits)
        for j:=0;j<ysplits;j++{
            rgbsumregions[i][j] = make([]uint32, 3)
        }
    }

    //var rgbsumregions [ysplits][xsplits][3]uint32
    avgregions := make([][][]float32, xsplits)
    for i:=0;i<xsplits;i++{
        avgregions[i] = make([][]float32, ysplits)
        for j:=0;j<ysplits;j++{
            avgregions[i][j] = make([]float32, 3)
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
            yregionbucket := y / yregionsize
            xregionbucket := x / xregionsize
            rgbsumregions[yregionbucket][xregionbucket][0] += r
            rgbsumregions[yregionbucket][xregionbucket][1] += g
            rgbsumregions[yregionbucket][xregionbucket][2] += b
		}
    }

    //Compute rgb average of each x,y bucket
    for i:=0; i< xsplits; i++{
        for j:=0; j<ysplits; j++ {
            for k := 0; k < 3; k++ {
                //d := float32(255.0 * yregionsize * xregionsize)
                d := float32(xregionsize * yregionsize)
                avgregions[i][j][k] = float32(rgbsumregions[i][j][k]) / d
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
                    int(avgregions[x][y][0]),
                    int(avgregions[x][y][1]),
                    int(avgregions[x][y][2]))

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


func sliceSave(){
    m := getImage("images/snowmandala.jpg")
    bounds := m.Bounds()

    yregionsize := bounds.Max.Y / ysplits
    xregionsize := bounds.Max.X / xsplits

    //make a 3D slice
    var rgbsumregions = make([][][]uint32, xsplits)
    for i:=0;i<xsplits;i++{
        rgbsumregions[i] = make([][]uint32, ysplits)
        for j:=0;j<ysplits;j++{
            rgbsumregions[i][j] = make([]uint32, 3)
        }
    }

    //var rgbsumregions [ysplits][xsplits][3]uint32
    avgregions := make([][][]float32, xsplits)
    for i:=0;i<xsplits;i++{
        avgregions[i] = make([][]float32, ysplits)
        for j:=0;j<ysplits;j++{
            avgregions[i][j] = make([]float32, 3)
        }
    }

    fmt.Printf("Region Splits: %dx%d %dx%d", xregionsize, yregionsize, xsplits, ysplits)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
            r = r>>8
            g = g>>8
            b = b>>8
            yregionbucket := y / yregionsize
            xregionbucket := x / xregionsize
            rgbsumregions[yregionbucket][xregionbucket][0] += r
            rgbsumregions[yregionbucket][xregionbucket][1] += g
            rgbsumregions[yregionbucket][xregionbucket][2] += b
		}
    }

    //Compute rgb average of each x,y bucket
    for i:=0; i< ysplits; i++{
        for j:=0; j<xsplits; j++ {
            for k := 0; k < 3; k++ {
                //d := float32(255.0 * yregionsize * xregionsize)
                d := float32(yregionsize * xregionsize)
                avgregions[i][j][k] = float32(rgbsumregions[i][j][k]) / d
            }
        }
    }


    fmt.Println(bounds)
    for y :=0; y < ysplits; y++{
        for x := 0; x <xsplits; x++{
            //myImage := image.NewRGBA(image.Rect(0, 0, xregionsize, yregionsize))

            b := image.Rect(0,0,
                xregionsize, yregionsize)


            rectPoint := image.Pt(x*xregionsize,y*yregionsize)
            b = b.Add(rectPoint)
            fmt.Println(b)

            //dst := image.NewRGBA(b)
            //draw.Draw(myImage, b, m, b.Min, draw.Src)
            //return dst
            dimensions := strconv.Itoa(xregionsize) + "x" + strconv.Itoa(yregionsize)
            region := strconv.Itoa(x) + "-"+ strconv.Itoa(y)
            fileName := "images/"+ dimensions + "-" + region + ".png"

            fmt.Println(fileName)

            myImage := CloneRectToRGBA(m, b)
            outputFile, err := os.Create(fileName)
            if err != nil {
                // Handle error
            }
            png.Encode(outputFile, myImage)
            // Don't forget to close files
            outputFile.Close()
        }

    }
}

func LowerResolution() {
	// Decode the JPEG data. If reading from file, create a reader with
	//
    //reader, err := os.Open("eyemazestyle.jpg")
    m := getImage("images/snowmandala.jpg")
    bounds := m.Bounds()
    myImage := CloneToRGBA(m)


    var rgbsumregions [ysplits][xsplits][3]uint32

    yregionsize := bounds.Max.Y / ysplits
    xregionsize := bounds.Max.X / xsplits
    fmt.Printf("%d %d is xy region sizes", xregionsize, yregionsize)


    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

            r = r>>8
            g = g>>8
            b = b>>8

            yregionbucket := y / yregionsize
            xregionbucket := x / xregionsize

            rgbsumregions[yregionbucket][xregionbucket][0] += r
            rgbsumregions[yregionbucket][xregionbucket][1] += g
            rgbsumregions[yregionbucket][xregionbucket][2] += b


		}
    }

    //Compute rgb average of each x,y bucket
    var avgregions [xsplits][ysplits][3]float32
    for i:=0; i< ysplits; i++{
        for j:=0; j<xsplits; j++ {
            for k := 0; k < 3; k++ {
                //d := float32(255.0 * yregionsize * xregionsize)
                d := float32(yregionsize * xregionsize)
                avgregions[i][j][k] = float32(rgbsumregions[i][j][k]) / d
            }
        }
    }

    //Print results
    //fmt.Println("sumregion [0][0]")
    //var buffer bytes.Buffer;
    //for j :=0; j < ysplits; j++{
    //    for i := 0; i<xsplits; i++{
    //        s := strconv.FormatFloat(float64(avgregions[j][i][0]), 'f',-1,32)
    //        buffer.WriteString( s + ", ")
    //        //fmt.Printf("%f\n", avgregions[j][0][0])
    //    }
    //    fmt.Println(buffer.String())
    //    buffer.Reset()
    //}



    //Use average to simplify input image
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

            r = r>>8
            g = g>>8
            b = b>>8

            yregionbucket := y / yregionsize
            xregionbucket := x / xregionsize


            newR := uint8(avgregions[yregionbucket][xregionbucket][0])
            newG := uint8(avgregions[yregionbucket][xregionbucket][1])
            newB := uint8(avgregions[yregionbucket][xregionbucket][2])

            //newColor :=
            myImage.Set(x,y, color.RGBA{newR, newG, newB, 255})

		}
    }


    // outputFile is a File type which satisfies Writer interface
    outputFile, err := os.Create("testRegion.png")
    if err != nil {
    	// Handle error
    }
    png.Encode(outputFile, myImage)
    // Don't forget to close files
    outputFile.Close()
}





