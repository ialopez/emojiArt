package emojiart

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"
)

//calls nearest Neighbor and return index
func nearestEmoji(subsection imageSlice, platform string) int {
	r0, g0, b0 := subsection.averageRGB(0, 0, len(subsection), false)
	//set smallest distance to first emoji
	_, nearestIndex := nearestNeighbor(r0, g0, b0, 0, kdTree[platform])

	return nearestIndex
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
					//find nearest emoji for subsection of input image and save result in mapping
					closestEmoji := nearestEmoji(subsection, p.outputPlatform)
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
