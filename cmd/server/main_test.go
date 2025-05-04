package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
