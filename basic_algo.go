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
	"sync"
)

func (p picToEmoji) nearestSimple(subsection imageSlice) int {
	var smallestDistance float64
	var nearestIndex int
	var distance float64

	//set smallest distance to first emoji
	r0, g0, b0 := subsection.averageRGB(0, 0, p.squareSize, false)
	r1, g1, b1 := emojiDictAvg[p.outputPlatform][0].R, emojiDictAvg[p.outputPlatform][0].G, emojiDictAvg[p.outputPlatform][0].B

	//use sum of square differences
	smallestDistance = math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
	nearestIndex = emojiDictAvg[p.outputPlatform][0].URLIndex

	for i := 1; i < len(emojiDictAvg[p.outputPlatform]); i++ {
		r1, g1, b1 := emojiDictAvg[p.outputPlatform][i].R, emojiDictAvg[p.outputPlatform][i].G, emojiDictAvg[p.outputPlatform][i].B
		distance = math.Pow(r0-r1, 2) + math.Pow(g0-g1, 2) + math.Pow(b0-b1, 2)
		if distance < smallestDistance {
			smallestDistance = distance
			nearestIndex = emojiDictAvg[p.outputPlatform][i].URLIndex
		}
	}

	return nearestIndex
}

func (p picToEmoji) nearestSimpleTree(subsection imageSlice) int {
	r0, g0, b0 := subsection.averageRGB(0, 0, p.squareSize, false)
	//set smallest distance to first emoji
	_, nearestIndex := nearestNeighbor(r0, g0, b0, 0, kdTree[p.outputPlatform])

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

	/*image is divided into slices based on number of threads being used like so
	each . is a pixel and = represents a border
	.....................
	.....................
	=====================
	.....................
	.....................
	=====================
	.....................
	.....................

	each slice has its output calculated on its on thread
	*/

	heightPerSection := (height / NUM_OF_THREADS) * p.squareSize //pixel height of each section

	var wg sync.WaitGroup
	for lowBoundY := p.inputImage.Bounds().Min.Y; lowBoundY < p.inputImage.Bounds().Max.Y; lowBoundY += heightPerSection {
		wg.Add(1)
		minY := lowBoundY
		initA := lowBoundY / p.squareSize
		var maxY int
		if lowBoundY+heightPerSection > p.inputImage.Bounds().Max.Y {
			//case where there is left over section in the image that has a smaller height than heightPerSection
			maxY = p.inputImage.Bounds().Max.Y
		} else {
			maxY = lowBoundY + heightPerSection
		}

		go func(minY, maxY, initA int, p *picToEmoji, wg *sync.WaitGroup) {
			var upperLeft image.Point
			//create 2d slice
			subsection := make([][]color.Color, p.squareSize)
			for i := 0; i < p.squareSize; i++ {
				subsection[i] = make([]color.Color, p.squareSize)
			}

			a := initA
			//divide image into blocks and find the nearest emoji for each block
			for upperLeft.Y = minY; upperLeft.Y+p.squareSize <= maxY; upperLeft.Y += p.squareSize {
				b := 0
				for upperLeft.X = p.inputImage.Bounds().Min.X; upperLeft.X+p.squareSize <= p.inputImage.Bounds().Max.X; upperLeft.X += p.squareSize {
					for y := upperLeft.Y; y < upperLeft.Y+p.squareSize; y++ {
						for x := upperLeft.X; x < upperLeft.X+p.squareSize; x++ {
							subsection[y%p.squareSize][x%p.squareSize] = p.inputImage.At(x, y)
						}
					}
					closestEmoji := (*p).nearestSimpleTree(subsection)
					resultMap.Mapping[a][b] = closestEmoji
					b++
				}
				a++
			}
			(*wg).Done()
		}(minY, maxY, initA, &p, &wg)
	}

	wg.Wait() //function waits here until all goroutines are finished

	//create dictionary with the URL paths of all the emojis used in resultMap.Mapping
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			key := strconv.Itoa(resultMap.Mapping[i][j])
			if _, contains := resultMap.Dictionary[key]; !contains {
				urlPath := emojiURLPath[p.outputPlatform][resultMap.Mapping[i][j]]
				resultMap.Dictionary[key] = urlPath
			}
		}
	}

	fmt.Println("finished")

	return resultMap
}
