package main

import "testing"
import "image/color"

func TestInitEmojiDict(t *testing.T) {
	initEmojiDict()
	for i := 0; i < len(emojiList); i++ {
		if emojiList[i].path == "" {
			t.Fail()
		}
		for x := 0; x < len(emojiList[i].vectorForm); x++ {
			for y := 0; y < len(emojiList[i].vectorForm[x]); y++ {
				if emojiList[i].vectorForm[x][y] == nil {
					t.Fail()
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
	wr, wg, wb, wa := whiteColor.RGBA()
	rr, rg, rb, ra := red.RGBA()
	br, bg, bb, ba := blue.RGBA()
	r1, g1, b1 := (wr+wr+wr+rr+rr+rr+br+br+br)/9, (wg+wg+wg+rg+rg+rg+bg+bg+bg)/9, (wb+wb+wb+rb+rb+rb+bb+bb+bb)/9
	if (r0 != r1) | (g0 != g1) | (b0 != b1) {
		t.Fail()
	}
}
