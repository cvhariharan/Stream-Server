package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"time"

	// "github.com/cvhariharan/screenshot"
	"github.com/kbinani/screenshot"
)

type Mjpeg struct {
}

var all image.Rectangle = image.Rect(0, 0, 0, 0)

func (m *Mjpeg) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mimeWriter := multipart.NewWriter(w)
	contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
	w.Header().Add("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// i := 0
	for {
		// start := time.Now()
		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Type", "image/jpeg")

		partWriter, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			log.Printf("Could not create Part Writer: %v\n", err)
			break
		}

		// imgName := strconv.Itoa(i) + ".jpg"
		data, err := screenshot.Capture(all.Min.X, all.Min.Y, all.Dx(), all.Dy())
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, data, nil)
		img := buf.Bytes()

		if err != nil {
			log.Printf("Could not get the next frame: %v\n", err)
			break
		}
		partHeader.Add("Content-Length", strconv.Itoa(len(img)))
		if _, err = io.Copy(partWriter, bytes.NewReader(img)); err != nil {
			log.Printf("Could not write the image to the response: %v\n", err)
			break
		}
		time.Sleep(16 * time.Millisecond)
		// end := time.Now()
		// i = (i + 1) % 300
		// dif := end.Sub(start)
		// time.Sleep()
		// os.Remove(imgName)

	}
}

func main() {
	// go Run()
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("Active display not found")
	}

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		all = bounds.Union(all)
	}
	mux := http.NewServeMux()
	v := &Mjpeg{}
	mux.Handle("/", v)
	http.ListenAndServe(":8080", mux)
}
