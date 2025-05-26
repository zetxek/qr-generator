package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()
	pingHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "hello" {
		t.Fatalf("expected body 'hello', got '%s'", rr.Body.String())
	}
}

func TestImageHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/image", nil)
	rr := httptest.NewRecorder()
	imageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestImageHandler_Defaults(t *testing.T) {
	req := httptest.NewRequest("GET", "/image", nil)
	rr := httptest.NewRecorder()
	imageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestImageHandler_CustomSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/image?size=123", nil)
	rr := httptest.NewRecorder()
	imageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestImageHandler_CustomColors(t *testing.T) {
	req := httptest.NewRequest("GET", "/image?color1=00ff00&color2=ff00ff", nil)
	rr := httptest.NewRecorder()
	imageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestImageHandler_InvalidSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/image?size=notanumber", nil)
	rr := httptest.NewRecorder()
	imageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestImageHandler_InvalidColors(t *testing.T) {
	req := httptest.NewRequest("GET", "/image?color1=badcolor&color2=alsoBad", nil)
	rr := httptest.NewRecorder()
	imageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestQRHandler_PNG(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestQRHandler_Base64(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello&base64=true", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/plain" {
		t.Fatalf("expected Content-Type text/plain, got %s", ct)
	}
	// Should be valid base64
	if _, err := base64.StdEncoding.DecodeString(rr.Body.String()); err != nil {
		t.Fatalf("expected valid base64, got error: %v", err)
	}
}

func TestQRHandler_MissingText(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "text") {
		t.Fatalf("expected error about 'text', got %s", rr.Body.String())
	}
}

func TestQRHandler_InvalidSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello&size=notanumber", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Size must be a valid number") {
		t.Fatalf("expected error about size, got %s", rr.Body.String())
	}
}

func TestQRHandler_SizeOutOfRange(t *testing.T) {
	cases := []string{"10", "2000"}
	for _, size := range cases {
		req := httptest.NewRequest("GET", "/qr?text=hello&size="+size, nil)
		rr := httptest.NewRecorder()
		qrHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for size %s, got %d", size, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Size must be between 50 and 1000") {
			t.Fatalf("expected error about size range, got %s", rr.Body.String())
		}
	}
}

func TestBarcodeHandler_Basic(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode?text=1234567890", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestBarcodeHandler_CustomSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode?text=1234567890&size=300", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestBarcodeHandler_Base64(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode?text=1234567890&base64=true", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/plain" {
		t.Fatalf("expected Content-Type text/plain, got %s", ct)
	}
	if _, err := base64.StdEncoding.DecodeString(rr.Body.String()); err != nil {
		t.Fatalf("expected valid base64, got error: %v", err)
	}
}

func TestBarcodeHandler_MissingText(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "text") {
		t.Fatalf("expected error about 'text', got %s", rr.Body.String())
	}
}

func TestBarcodeHandler_InvalidSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode?text=1234567890&size=notanumber", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Size must be a valid number") {
		t.Fatalf("expected error about size, got %s", rr.Body.String())
	}
}

func TestQRHandler_Cache(t *testing.T) {
	// First request should generate and cache
	req1 := httptest.NewRequest("GET", "/qr?text=testcache", nil)
	rr1 := httptest.NewRecorder()
	qrHandler(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr1.Code)
	}

	// Second request should use cache
	req2 := httptest.NewRequest("GET", "/qr?text=testcache", nil)
	rr2 := httptest.NewRecorder()
	qrHandler(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr2.Code)
	}

	// Verify both responses are identical
	if !bytes.Equal(rr1.Body.Bytes(), rr2.Body.Bytes()) {
		t.Fatal("cached response differs from original")
	}
}

func TestQRHandler_Cache_DifferentSizes(t *testing.T) {
	// Test that different sizes create different cache entries
	req1 := httptest.NewRequest("GET", "/qr?text=testcache&size=100", nil)
	rr1 := httptest.NewRecorder()
	qrHandler(rr1, req1)

	req2 := httptest.NewRequest("GET", "/qr?text=testcache&size=200", nil)
	rr2 := httptest.NewRecorder()
	qrHandler(rr2, req2)

	if bytes.Equal(rr1.Body.Bytes(), rr2.Body.Bytes()) {
		t.Fatal("different sizes should not have same cache entry")
	}
}

func TestRateLimiter(t *testing.T) {
	// Create a new rate limiter with 2 requests per second
	limiter := NewIPRateLimiter(2, 2)
	ip := "127.0.0.1"

	// Should allow first two requests
	if !limiter.Allow(ip) {
		t.Fatal("first request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Fatal("second request should be allowed")
	}

	// Third request should be rate limited
	if limiter.Allow(ip) {
		t.Fatal("third request should be rate limited")
	}

	// Wait for token refill
	time.Sleep(time.Second)

	// Should allow one more request
	if !limiter.Allow(ip) {
		t.Fatal("request after refill should be allowed")
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	limiter := NewIPRateLimiter(1, 1)
	ip1 := "127.0.0.1"
	ip2 := "127.0.0.2"

	// Both IPs should be able to make a request
	if !limiter.Allow(ip1) {
		t.Fatal("first IP should be allowed")
	}
	if !limiter.Allow(ip2) {
		t.Fatal("second IP should be allowed")
	}

	// Both IPs should be rate limited
	if limiter.Allow(ip1) {
		t.Fatal("first IP should be rate limited")
	}
	if limiter.Allow(ip2) {
		t.Fatal("second IP should be rate limited")
	}
}

func resetRateLimiter() {
	ipRateLimiter = NewIPRateLimiter(10, 20)
}

func TestQRHandler_RateLimit(t *testing.T) {
	resetRateLimiter() // Reset rate limiter before test

	// Create a request with a specific IP
	req := httptest.NewRequest("GET", "/qr?text=test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// Make requests up to the limit (20 requests to ensure we hit the bucket size)
	for i := 0; i < 20; i++ {
		rr := httptest.NewRecorder()
		qrHandler(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d should succeed, got status %d", i+1, rr.Code)
		}
	}

	// Next request should be rate limited
	rr := httptest.NewRecorder()
	qrHandler(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatal("request should be rate limited")
	}
}

func TestQRHandler_RateLimit_DifferentIPs(t *testing.T) {
	resetRateLimiter() // Reset rate limiter before test

	// Create two requests with different IPs
	req1 := httptest.NewRequest("GET", "/qr?text=test1", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	req2 := httptest.NewRequest("GET", "/qr?text=test2", nil)
	req2.RemoteAddr = "127.0.0.2:12345"

	// Both IPs should be able to make requests
	rr1 := httptest.NewRecorder()
	qrHandler(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatal("first IP should be allowed")
	}

	rr2 := httptest.NewRecorder()
	qrHandler(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatal("second IP should be allowed")
	}
}

func TestQRHandler_XForwardedFor(t *testing.T) {
	resetRateLimiter() // Reset rate limiter before test

	// Test rate limiting with X-Forwarded-For header
	req := httptest.NewRequest("GET", "/qr?text=test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	// Make requests up to the limit (20 requests to ensure we hit the bucket size)
	for i := 0; i < 20; i++ {
		rr := httptest.NewRecorder()
		qrHandler(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d should succeed, got status %d", i+1, rr.Code)
		}
	}

	// Next request should be rate limited
	rr := httptest.NewRecorder()
	qrHandler(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatal("request should be rate limited")
	}
}

// Tests for Shape Parameter

func TestQRHandler_Shape_Square(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello&shape=square", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestQRHandler_Shape_Rectangle(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello&shape=rectangle", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestQRHandler_Shape_Default(t *testing.T) {
	// Test that default shape is square
	req := httptest.NewRequest("GET", "/qr?text=hello", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
}

func TestQRHandler_Shape_Invalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello&shape=invalid", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Shape must be 'square' or 'rectangle'") {
		t.Fatalf("expected error about shape, got %s", rr.Body.String())
	}
}

func TestQRHandler_Type_QR(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello&type=qr", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
}

func TestQRHandler_Type_Barcode(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=1234567890&type=barcode", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
}

func TestQRHandler_Type_Default(t *testing.T) {
	// Test that default type is qr
	req := httptest.NewRequest("GET", "/qr?text=hello", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestQRHandler_Type_Invalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/qr?text=hello&type=invalid", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Type must be 'qr' or 'barcode'") {
		t.Fatalf("expected error about type, got %s", rr.Body.String())
	}
}

func TestQRHandler_Barcode_Rectangle(t *testing.T) {
	// Test barcode with rectangle shape via QR endpoint
	req := httptest.NewRequest("GET", "/qr?text=1234567890&type=barcode&shape=rectangle", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
}

func TestQRHandler_Barcode_Square(t *testing.T) {
	// Test barcode with square shape via QR endpoint
	req := httptest.NewRequest("GET", "/qr?text=1234567890&type=barcode&shape=square", nil)
	rr := httptest.NewRecorder()
	qrHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
}

func TestBarcodeHandler_Shape_Rectangle(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode?text=1234567890&shape=rectangle", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestBarcodeHandler_Shape_Square(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode?text=1234567890&shape=square", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	if len(rr.Body.Bytes()) < 100 {
		t.Fatalf("expected image data, got too few bytes (%d)", len(rr.Body.Bytes()))
	}
}

func TestBarcodeHandler_Shape_Default(t *testing.T) {
	// Test that default shape for barcode is rectangle
	req := httptest.NewRequest("GET", "/barcode?text=1234567890", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestBarcodeHandler_Shape_Invalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/barcode?text=1234567890&shape=invalid", nil)
	rr := httptest.NewRecorder()
	barcodeHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Shape must be 'square' or 'rectangle'") {
		t.Fatalf("expected error about shape, got %s", rr.Body.String())
	}
}

func TestQRHandler_Cache_DifferentShapes(t *testing.T) {
	// Test that different shapes create different cache entries
	req1 := httptest.NewRequest("GET", "/qr?text=testcache&shape=square", nil)
	rr1 := httptest.NewRecorder()
	qrHandler(rr1, req1)

	req2 := httptest.NewRequest("GET", "/qr?text=testcache&shape=rectangle", nil)
	rr2 := httptest.NewRecorder()
	qrHandler(rr2, req2)

	if bytes.Equal(rr1.Body.Bytes(), rr2.Body.Bytes()) {
		t.Fatal("different shapes should not have same cache entry")
	}
}

func TestQRHandler_Cache_DifferentTypes(t *testing.T) {
	// Test that different types create different cache entries
	req1 := httptest.NewRequest("GET", "/qr?text=1234567890&type=qr", nil)
	rr1 := httptest.NewRecorder()
	qrHandler(rr1, req1)

	req2 := httptest.NewRequest("GET", "/qr?text=1234567890&type=barcode", nil)
	rr2 := httptest.NewRecorder()
	qrHandler(rr2, req2)

	if bytes.Equal(rr1.Body.Bytes(), rr2.Body.Bytes()) {
		t.Fatal("different types should not have same cache entry")
	}
}

func TestQRHandler_Rectangle_DifferentFromSquare(t *testing.T) {
	// Verify that rectangle and square shapes produce different images
	req1 := httptest.NewRequest("GET", "/qr?text=hello&shape=square", nil)
	rr1 := httptest.NewRecorder()
	qrHandler(rr1, req1)

	req2 := httptest.NewRequest("GET", "/qr?text=hello&shape=rectangle", nil)
	rr2 := httptest.NewRecorder()
	qrHandler(rr2, req2)

	if rr1.Code != http.StatusOK || rr2.Code != http.StatusOK {
		t.Fatal("both requests should succeed")
	}

	if bytes.Equal(rr1.Body.Bytes(), rr2.Body.Bytes()) {
		t.Fatal("square and rectangle shapes should produce different images")
	}
}

func TestBarcodeHandler_Rectangle_DifferentFromSquare(t *testing.T) {
	// Verify that rectangle and square shapes produce different images for barcodes
	req1 := httptest.NewRequest("GET", "/barcode?text=1234567890&shape=square", nil)
	rr1 := httptest.NewRecorder()
	barcodeHandler(rr1, req1)

	req2 := httptest.NewRequest("GET", "/barcode?text=1234567890&shape=rectangle", nil)
	rr2 := httptest.NewRecorder()
	barcodeHandler(rr2, req2)

	if rr1.Code != http.StatusOK || rr2.Code != http.StatusOK {
		t.Fatal("both requests should succeed")
	}

	if bytes.Equal(rr1.Body.Bytes(), rr2.Body.Bytes()) {
		t.Fatal("square and rectangle shapes should produce different images")
	}
}
