package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var emojiList [2661]Emoji //for now we hard code this value
var numOfEmojisForResult = 30
var inputRatio = 0
var outputRatio = 0
var totalDistance = uint32(0)
var currentSquare image.Point
var resultHeight int
var resultWidth int

//emoji png files are saved internally in the program as 72x72 arrays of colors, also the path of the emoji is saved
type Emoji struct {
	path       string
	vectorForm [72][72]color.Color
}

//This builds an array of emoji structs that represents every png file in the 72x72 directory
func initEmojiDict() {
	file, err := os.Open("./72x72")
	if err != nil {
		fmt.Println(err)
	}
	//read up to 3000 png files from the directory "72x72"
	names, err := file.Readdirnames(3000)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(names); i++ {
		file, err := os.Open("./72x72/" + names[i])
		if err != nil {
			fmt.Println(err)
		}

		img, err := png.Decode(file)
		if err != nil {
			fmt.Println(err)
		}
		emojiList[i].path = "./72x72/" + names[i]

		size := img.Bounds()

		whiteColor := color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}

		//store vector form for emojis
		for x := size.Min.X; x < size.Max.X; x++ {
			for y := size.Min.Y; y < size.Max.Y; y++ {
				//check if pixel is transparent if it is set it to white
				//pixel := img.At(x, y)
				//_, _, _, a := pixel.RGBA()
				//debug
				if 1 == 0 {
					//if color is transparent set color to white
					emojiList[i].vectorForm[x][y] = whiteColor
				} else {
					emojiList[i].vectorForm[x][y] = img.At(x, y)
				}
			}
		}
	}
	fmt.Println("emoji dict initialized")
}

/*this function draws out the specified emoji onto the output image
emojiIndex: the index of a emoji in emojiDict
img: the image to draw onto
*/
func drawEmoji(img *image.RGBA, emojiIndex int) {
	for x := 0; x < 72; x++ {
		for y := 0; y < 72; y++ {
			img.Set(x+currentSquare.X, y+currentSquare.Y, emojiList[emojiIndex].vectorForm[x][y])
		}
	}
	if currentSquare.Y+72 < resultHeight {
		currentSquare.Y += 72
	} else {
		currentSquare.Y = 0
		currentSquare.X += 72
	}
}

func testDrawSubregion(img *image.RGBA, subsection [][]color.Color) {
	for x := 0; x < 72; x++ {
		for y := 0; y < 72; y++ {
			img.Set(x+currentSquare.X, y+currentSquare.Y, subsection[x][y])
		}
	}
	if currentSquare.Y+72 < resultHeight {
		currentSquare.Y += 72
	} else {
		currentSquare.Y = 0
		currentSquare.X += 72
	}
}

//test this
/*finds the average rgb value of a specified sub square of the image
subsection: an image
x, y: these compose the upperLeft corner of the square region the function finds the average rgb value of
squareSize: the size of the square region
*/
func averageRGBSlice(subsection [][]color.Color, x int, y int, squareSize int) (r, g, b float64) {
	r0, g0, b0 := uint32(0), uint32(0), uint32(0)
	count := 0
	for i := x; i < x+squareSize; i++ {
		for j := y; j < y+squareSize; j++ {
			r1, g1, b1, a := subsection[i][j].RGBA()
			if a != 0 {
				r0, g0, b0 = r0+r1, g0+g1, b0+b1
				count++
			}
		}
	}
	if count == 0 {
		r, g, b = 0xffff, 0xffff, 0xffff
	} else {
		r = float64(r0) / float64(count)
		g = float64(g0) / float64(count)
		b = float64(b0) / float64(count)
	}
	return
}

//test this
/*same as above function but takes arrays as input instead of slices, see averageRGBSlice
 */
func averageRGBArray(subsection [72][72]color.Color, x int, y int, squareSize int) (r, g, b float64) {
	r0, g0, b0 := uint32(0), uint32(0), uint32(0)
	count := 0
	for i := x; i < x+squareSize; i++ {
		for j := y; j < y+squareSize; j++ {
			r1, g1, b1, a := subsection[i][j].RGBA()
			if a != 0 {
				r0, g0, b0 = r0+r1, g0+g1, b0+b1
				count++
			}
		}
	}
	if count == 0 {
		r, g, b = 255, 255, 255
	} else {
		r = float64(r0) / float64(count)
		g = float64(g0) / float64(count)
		b = float64(b0) / float64(count)
	}
	return
}

/*TODO
 */
func nearestEmoji(subsection [][]color.Color, squareSize int) int {
	var smallestDistance float64
	var nearestIndex int
	var distance float64

	//set smallest distance to first emoji
	for eX, sX := 0, 0; eX < len(emojiList[0].vectorForm); eX, sX = eX+outputRatio, sX+inputRatio {
		for eY, sY := 0, 0; eY < len(emojiList[0].vectorForm[eX]); eY, sY = eY+outputRatio, sY+inputRatio {
			r0, g0, b0 := averageRGBSlice(subsection, sX, sY, inputRatio)
			r1, g1, b1 := averageRGBArray(emojiList[0].vectorForm, eX, eY, outputRatio)
			//use sum of square differences
			distance += math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
		}
	}
	smallestDistance = distance
	nearestIndex = 0

	for i := 1; i < len(emojiList); i++ {
		distance = 0
		for eX, sX := 0, 0; eX < len(emojiList[0].vectorForm); eX, sX = eX+outputRatio, sX+inputRatio {
			for eY, sY := 0, 0; eY < len(emojiList[0].vectorForm[eX]); eY, sY = eY+outputRatio, sY+inputRatio {
				r0, g0, b0 := averageRGBSlice(subsection, sX, sY, inputRatio)
				r1, g1, b1 := averageRGBArray(emojiList[i].vectorForm, eX, eY, outputRatio)
				//use sum of square differences
				distance += math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
			}
		}
		if distance < smallestDistance {
			smallestDistance = distance
			nearestIndex = i
		}
	}
	/*
	   //set smallest distance to first emoji
	   for x := 0; x < len(emojiList[0].vectorForm); x++ {
	       for y := 0; y < len(emojiList[0].vectorForm[x]); y++ {
	           r0, g0, b0, a0 := subsection[x][y].RGBA()
	           r1, g1, b1, a1 := emojiList[0].vectorForm[x][y].RGBA()
	           //use sum of square differences
	           distance += math.Pow((r0-r1)**, 2)2 + math.Pow((g0-g1)**, 2)2 + math.Pow((b0-b1)**, 2)2 + math.Pow((a0-a1)**, 2)2
	       }
	   }
	   smallestDistance = distance
	   nearestIndex = 0

	   for i := 1; i < len(emojiList); i++ {
	       distance = 0
	       for x := 0; x < len(emojiList[i].vectorForm); x++ {
	           for y := 0; y < len(emojiList[i].vectorForm[x]); y++ {
	               r0, g0, b0, a0 := subsection[x][y].RGBA()
	               r1, g1, b1, a1 := emojiList[i].vectorForm[x][y].RGBA()
	               //use sum of square differencet
	               distance += math.Pow((r0-r1)**, 2)2 + math.Pow((g0-g1)**, 2)2 + math.Pow((b0-b1)**, 2)2 + math.Pow((a0-a1)**, 2)2
	           }
	       }
	       if distance < smallestDistance {
	               smallestDistance = distance
	               nearestIndex = i
	       }
	   }
	*/
	return nearestIndex
}

/*not sure if necessary anymore still debating on some details of final algo
 */
func findRatio(a, b int) (x, y int) {
	//find gcd
	temp0 := a
	temp1 := b
	for temp0 != temp1 {
		if temp0 < temp1 {
			temp1 -= temp0
		} else {
			temp0 -= temp1
		}
	}

	gcd := temp0

	x = a / gcd
	y = b / gcd

	return
}

/* logic for creating emoji art from image starts here */
func createEmojiArt(img image.Image) {
	size := img.Bounds()

	//hard code for now, square size is 15 pixels long
	squareSize := (size.Max.X - size.Min.X) / numOfEmojisForResult
	inputRatio, outputRatio = findRatio(squareSize, 72)

	resultWidth := 72 * numOfEmojisForResult
	resultHeight := 72 * (size.Max.Y - size.Min.Y) / squareSize
	if (size.Max.Y-size.Min.Y)%squareSize != 0 {
		resultWidth += 72
	}
	var resultImg *image.RGBA = image.NewRGBA(image.Rect(0, 0, resultWidth, resultHeight))

	var upperLeft image.Point
	currentSquare.X = 0
	currentSquare.Y = 0

	//create 2d slice
	subsection := make([][]color.Color, squareSize)
	for i := 0; i < squareSize; i++ {
		subsection[i] = make([]color.Color, squareSize)
	}

	whiteColor := color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}

	for upperLeft.X = size.Min.X; upperLeft.X+squareSize <= size.Max.X; upperLeft.X = upperLeft.X + squareSize {
		for upperLeft.Y = size.Min.Y; upperLeft.Y+squareSize <= size.Max.Y; upperLeft.Y = upperLeft.Y + squareSize {
			for x := upperLeft.X; x < upperLeft.X+squareSize; x++ {
				for y := upperLeft.Y; y < upperLeft.Y+squareSize; y++ {
					pixel := img.At(x, y)
					_, _, _, a := pixel.RGBA()
					if a == 0 {
						subsection[x%squareSize][y%squareSize] = whiteColor
					} else {
						subsection[x%squareSize][y%squareSize] = img.At(x, y)
					}
				}
			}
			closestEmoji := nearestEmoji(subsection, squareSize)
			drawEmoji(resultImg, closestEmoji)
		}
	}

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

	fmt.Println("total distance is ", totalDistance)

}
