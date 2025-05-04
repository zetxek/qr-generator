package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

func parseHexColor(s string) (color.RGBA, error) {
	c := color.RGBA{A: 255}
	if strings.HasPrefix(s, "#") {
		s = s[1:]
	}
	if len(s) != 6 {
		return c, fmt.Errorf("invalid color length")
	}
	_, err := fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	return c, err
}

func generateImage(size int, c1, c2 color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			ratio := float64(x) / float64(size-1)
			r := uint8(float64(c1.R)*(1-ratio) + float64(c2.R)*ratio)
			g := uint8(float64(c1.G)*(1-ratio) + float64(c2.G)*ratio)
			b := uint8(float64(c1.B)*(1-ratio) + float64(c2.B)*ratio)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse size
	size := 200
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s >= 10 && s <= 2000 {
			size = s
		}
	}

	// Parse colors
	c1, err1 := parseHexColor(r.URL.Query().Get("color1"))
	if err1 != nil {
		c1 = color.RGBA{R: 0, G: 0, B: 255, A: 255} // default blue
	}
	c2, err2 := parseHexColor(r.URL.Query().Get("color2"))
	if err2 != nil {
		c2 = color.RGBA{R: 255, G: 0, B: 0, A: 255} // default red
	}

	img := generateImage(size, c1, c2)

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		return
	}
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	// Get the text parameter from the query string
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "Please provide a 'text' parameter", http.StatusBadRequest)
		return
	}

	// Get and validate the size parameter
	size := 256 // default size
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "Size must be a valid number", http.StatusBadRequest)
			return
		}
		if size < 50 || size > 1000 {
			http.Error(w, "Size must be between 50 and 1000 pixels", http.StatusBadRequest)
			return
		}
	}

	// Generate QR code
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	// Create a buffer to store the PNG
	var buf bytes.Buffer
	if err := qr.Write(size, &buf); err != nil {
		http.Error(w, "Failed to write QR code", http.StatusInternalServerError)
		return
	}

	// Check if base64 encoding is requested
	if r.URL.Query().Get("base64") == "true" {
		// Encode the image to base64
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(base64Str))
		return
	}

	// If not base64, return the PNG image
	w.Header().Set("Content-Type", "image/png")
	w.Write(buf.Bytes())
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	// Register the image handler
	http.HandleFunc("/image", imageHandler)
	// Register the QR code handler
	http.HandleFunc("/qr", qrHandler)
	// Register the ping handler
	http.HandleFunc("/ping", pingHandler)

	// Start the server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
