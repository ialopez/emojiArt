package emojiart

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

/*this is a map that contains array of emoji structs
each key is the platform the emojis belong to ie. apple, facebook
*/
var emojiDict map[string][]emoji
var emojiDictAvg map[string][][3]float64
var platforms = [6]string{"apple", "emojione", "facebook", "facebook-messenger", "google", "twitter"}
var emojiURLPath map[string][]string

const emojiSize = 64

//emoji png files are saved internally in the program as 64x64 arrays of colors, also the path of the emoji is saved
type emoji struct {
	urlpath    string
	vectorForm [emojiSize][emojiSize]color.Color
}

type picToEmoji struct {
	squareSize, resultWidth, resultHeight int
	outputPlatform                        string
	useAdvancedAlgo                       bool
	inputImage                            image.Image
	outputImage                           *image.RGBA
	currentSquare                         image.Point
}

type emojiMap struct {
	Dictionary map[string]string `json:"dictionary"`
	Mapping    [][]int           `json:"mapping"`
}

func NewPicToEmoji(squareSize int, outputPlatform string, useAdvancedAlgo bool, inputImage image.Image) *picToEmoji {
	p := new(picToEmoji)
	p.squareSize = squareSize
	p.outputPlatform = outputPlatform
	p.useAdvancedAlgo = useAdvancedAlgo
	p.inputImage = inputImage
	p.currentSquare = image.Point{X: 0, Y: 0}
	return p
}

func newEmojiMap(width, height int) *emojiMap {
	e := new(emojiMap)
	e.Dictionary = make(map[string]string)
	e.Mapping = make([][]int, width)
	for i := 0; i < width; i++ {
		e.Mapping[i] = make([]int, height)
	}
	return e
}

/*this Dictionary holds the precomputed average RGB values of Emojis
it is a map that takes in the platform as a key ie.. "apple" "facebook" and returns a [][3] slice of floats
where [][0] is red [][1] is green and [][2] is blue
*/

func InitEmojiDictAvg() {
	//emojiDict = make(map[string][]emoji)
	emojiURLPath = make(map[string][]string)
	emojiDictAvg = make(map[string][][3]float64)

	for i := 0; i < len(platforms); i++ {
		currentDir := "../emojiart/" + platforms[i] + "/"
		folder, err := os.Open(currentDir)
		if err != nil {
			log.Fatal(err)
		}
		defer folder.Close()

		//read up to 3000 png files from the current folder
		names, err := folder.Readdirnames(3000)
		if err != nil {
			log.Fatal(err)
		}

		emojiURLPath[platforms[i]] = make([]string, len(names))
		emojiDictAvg[platforms[i]] = make([][3]float64, len(names))
		//emojiDict[platforms[i]] = make([]emoji, len(names))

		var currentEmoji [emojiSize][emojiSize]color.Color

		for j := 0; j < len(names); j++ {
			file, err := os.Open(currentDir + names[j])
			if err != nil {
				log.Fatal(err)
			}
			img, err := png.Decode(file)
			file.Close()
			if err != nil {
				log.Fatal(err)
			}
			emojiURLPath[platforms[i]][j] = "/images/" + platforms[i] + "/" + names[j]

			size := img.Bounds()

			//store vector form for emojis
			for x := size.Min.X; x < size.Max.X; x++ {
				for y := size.Min.Y; y < size.Max.Y; y++ {
					currentEmoji[x][y] = img.At(x, y)
				}
			}
			r, g, b := averageRGBArray(currentEmoji, 0, 0, emojiSize)
			emojiDictAvg[platforms[i]][j][0] = r
			emojiDictAvg[platforms[i]][j][1] = g
			emojiDictAvg[platforms[i]][j][2] = b
		}
	}
	fmt.Println("emoji dict initialized")
}

/*this function draws out the specified emoji onto the output image
emojiIndex: the index of a emoji in emojiDict
img: the image to draw onto
*/
/*
func (p *picToEmoji) drawEmoji(emojiIndex int) {
	for x := 0; x < emojiSize; x++ {
		for y := 0; y < emojiSize; y++ {
			p.outputImage.Set(x+p.currentSquare.X, y+p.currentSquare.Y, emojiDict[p.outputPlatform][emojiIndex].vectorForm[x][y])
		}
	}
	if p.currentSquare.Y+emojiSize < p.resultHeight {
		p.currentSquare.Y += emojiSize
	} else {
		p.currentSquare.Y = 0
		p.currentSquare.X += emojiSize
	}
}
*/

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

/*same as above function but takes arrays as input instead of slices, see averageRGBSlice
 */
func averageRGBArray(subsection [64][64]color.Color, x int, y int, squareSize int) (r, g, b float64) {
	r0, g0, b0 := uint32(0), uint32(0), uint32(0)
	count := 0
	for i := x; i < x+squareSize; i++ {
		for j := y; j < y+squareSize; j++ {
			if subsection[i][j] == nil {
				fmt.Println("null color")
			}
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
func (p picToEmoji) CreateEmojiArt() image.Image {
	if p.useAdvancedAlgo {
		return p.advancedAlgo()
	} else {
		return p.basicAlgo()
	}
}

//debug purposes
func (p picToEmoji) DrawInputImage() {
	fmt.Println("drawing image")
	f, err := os.Create("./image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, p.inputImage); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("finished")
}

func (p picToEmoji) CreateEmojiArtMap() *emojiMap {
	return p.basicAlgoGenMap()
}
