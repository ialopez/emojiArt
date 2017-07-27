package main

import "testing"
import "os"
import "image/png"
import "image/color"

func TestInitEmojiDictAvg(t *testing.T) {
	initEmojiDict()
	initEmojiDictAvg()

}

func TestNearestSimple(t *testing.T) {
	square := make([][]color.Color, emojiSize)
	for i := 0; i < emojiSize; i++ {
		square[i] = make([]color.Color, emojiSize)
	}

	expectedPath := "./apple/00a9.png"

	file, err := os.Open(expectedPath)
	if err != nil {
		t.Error(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		t.Error(err)
	}

	for x := 0; x < emojiSize; x++ {
		for y := 0; y < emojiSize; y++ {
			square[x][y] = img.At(x, y)
		}
	}
	currentPlatform = "apple"
	closestEmoji := nearestSimple(square, emojiSize)
	if emojiDict[currentPlatform][closestEmoji].path != "./apple/00a9.png" {
		t.Errorf("expected %s, instead got %s\n", expectedPath, emojiDict[currentPlatform][closestEmoji].path)
	}
	for x := 0; x < emojiSize; x++ {
		for y := 0; y < emojiSize; y++ {
			if emojiDict[currentPlatform][closestEmoji].vectorForm[x][y] != img.At(x, y) {
				t.Fail()
			}
		}
	}

}

func TestDistance(t *testing.T) {
	copyRightEmoji := make([][]color.Color, emojiSize)
	for i := 0; i < emojiSize; i++ {
		copyRightEmoji[i] = make([]color.Color, emojiSize)
	}

	pathCopy := "./apple/00a9.png"

	file, err := os.Open(pathCopy)
	if err != nil {
		t.Error(err)
	}
	imgCopy, err := png.Decode(file)
	if err != nil {
		t.Error(err)
	}

	for x := 0; x < emojiSize; x++ {
		for y := 0; y < emojiSize; y++ {
			copyRightEmoji[x][y] = imgCopy.At(x, y)
		}
	}

	swimmerEmoji := make([][]color.Color, emojiSize)
	for i := 0; i < emojiSize; i++ {
		swimmerEmoji[i] = make([]color.Color, emojiSize)
	}

	pathSwimmer := "./apple/1f3ca-1f3fd.png"

	file, err = os.Open(pathSwimmer)
	if err != nil {
		t.Error(err)
	}
	imgSwimmer, err := png.Decode(file)
	if err != nil {
		t.Error(err)
	}

	for x := 0; x < emojiSize; x++ {
		for y := 0; y < emojiSize; y++ {
			swimmerEmoji[x][y] = imgSwimmer.At(x, y)
		}
	}

	//test if emojis are the same for some reason
	same := true
	for x := 0; x < emojiSize; x++ {
		for y := 0; y < emojiSize; y++ {
			if copyRightEmoji[x][y] != swimmerEmoji[x][y] {
				same = false
			}
		}
	}
	if same {
		t.Error("emojis are the same")
	}

	r0, g0, b0 := averageRGBSlice(copyRightEmoji, 0, 0, emojiSize)
	r1, g1, b1 := averageRGBSlice(swimmerEmoji, 0, 0, emojiSize)

	copyIndex := -1
	for i := 0; i < len(emojiDict["apple"]); i++ {
		if emojiDict["apple"][i].path == pathCopy {
			copyIndex = i
		}
	}

	r2, g2, b2 := emojiDictAvg["apple"][copyIndex][0], emojiDictAvg["apple"][copyIndex][1], emojiDictAvg["apple"][copyIndex][2]

	if r0 != r2 || g0 != g2 || b0 != b2 {
		t.Errorf("expected %f, %f, %f, instead got %f, %f, %f\n", r0, g0, b0, r2, g2, b2)
	}

	if r0 != r1 || g0 != g1 || b0 != b1 {
		t.Errorf("expected %f, %f, %f, instead got %f, %f, %f\n", r0, g0, b0, r1, g1, b1)
	}
}
