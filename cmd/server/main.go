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
	"sync"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	qrcode "github.com/skip2/go-qrcode"
)

// Rate limiter using token bucket algorithm
type RateLimiter struct {
	rate       float64 // tokens per second
	bucketSize float64
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(rate, bucketSize float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     bucketSize,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens = min(rl.bucketSize, rl.tokens+elapsed*rl.rate)
	rl.lastRefill = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// IP-based rate limiter
type IPRateLimiter struct {
	ips    map[string]*RateLimiter
	mu     sync.RWMutex
	rate   float64
	bucket float64
}

func NewIPRateLimiter(rate, bucket float64) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*RateLimiter),
		rate:   rate,
		bucket: bucket,
	}
}

func (rl *IPRateLimiter) getLimiter(ip string) *RateLimiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = NewRateLimiter(rl.rate, rl.bucket)
		rl.ips[ip] = limiter
	}
	return limiter
}

func (rl *IPRateLimiter) Allow(ip string) bool {
	return rl.getLimiter(ip).Allow()
}

// Cache for QR codes
var (
	qrCache      = make(map[string][]byte)
	qrCacheMutex sync.RWMutex
)

// Global IP-based rate limiter: 10 requests per second with a bucket size of 20
var ipRateLimiter = NewIPRateLimiter(10, 20)

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

func getIP(r *http.Request) string {
	// Try to get IP from X-Forwarded-For header first
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Check rate limit per IP
	if !ipRateLimiter.Allow(getIP(r)) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

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
	// Check rate limit per IP
	if !ipRateLimiter.Allow(getIP(r)) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

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

	// Get and validate the shape parameter
	shape := r.URL.Query().Get("shape")
	if shape == "" {
		shape = "square" // default shape
	}
	if shape != "square" && shape != "rectangle" {
		http.Error(w, "Shape must be 'square' or 'rectangle'", http.StatusBadRequest)
		return
	}

	// Get and validate the type parameter
	codeType := r.URL.Query().Get("type")
	if codeType == "" {
		codeType = "qr" // default type
	}
	if codeType != "qr" && codeType != "barcode" {
		http.Error(w, "Type must be 'qr' or 'barcode'", http.StatusBadRequest)
		return
	}

	// Create cache key
	cacheKey := fmt.Sprintf("%s:%d:%s:%s", text, size, shape, codeType)

	// Check cache first
	qrCacheMutex.RLock()
	cachedQR, found := qrCache[cacheKey]
	qrCacheMutex.RUnlock()

	if found {
		if r.URL.Query().Get("base64") == "true" {
			base64Str := base64.StdEncoding.EncodeToString(cachedQR)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(base64Str))
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(cachedQR)
		return
	}

	// Create a buffer to store the PNG
	var buf bytes.Buffer
	var codeImg image.Image

	if codeType == "barcode" {
		// Generate barcode
		bar, err := code128.Encode(text)
		if err != nil {
			http.Error(w, "Failed to generate barcode", http.StatusInternalServerError)
			return
		}

		if shape == "rectangle" {
			// For barcodes, use natural barcode proportions
			codeImg, err = barcode.Scale(bar, size*4, size)
			if err != nil {
				http.Error(w, "Failed to scale barcode", http.StatusInternalServerError)
				return
			}
		} else {
			// Square shape for barcode
			codeImg, err = barcode.Scale(bar, size, size)
			if err != nil {
				http.Error(w, "Failed to scale barcode", http.StatusInternalServerError)
				return
			}
		}
	} else {
		// Generate QR code
		qr, err := qrcode.New(text, qrcode.Medium)
		if err != nil {
			http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
			return
		}

		if shape == "rectangle" {
			// For rectangle shape, use barcode proportions (approx 4:1 ratio)
			qrImg := qr.Image(size)
			width := size * 4
			height := size

			// Create a new rectangular image
			rectImg := image.NewRGBA(image.Rect(0, 0, width, height))

			// Fill background with white
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					rectImg.Set(x, y, color.RGBA{255, 255, 255, 255})
				}
			}

			// Center the QR code in the rectangular image
			offsetX := (width - size) / 2
			offsetY := (height - size) / 2

			// Draw the QR code in the center
			for y := 0; y < size; y++ {
				for x := 0; x < size; x++ {
					c := qrImg.At(x, y)
					rectImg.Set(x+offsetX, y+offsetY, c)
				}
			}

			codeImg = rectImg
		} else {
			// Default square shape
			codeImg = qr.Image(size)
		}
	}

	// Encode the image to PNG
	if err := png.Encode(&buf, codeImg); err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}

	// Store in cache
	qrCacheMutex.Lock()
	qrCache[cacheKey] = buf.Bytes()
	qrCacheMutex.Unlock()

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

func barcodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check rate limit per IP
	if !ipRateLimiter.Allow(getIP(r)) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "Please provide a 'text' parameter", http.StatusBadRequest)
		return
	}

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

	// Get and validate the shape parameter
	shape := r.URL.Query().Get("shape")
	if shape == "" {
		shape = "rectangle" // default shape for barcodes is rectangle
	}
	if shape != "square" && shape != "rectangle" {
		http.Error(w, "Shape must be 'square' or 'rectangle'", http.StatusBadRequest)
		return
	}

	// Generate barcode
	bar, err := code128.Encode(text)
	if err != nil {
		http.Error(w, "Failed to generate barcode", http.StatusInternalServerError)
		return
	}

	// Scale barcode to requested size based on shape
	var scaledBar image.Image
	if shape == "rectangle" {
		// For rectangle shape, use natural barcode proportions (4:1 ratio)
		scaledBar, err = barcode.Scale(bar, size*4, size)
	} else {
		// Square shape
		scaledBar, err = barcode.Scale(bar, size, size)
	}
	if err != nil {
		http.Error(w, "Failed to scale barcode", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaledBar); err != nil {
		http.Error(w, "Failed to encode barcode", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("base64") == "true" {
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(base64Str))
		return
	}

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
	// Register the barcode handler
	http.HandleFunc("/barcode", barcodeHandler)
	// Register the ping handler
	http.HandleFunc("/ping", pingHandler)

	// Start the server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
