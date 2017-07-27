package main

import (
	"bytes"
	"encoding/base64"
	"fmt" //reader writer
	"html/template"
	"image"
	"image/png"
	"io"        //used to write read files
	"io/ioutil" //read files
	"log"
	"net/http" //used to handle serve http requests
	"os"       //used to create files in server
	"strconv"
)

var mainPage, _ = ioutil.ReadFile("main.html") //read in main page from main.html
var resultPage = template.Must(template.ParseFiles("result.html"))

func openPNG(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	pic, err := png.Decode(file)
	if err != nil {
		fmt.Println("cannot decode")
		log.Fatal(err)
	}
	return pic
}

//serve main page
func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string(mainPage))
}

//serve result page, seen after submitting picture to server
func resultHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("pic")
	temp := r.FormValue("squareSize")
	squareSize, err = strconv.Atoi(temp)
	if err != nil {
		log.Fatal(err)
	}
	currentPlatform = r.FormValue("platform")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(f, file)
	img := openPNG("./" + handler.Filename)
	resultImg := simpleAlgo(img)

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, resultImg)
	if err != nil {
		log.Fatal(err)
	}

	title := handler.Filename
	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	data := map[string]string{"Image": str, "Title": title}
	err = resultPage.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	currentPlatform = "apple"

	initEmojiDict()
	initEmojiDictAvg()
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/view/", resultHandler)
	http.ListenAndServe(":8080", nil)
}
