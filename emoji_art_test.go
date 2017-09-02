package emojiart

import "testing"
import "image"
import "os"
import _ "image/png"

func TestInitEmojiDict(t *testing.T) {
	InitEmojiDictAvg(false)
}

func TestAverageRGBSlice(t *testing.T) {
}

func BenchmarkCreateEmojiMap(b *testing.B) {
	file, _ := os.Open("../emojiart/obama.png")
	img, _, _ := image.Decode(file)
	picToEmoji := NewPicToEmoji(22, "apple", img)
	b.ResetTimer()
	picToEmoji.CreateEmojiArtMap()
}
