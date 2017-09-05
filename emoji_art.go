package emojiart

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

var (
	emojiDictAvg   map[string][]Emoji
	emojiURLPath   map[string][]string
	platforms      = [6]string{"apple", "emojione", "facebook", "facebook-messenger", "google", "twitter"}
	NUM_OF_THREADS int
)

const (
	EMOJI_SIZE              = 64
	KD_TREE_JSON_LOC        = "../emojiart/kdTree.json"
	EMOJI_URL_PATH_JSON_LOC = "../emojiart/emojiURLPath.json"
)

type emojiImage [EMOJI_SIZE][EMOJI_SIZE]color.Color
type imageSlice [][]color.Color

type Emoji struct {
	R, G, B  float64
	URLIndex int
}

type picToEmoji struct {
	squareSize, resultWidth, resultHeight int
	outputPlatform                        string
	inputImage                            image.Image
	outputImage                           *image.RGBA
	currentSquare                         image.Point
}

type emojiMap struct {
	Dictionary map[string]string `json:"dictionary"`
	Mapping    [][]int           `json:"mapping"`
}

func NewPicToEmoji(squareSize int, outputPlatform string, inputImage image.Image) *picToEmoji {
	p := new(picToEmoji)
	p.squareSize = squareSize
	p.outputPlatform = outputPlatform
	p.inputImage = inputImage
	p.currentSquare = image.Point{X: 0, Y: 0}
	return p
}

func newEmojiMap(width, height int) *emojiMap {
	e := new(emojiMap)
	e.Dictionary = make(map[string]string)
	e.Mapping = make([][]int, height)
	for i := 0; i < height; i++ {
		e.Mapping[i] = make([]int, width)
	}
	return e
}

/*this builds emojiDictAvg(holds the average rgb values of emoji images) and emojiURLPath(holds the url addresses to images on server)
emojiDictAvg is then uses to build the kdTree and is freed once its built
this function can also rebuild emojiURLPath and kdTree from a json file if they were previously saved
*/
func InitDataStructs(recalculateEmojiDict bool) {

	//read previously saved emoji dict from json files
	if !recalculateEmojiDict {
		buffer, err := ioutil.ReadFile(EMOJI_URL_PATH_JSON_LOC)
		if err == nil {
			err = json.Unmarshal(buffer, &emojiURLPath)
			if err == nil {
				buffer, err := ioutil.ReadFile(KD_TREE_JSON_LOC)
				if err == nil {
					err = json.Unmarshal(buffer, &kdTree)
					if err == nil {
						fmt.Println("kd tree initialized from json")
						return
					}
				}
			}
		}
	}

	//if recalculateEmojiDict is true or an error occured above make dictionary from scratch
	//emojiDict = make(map[string][]emoji)
	emojiDictAvg = make(map[string][]Emoji)
	emojiURLPath = make(map[string][]string)
	kdTree = make(map[string]*Tree)

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

		emojiDictAvg[platforms[i]] = make([]Emoji, len(names))
		emojiURLPath[platforms[i]] = make([]string, len(names))

		var currentEmoji emojiImage

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
			emojiDictAvg[platforms[i]][j].URLIndex = j
			emojiURLPath[platforms[i]][j] = "/images/" + platforms[i] + "/" + names[j]

			size := img.Bounds()

			//store vector form for emojis
			for x := size.Min.X; x < size.Max.X; x++ {
				for y := size.Min.Y; y < size.Max.Y; y++ {
					currentEmoji[x][y] = img.At(x, y)
				}
			}
			r, g, b := currentEmoji.averageRGB(0, 0, EMOJI_SIZE, false)
			emojiDictAvg[platforms[i]][j].R = r
			emojiDictAvg[platforms[i]][j].G = g
			emojiDictAvg[platforms[i]][j].B = b
		}
	}

	for i := 0; i < len(platforms); i++ {
		kdTree[platforms[i]] = BuildTree(platforms[i], 0, len(emojiDictAvg[platforms[i]])-1, 0)
	}

	emojiDictAvg = nil //we can free emojiDictavg once the kdtree is built

	outputDataStructsJSON()

	fmt.Println("kd tree initialized from scratch")
}

/*finds the average rgb value of a specified sub square of the image
subsection: an image
x, y: these compose the coordinates for the upperleft corner of the square region the function will operate on
squareSize: the length of the square region
*/
func (s *imageSlice) averageRGB(x int, y int, squareSize int, ignoreTransparentPixels bool) (r, g, b float64) {
	r0, g0, b0 := uint32(0), uint32(0), uint32(0)
	count := 0
	for i := x; i < x+squareSize; i++ {
		for j := y; j < y+squareSize; j++ {
			r1, g1, b1, a := (*s)[i][j].RGBA()
			if a != 0 {
				r0, g0, b0 = r0+r1, g0+g1, b0+b1
				count++
			} else if !ignoreTransparentPixels {
				//add the color white to the sum
				r0, g0, b0 = r0+0xffff, g0+0xffff, b0+0xffff
				count++
			}
		}
	}
	// if all pixels are transparent return white as average color
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
func (s *emojiImage) averageRGB(x int, y int, squareSize int, ignoreTransparentPixels bool) (r, g, b float64) {
	r0, g0, b0 := uint32(0), uint32(0), uint32(0)
	count := 0
	for i := x; i < x+squareSize; i++ {
		for j := y; j < y+squareSize; j++ {
			if s[i][j] == nil {
				fmt.Println("null color")
			}
			r1, g1, b1, a := s[i][j].RGBA()
			if a != 0 {
				r0, g0, b0 = r0+r1, g0+g1, b0+b1
				count++
			} else if !ignoreTransparentPixels {
				//add the color white to the sum
				r0, g0, b0 = r0+0xffff, g0+0xffff, b0+0xffff
				count++
			}
		}
	}
	// if all pixels are transparent return white as average color
	if count == 0 {
		r, g, b = 0xffff, 0xffff, 0xffff
	} else {
		r = float64(r0) / float64(count)
		g = float64(g0) / float64(count)
		b = float64(b0) / float64(count)
	}
	return
}

func InitNumOfThreads() {
	NUM_OF_THREADS = runtime.NumCPU() + 1
	runtime.GOMAXPROCS(NUM_OF_THREADS - 1)
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

//outputs kdTree and emojiURLPath to json file
func outputDataStructsJSON() {
	output, _ := json.Marshal(kdTree)
	file, err := os.Create(KD_TREE_JSON_LOC)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write(output)

	output, _ = json.Marshal(&emojiURLPath)
	file, err = os.Create(EMOJI_URL_PATH_JSON_LOC)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write(output)
}

func (p picToEmoji) CreateEmojiArtMap() *emojiMap {
	return p.basicAlgoGenMap()
}
