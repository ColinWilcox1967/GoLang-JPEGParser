package jpegparser

import (
	"errors"
	"os"
)

// JPEGHeader represents the structure of a JPEG file header.
type JPEGHeader struct {
	SOIMarker   [2]byte // Start of Image marker
	APP0Marker  [2]byte // Application marker
	Length      uint16  // Length of the APP0 segment
	Identifier  string  // Identifier, typically "JFIF\0"
	Version     [2]byte // Version, e.g., 1.01
	Units       byte    // Units for X and Y density
	XDensity    uint16  // Horizontal pixel density
	YDensity    uint16  // Vertical pixel density
	XThumbnail  byte    // Thumbnail width
	YThumbnail  byte    // Thumbnail height
}

// ParseJPEGHeader parses the byte slice representing a JPEG header and interprets each field.
func ParseJPEGHeader(data []byte) (*JPEGHeader, error) {
	if len(data) < 20 {
		return nil, errors.New("data too short to contain a valid JPEG header")
	}

	header := &JPEGHeader{
		SOIMarker:  [2]byte{data[0], data[1]},
		APP0Marker: [2]byte{data[2], data[3]},
		Length:     uint16(data[4])<<8 | uint16(data[5]),
		Identifier: string(data[6:11]),
		Version:    [2]byte{data[11], data[12]},
		Units:      data[13],
		XDensity:   uint16(data[14])<<8 | uint16(data[15]),
		YDensity:   uint16(data[16])<<8 | uint16(data[17]),
		XThumbnail: data[18],
		YThumbnail: data[19],
	}

	// Validate SOI Marker
	if header.SOIMarker != [2]byte{0xFF, 0xD8} {
		return nil, errors.New("invalid SOI marker")
	}

	// Validate APP0 Marker
	if header.APP0Marker != [2]byte{0xFF, 0xE0} {
		return nil, errors.New("invalid APP0 marker")
	}

	// Validate Identifier
	if header.Identifier != "JFIF\u0000" {
		return nil, errors.New("invalid identifier, expected 'JFIF\\0'")
	}

	return header, nil
}

// WriteJPEGFile writes the parsed header and raw data into a file.
func WriteJPEGFile(header *JPEGHeader, rawData []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write SOI Marker
	if _, err := file.Write(header.SOIMarker[:]); err != nil {
		return err
	}

	// Write APP0 Marker
	if _, err := file.Write(header.APP0Marker[:]); err != nil {
		return err
	}

	// Write Length
	length := []byte{byte(header.Length >> 8), byte(header.Length & 0xFF)}
	if _, err := file.Write(length); err != nil {
		return err
	}

	// Write Identifier
	if _, err := file.Write([]byte(header.Identifier)); err != nil {
		return err
	}

	// Write Version
	if _, err := file.Write(header.Version[:]); err != nil {
		return err
	}

	// Write Units
	if _, err := file.Write([]byte{header.Units}); err != nil {
		return err
	}

	// Write Densities
	density := []byte{
		byte(header.XDensity >> 8), byte(header.XDensity & 0xFF),
		byte(header.YDensity >> 8), byte(header.YDensity & 0xFF),
	}
	if _, err := file.Write(density); err != nil {
		return err
	}

	// Write Thumbnail Dimensions
	if _, err := file.Write([]byte{header.XThumbnail, header.YThumbnail}); err != nil {
		return err
	}

	// Write Remaining Data
	if _, err := file.Write(rawData); err != nil {
		return err
	}

	return nil
}

