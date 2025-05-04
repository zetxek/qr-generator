# QR Code Generator

A simple HTTP server that generates QR codes on demand. The server provides an endpoint that can generate QR codes with customizable size and format.

## Features

- Generate QR codes from text input
- Customize QR code size
- Option to receive QR code as PNG image or base64-encoded string
- Simple HTTP interface

## Running the Project

1. Make sure you have Go installed on your system (version 1.16 or higher recommended)

2. Clone the repository:
```bash
git clone <repository-url>
cd qr-generator
```

3. Install dependencies:
```bash
go mod tidy
```

4. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080.

## API Usage

### Health Check

```
GET /ping
```

Returns:
```
hello
```

### Generate QR Code

```
GET /qr?text=<text>&size=<size>&base64=<true|false>
```

Parameters:
- `text` (required): The text to encode in the QR code
- `size` (optional): Size of the QR code in pixels (default: 256, min: 50, max: 1000)
- `base64` (optional): Set to "true" to receive the QR code as a base64-encoded string

Examples:
- Basic usage: `http://localhost:8080/qr?text=HelloWorld`
- Custom size: `http://localhost:8080/qr?text=HelloWorld&size=500`
- Base64 output: `http://localhost:8080/qr?text=HelloWorld&base64=true`
- Custom size with base64: `http://localhost:8080/qr?text=HelloWorld&size=500&base64=true`

## Response

- When `base64=false` (default): Returns a PNG image
- When `base64=true`: Returns a base64-encoded string of the PNG image

## Error Handling

The server returns appropriate HTTP status codes and error messages for:
- Missing required parameters
- Invalid size values
- QR code generation failures








