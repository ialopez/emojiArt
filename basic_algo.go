package emojiart

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
)

func (p picToEmoji) nearestSimple(subsection imageSlice) int {
	var smallestDistance float64
	var nearestIndex int
	var distance float64

	//set smallest distance to first emoji
	r0, g0, b0 := subsection.averageRGB(0, 0, p.squareSize, false)
	r1, g1, b1 := emojiDictAvg[p.outputPlatform][0][0], emojiDictAvg[p.outputPlatform][0][1], emojiDictAvg[p.outputPlatform][0][2]

	//use sum of square differences
	smallestDistance = math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
	nearestIndex = 0

	for i := 1; i < len(emojiDictAvg[p.outputPlatform]); i++ {
		r1, g1, b1 := emojiDictAvg[p.outputPlatform][i][0], emojiDictAvg[p.outputPlatform][i][1], emojiDictAvg[p.outputPlatform][i][2]
		distance = math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
		if distance < smallestDistance {
			smallestDistance = distance
			nearestIndex = i
		}
	}

	return nearestIndex
}

/*returns an emojiArt Image of the input image field contained in the picToEmoji struct
 */
func (p picToEmoji) basicAlgo() image.Image {
	//hardcode square size for now
	/*fmt.Println("enter square size")
	var squareSize int
	fmt.Scanf("%d\n", &squareSize)
	*/

	imgWidth := p.inputImage.Bounds().Max.X - p.inputImage.Bounds().Min.X
	imgHeight := p.inputImage.Bounds().Max.Y - p.inputImage.Bounds().Min.X
	p.resultWidth = EMOJI_SIZE * imgWidth / p.squareSize
	p.resultHeight = EMOJI_SIZE * imgHeight / p.squareSize
	p.outputImage = image.NewRGBA(image.Rect(0, 0, p.resultWidth, p.resultHeight))

	//create 2d slice
	subsection := make([][]color.Color, p.squareSize)
	for i := 0; i < p.squareSize; i++ {
		subsection[i] = make([]color.Color, p.squareSize)
	}

	var upperLeft image.Point

	for upperLeft.X = p.inputImage.Bounds().Min.X; upperLeft.X < p.inputImage.Bounds().Max.X; upperLeft.X += p.squareSize {
		for upperLeft.Y = p.inputImage.Bounds().Min.Y; upperLeft.Y < p.inputImage.Bounds().Max.Y; upperLeft.Y += p.squareSize {
			for x := upperLeft.X; x < upperLeft.X+p.squareSize; x++ {
				for y := upperLeft.Y; y < upperLeft.Y+p.squareSize; y++ {
					subsection[x%p.squareSize][y%p.squareSize] = p.inputImage.At(x, y)
				}
			}
			//closestEmoji := p.nearestSimple(subsection)
			//p.drawEmoji(closestEmoji)
		}
		fmt.Print(upperLeft.X)
		fmt.Print("out of ")
		fmt.Println(p.inputImage.Bounds().Max.X)
	}

	fmt.Println("drawing image")
	f, err := os.Create("./image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, p.outputImage); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("finished")

	return p.outputImage

}

/*similar to function above instead returns a struct containing a compressed representation of the emojiArt image
 */
func (p picToEmoji) basicAlgoGenMap() *emojiMap {
	imgWidth := p.inputImage.Bounds().Max.X - p.inputImage.Bounds().Min.X
	imgHeight := p.inputImage.Bounds().Max.Y - p.inputImage.Bounds().Min.X
	width := imgWidth / p.squareSize
	height := imgHeight / p.squareSize

	resultMap := newEmojiMap(width, height)

	//create 2d slice
	subsection := make([][]color.Color, p.squareSize)
	for i := 0; i < p.squareSize; i++ {
		subsection[i] = make([]color.Color, p.squareSize)
	}

	var upperLeft image.Point

	a := 0
	for upperLeft.X = p.inputImage.Bounds().Min.X; upperLeft.X+p.squareSize < p.inputImage.Bounds().Max.X; upperLeft.X += p.squareSize {
		b := 0
		for upperLeft.Y = p.inputImage.Bounds().Min.Y; upperLeft.Y+p.squareSize < p.inputImage.Bounds().Max.Y; upperLeft.Y += p.squareSize {
			for x := upperLeft.X; x < upperLeft.X+p.squareSize; x++ {
				for y := upperLeft.Y; y < upperLeft.Y+p.squareSize; y++ {
					subsection[x%p.squareSize][y%p.squareSize] = p.inputImage.At(x, y)
				}
			}
			closestEmoji := p.nearestSimple(subsection)
			//fmt.Printf("a = %d, b = %d\n", a, b)
			//reversed the order of indexes "a" and "b" because images were coming out rotated
			resultMap.Mapping[b][a] = closestEmoji
			b++
		}
		fmt.Print(upperLeft.X)
		fmt.Print("out of ")
		fmt.Println(p.inputImage.Bounds().Max.X)
		a++
	}

	//create dictionary with the URL paths of all the emojis used in resultMap.Mapping
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			key := strconv.Itoa(resultMap.Mapping[i][j])
			if _, contains := resultMap.Dictionary[key]; !contains {
				resultMap.Dictionary[key] = emojiURLPath[p.outputPlatform][resultMap.Mapping[i][j]]
			}
		}
	}

	fmt.Println("finished")

	return resultMap
}
