package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"

	qrcode "github.com/skip2/go-qrcode"
)

func generateImage() image.Image {
	// Create a 200x200 image
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	// Fill the image with a gradient
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			// Create a gradient from blue to red
			img.Set(x, y, color.RGBA{
				R: uint8(x),
				G: uint8(y),
				B: 255,
				A: 255,
			})
		}
	}

	return img
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Generate the image
	img := generateImage()

	// Set the content type header
	w.Header().Set("Content-Type", "image/png")

	// Encode the image to PNG and write it to the response
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

	// Generate QR code
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", "image/png")

	// Write the QR code as PNG
	if err := qr.Write(256, w); err != nil {
		http.Error(w, "Failed to write QR code", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Register the image handler
	http.HandleFunc("/image", imageHandler)
	// Register the QR code handler
	http.HandleFunc("/qr", qrHandler)

	// Start the server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
