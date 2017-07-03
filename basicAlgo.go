package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var emojiListAvg [2661][3]float64

func initEmojiDictAvg() {
	for i := 0; i < len(emojiList); i++ {
		r, g, b := averageRGBArray(emojiList[i].vectorForm, 0, 0, 72)
		emojiListAvg[i][0] = r
		emojiListAvg[i][1] = g
		emojiListAvg[i][2] = b
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
	r1, g1, b1 := emojiListAvg[0][0], emojiListAvg[0][1], emojiListAvg[0][2]
	//r1, g1, b1 := averageRGBArray(emojiList[0].vectorForm, 0, 0, 72)

	//use sum of square differences
	smallestDistance += math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
	nearestIndex = 0

	for i := 1; i < len(emojiList); i++ {
		r1, g1, b1 := emojiListAvg[i][0], emojiListAvg[i][1], emojiListAvg[i][2]
		//r1, g1, b1 := averageRGBArray(emojiList[i].vectorForm, 0, 0, 72)
		distance = math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
		if distance < smallestDistance {
			smallestDistance = distance
			nearestIndex = i
		}
	}

	return nearestIndex
}

func simpleAlgo(img image.Image) {
	//hardcode square size for now
	fmt.Println("enter square size")
	var squareSize int
	fmt.Scanf("%d\n", &squareSize)

	imgWidth := img.Bounds().Max.X - img.Bounds().Min.X
	imgHeight := img.Bounds().Max.Y - img.Bounds().Min.X

	resultWidth = 72 * imgWidth / squareSize
	resultHeight = 72 * imgHeight / squareSize
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
		fmt.Println(err)
	}

	if err := png.Encode(f, resultImg); err != nil {
		f.Close()
		fmt.Println(err)
	}

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("finished")

}
