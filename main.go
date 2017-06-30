package main

import (
	"fmt" //reader writer
	"image"
	"image/png"
	"io"        //used to write read files
	"io/ioutil" //read files
	"net/http"  //used to handle serve http requests
	"os"        //used to create files in server
)

var mainPage, _ = ioutil.ReadFile("main.html") //read in main page from main.html

func openPNG(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	pic, err := png.Decode(file)
	if err != nil {
		fmt.Println("cannot decode")
		fmt.Println(err)
		return nil
	}
	return pic
}

//serve main page
func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string(mainPage))
}

//serve result page, seen after submitting picture to server
func resultHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("pic")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "<h1>Result %s</h1><div>%s</div>", title, handler.Header)
	f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	img := openPNG("./" + handler.Filename)
	simpleAlgo(img)
}

func main() {
	initEmojiDict()
	initEmojiDictAvg()
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/view/", resultHandler)
	http.ListenAndServe(":8080", nil)
}
