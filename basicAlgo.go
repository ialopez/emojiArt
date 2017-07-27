package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

/*this dictionary holds the precomputed average RGB values of Emojis
it is a map that takes in the platform as a key ie.. "apple" "facebook" and returns a [][3] slice of floats
where [][0] is red [][1] is green and [][2] is blue
*/
var emojiDictAvg map[string][][3]float64
var squareSize int

func initEmojiDictAvg() {
	emojiDictAvg = make(map[string][][3]float64)
	for i := 0; i < len(platforms); i++ {
		emojiDictAvg[platforms[i]] = make([][3]float64, len(emojiDict[platforms[i]]))
		for j := 0; j < len(emojiDict[platforms[i]]); j++ {
			r, g, b := averageRGBArray(emojiDict[platforms[i]][j].vectorForm, 0, 0, 64)
			emojiDictAvg[platforms[i]][j][0] = r
			emojiDictAvg[platforms[i]][j][1] = g
			emojiDictAvg[platforms[i]][j][2] = b
		}
	}
	fmt.Println("average emoji dict initialized")
}

//test this
func nearestSimple(subsection [][]color.Color, squareSize int) int {

	var smallestDistance float64
	var nearestIndex int
	var distance float64

	//set smallest distance to first emoji
	r0, g0, b0 := averageRGBSlice(subsection, 0, 0, squareSize)
	r1, g1, b1 := emojiDictAvg[currentPlatform][0][0], emojiDictAvg[currentPlatform][0][1], emojiDictAvg[currentPlatform][0][2]
	//r1, g1, b1 := averageRGBArray(emojiDict[0].vectorForm, 0, 0, 64)

	//use sum of square differences
	smallestDistance = math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
	nearestIndex = 0

	for i := 1; i < len(emojiDictAvg[currentPlatform]); i++ {
		r1, g1, b1 := emojiDictAvg[currentPlatform][i][0], emojiDictAvg[currentPlatform][i][1], emojiDictAvg[currentPlatform][i][2]
		//r1, g1, b1 := averageRGBArray(emojiDict[i].vectorForm, 0, 0, 64)
		distance = math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
		if distance < smallestDistance {
			smallestDistance = distance
			nearestIndex = i
		}
	}

	return nearestIndex
}

func simpleAlgo(img image.Image) image.Image {
	//hardcode square size for now
	/*fmt.Println("enter square size")
	var squareSize int
	fmt.Scanf("%d\n", &squareSize)
	*/

	imgWidth := img.Bounds().Max.X - img.Bounds().Min.X
	imgHeight := img.Bounds().Max.Y - img.Bounds().Min.X

	resultWidth = emojiSize * imgWidth / squareSize
	resultHeight = emojiSize * imgHeight / squareSize
	var resultImg *image.RGBA = image.NewRGBA(image.Rect(0, 0, resultWidth, resultHeight))

	//create 2d slice
	subsection := make([][]color.Color, squareSize)
	for i := 0; i < squareSize; i++ {
		subsection[i] = make([]color.Color, squareSize)
	}

	var upperLeft image.Point
	currentSquare.X = 0
	currentSquare.Y = 0

	for upperLeft.X = img.Bounds().Min.X; upperLeft.X < img.Bounds().Max.X; upperLeft.X += squareSize {
		for upperLeft.Y = img.Bounds().Min.Y; upperLeft.Y < img.Bounds().Max.Y; upperLeft.Y += squareSize {
			for x := upperLeft.X; x < upperLeft.X+squareSize; x++ {
				for y := upperLeft.Y; y < upperLeft.Y+squareSize; y++ {
					subsection[x%squareSize][y%squareSize] = img.At(x, y)
				}
			}
			closestEmoji := nearestSimple(subsection, squareSize)
			drawEmoji(resultImg, closestEmoji)
		}
		fmt.Print(upperLeft.X)
		fmt.Print("out of ")
		fmt.Println(img.Bounds().Max.X)
	}

	fmt.Println("drawing image")
	f, err := os.Create("./image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, resultImg); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("finished")

	return resultImg

}
