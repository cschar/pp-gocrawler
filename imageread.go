//https://www.devdungeon.com/content/working-images-go
// This example demonstrates decoding a JPEG image and examining its pixels.
package main

import (
	//"encoding/base64"
	"fmt"
	"image"
	"log"
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

func CloneRectToRGBA(src image.Image, rect image.Rectangle, rectPoint image.Point) draw.Image{
    dst := image.NewRGBA(rect)
    draw.Draw(dst, rect, src, rect.Min, draw.Src)
    return dst
}


func getImage(imageName string) image.Image{
    reader, err := os.Open("snowmandala.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()
    m, _, err := image.Decode(reader)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("bounds of opened image", m.Bounds())
    return m
}
//
//
//func getAVGRGBregions(m image.Image, xs, ys int) [][][]uint32{
//
//    const xsplits = xs
//    const ysplits = ys
//
//    bounds := m.Bounds()
//    var rgbsumregions [ysplits][xsplits][3]uint32
//    var avgregions [xsplits][ysplits][3]float32
//
//    yregionsize := bounds.Max.Y / ysplits
//    xregionsize := bounds.Max.X / xsplits
//    fmt.Printf("%d %d is xy region sizes", xregionsize, yregionsize)
//
//    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
//		for x := bounds.Min.X; x < bounds.Max.X; x++ {
//			r, g, b, _ := m.At(x, y).RGBA()
//            r = r>>8
//            g = g>>8
//            b = b>>8
//            yregionbucket := y / yregionsize
//            xregionbucket := x / xregionsize
//            rgbsumregions[yregionbucket][xregionbucket][0] += r
//            rgbsumregions[yregionbucket][xregionbucket][1] += g
//            rgbsumregions[yregionbucket][xregionbucket][2] += b
//		}
//    }
//
//    //Compute rgb average of each x,y bucket
//    for i:=0; i< ysplits; i++{
//        for j:=0; j<xsplits; j++ {
//            for k := 0; k < 3; k++ {
//                //d := float32(255.0 * yregionsize * xregionsize)
//                d := float32(yregionsize * xregionsize)
//                avgregions[i][j][k] = float32(rgbsumregions[i][j][k]) / d
//            }
//        }
//    }
//    return avgregions
//}

func main() {
	// Decode the JPEG data. If reading from file, create a reader with
	//
    //reader, err := os.Open("eyemazestyle.jpg")
    m := getImage("snowmandala.jpg")
    bounds := m.Bounds()
    //myImage := CloneToRGBA(m)


    const xsplits = 20
    const ysplits = 20

    //const xsplits = 4
    //const ysplits = 4
    //const xsplits = 10
    //const ysplits = 10

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
            fileName := "images/"+ strconv.Itoa(x) + "-"+ strconv.Itoa(y) + ".png"

            fmt.Println(fileName)

            myImage := CloneRectToRGBA(m, b, rectPoint)
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
    m := getImage("snowmandala.jpg")
    bounds := m.Bounds()
    myImage := CloneToRGBA(m)


    const xsplits = 20
    const ysplits = 20
    //const xsplits = 4
    //const ysplits = 4
    //const xsplits = 10
    //const ysplits = 10
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


func getImageStats(){
    //reader, err := os.Open("eyemazestyle.jpg")
    reader, err := os.Open("snowmandala.jpg")
    if err != nil {
        log.Fatal(err)
    }

    defer reader.Close()
    m, _, err := image.Decode(reader)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("bounds of opened image", m.Bounds())
    bounds := m.Bounds()
    var rgbsum [3]uint32
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
            rgbsum[0] += r>>8
            rgbsum[1] += g>>8
            rgbsum[2] += b>>8
		}
    }

    fmt.Println("RGBSUM: ", rgbsum)
    fmt.Println(rgbsum[0])
    fmt.Println(rgbsum[1])
    fmt.Println(rgbsum[2])

    fmt.Println("% 0 -> 100% avg of rgb value in total image")
    d := float32(255 * bounds.Max.Y * bounds.Max.X)
    fmt.Println(float32(rgbsum[0]) / d)
    fmt.Println(float32(rgbsum[1]) / d)
    fmt.Println(float32(rgbsum[2]) / d)

}

func createImage(){
    // Create a blank image 100x200 pixels
    myImage := image.NewRGBA(image.Rect(0, 0, 100, 200))

    // You can access the pixels through myImage.Pix[i]
    // One pixel takes up four bytes/uint8. One for each: RGBA
    // So the first pixel is controlled by the first 4 elements
    // Values for color are 0 black - 255 full color
    // Alpha value is 0 transparent - 255 opaque
    myImage.Pix[0] = 255 // 1st pixel red
    myImage.Pix[1] = 0 // 1st pixel green
    myImage.Pix[2] = 0 // 1st pixel blue
    myImage.Pix[3] = 255 // 1st pixel alpha

    i := 0

    // 100x200 = 20 000 pixels x 4bytes == 80 000
    for i < 20000{

        if i % 5 == 0 {
            myImage.Pix[i] = 110
            myImage.Pix[i+1] = 255
            myImage.Pix[i+2] = 30
            myImage.Pix[i+3] = 255
        } else {
            myImage.Pix[i] = 170
            myImage.Pix[i+1] = 255
            myImage.Pix[i+2] = 30
            myImage.Pix[i+3] = 255
        }


        i += 4;

    }

    fmt.Println(myImage.Stride) // 40 for an image 10 pixels wide



    // outputFile is a File type which satisfies Writer interface
    outputFile, err := os.Create("test.png")
    if err != nil {
    	// Handle error
    }

    // Encode takes a writer interface and an image interface
    // We pass it the File and the RGBA
    png.Encode(outputFile, myImage)

    // Don't forget to close files
    outputFile.Close()
}


func drawRGBOnCopy() {
	// Decode the JPEG data. If reading from file, create a reader with
	//
    //reader, err := os.Open("eyemazestyle.jpg")
    reader, err := os.Open("snowmandala.jpg")
    if err != nil {
        log.Fatal(err)
    }

    defer reader.Close()
    m, _, err := image.Decode(reader)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("bounds of opened image", m.Bounds())
    bounds := m.Bounds()
    myImage := CloneToRGBA(m)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

            if r>>8 < 100 {
                myImage.Set(x,y, color.RGBA{255, 0, 0, 255})
            }else if g>>8 < 100 {
                myImage.Set(x,y, color.RGBA{0, 255, 0, 255})
            }else if b>>8 < 100 {
                myImage.Set(x,y, color.RGBA{0, 0, 255, 255})
            }

		}
    }


    // outputFile is a File type which satisfies Writer interface
    outputFile, err := os.Create("testdrawrgb.png")
    if err != nil {
    	// Handle error
    }

    png.Encode(outputFile, myImage)

    // Don't forget to close files
    outputFile.Close()
}