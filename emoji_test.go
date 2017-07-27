package main

import "testing"
import "image/color"
import "os"
import "image/png"

func TestInitEmojiDict(t *testing.T) {
	platforms := [6]string{"apple", "emojione", "facebook", "facebook-messenger", "google", "twitter"}
	for i := 0; i < len(platforms); i++ {
		for j := 0; j < len(emojiDict[platforms[i]]); j++ {
			if emojiDict[platforms[i]][j].path == "" {
				t.Errorf("path string is empty for the emoji index %d from %s\n", j, platforms[i])
			}
			file, err := os.Open(emojiDict[platforms[i]][j].path)
			if err != nil {
				t.Error(err)
			}
			img, err := png.Decode(file)
			if err != nil {
				t.Error(err)
			}

			for x := 0; x < len(emojiDict[platforms[i]][j].vectorForm); x++ {
				for y := 0; y < len(emojiDict[platforms[i]][j].vectorForm[x]); y++ {
					if emojiDict[platforms[i]][j].vectorForm[x][y] == nil {
						t.Errorf("null color for the emoji index %d from %s\n", j, platforms[i])
					}
					if emojiDict[platforms[i]][j].vectorForm[x][y] != img.At(x, y) {
						t.Fail()
					}
				}
			}
			file.Close()
		}
	}
}

func TestAverageRGBSlice(t *testing.T) {
	square := make([][]color.Color, emojiSize)
	for i := 0; i < emojiSize; i++ {
		square[i] = make([]color.Color, emojiSize)
	}

	file, err := os.Open("./apple/00a9.png")
	if err != nil {
		t.Error(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		t.Error(err)
	}

	subsection := make([][]color.Color, 64)
	for i := 0; i < 64; i++ {
		subsection[i] = make([]color.Color, 64)
	}

	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			subsection[x][y] = img.At(x, y)
		}
	}
	currentPlatform = "apple"
	closestEmoji := nearestSimple(subsection, 64)
	for x := 0; x < emojiSize; x++ {
		for y := 0; y < emojiSize; y++ {
			if emojiDict[currentPlatform][closestEmoji].vectorForm[x][y] != img.At(x, y) {
				t.Fail()
			}
		}
	}
}
