# qr-generator

# QR Code Generator

A simple web service that generates QR codes and gradient images.

## Features

- Generate QR codes from text
- Customize QR code size
- Generate colorful gradient images

## API Endpoints

### Generate QR Code

Endpoint: `/qr`
Method: `GET`

Parameters:
- `text`: The text to encode in the QR code (required)
- `size`: The size of the QR code (optional, default is 256)
- `base64`: Whether to return the image as a base64 encoded string (optional, default is false)

Example:
```
/qr?text=https://example.com&size=512&base64=true
``` 




