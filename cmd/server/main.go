package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
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

func main() {
	// Register the image handler
	http.HandleFunc("/image", imageHandler)

	// Start the server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
