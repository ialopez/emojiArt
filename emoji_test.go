package main

import "testing"
import "image/color"

func TestInitEmojiDict(t *testing.T) {
	initEmojiDict()
	platforms := [6]string{"apple", "emojione", "facebook", "facebook-messenger", "google", "twitter"}
	for i := 0; i < len(platforms); i++ {
		for j := 0; j < len(emojiDict[platforms[i]]); j++ {
			if emojiDict[platforms[i]][j].path == "" {
				t.Errorf("path string is empty for the emoji index %d from %s\n", j, platforms[i])
			}
			for x := 0; x < len(emojiDict[platforms[i]][j].vectorForm); x++ {
				for y := 0; y < len(emojiDict[platforms[i]][j].vectorForm[x]); y++ {
					if emojiDict[platforms[i]][j].vectorForm[x][y] == nil {
						t.Errorf("null color for the emoji index %d from %s\n", j, platforms[i])
					}
				}
			}
		}
	}
}

func TestAverageRGBSlice(t *testing.T) {
	square := make([][]color.Color, 3)
	for i := 0; i < 3; i++ {
		square[i] = make([]color.Color, 3)
	}
	whiteColor := color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}
	red := color.NRGBA{R: uint8(255), G: uint8(0), B: uint8(0), A: uint8(255)}
	blue := color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(255), A: uint8(255)}

	square[0][0] = whiteColor
	square[0][1] = whiteColor
	square[0][2] = whiteColor
	square[1][0] = red
	square[1][1] = red
	square[1][2] = red
	square[2][0] = blue
	square[2][1] = blue
	square[2][2] = blue

	r0, g0, b0 := averageRGBSlice(square, 0, 0, 3)
	wr, wg, wb, _ := whiteColor.RGBA()
	rr, rg, rb, _ := red.RGBA()
	br, bg, bb, _ := blue.RGBA()
	r1, g1, b1 := float64((wr+wr+wr+rr+rr+rr+br+br+br))/9, float64((wg+wg+wg+rg+rg+rg+bg+bg+bg))/9, float64((wb+wb+wb+rb+rb+rb+bb+bb+bb))/9
	if (r0 != r1) || (g0 != g1) || (b0 != b1) {
		t.Fail()
	}
}
