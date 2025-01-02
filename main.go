package main

import (
	"./jpegparser"
)

func main () {

	// Example usage:
 	jpegHeader := []byte{
 		0xFF, 0xD8, // SOI marker
 		0xFF, 0xE0, // APP0 marker
 		0x00, 0x10, // Length
 		0x4A, 0x46, 0x49, 0x46, 0x00, // "JFIF\0"
 		0x01, 0x01, // Version
 		0x00,       // Units
 		0x00, 0x48, // XDensity
 		0x00, 0x48, // YDensity
 		0x00,       // XThumbnail
 		0x00,       // YThumbnail
 	}

 	header, err := jpegparser.ParseJPEGHeader(jpegHeader)
 	if err != nil {
 		panic(err)
 	}

 	rawData := []byte{} // Additional image data if available
 	filename := "output.jpg"
 	if err := jpegparser.WriteJPEGFile(header, rawData, filename); err != nil {
 		panic(err)
 	}
 	println("JPEG file written successfully")
 }

 // end of file
