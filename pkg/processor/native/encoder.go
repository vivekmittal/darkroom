package native

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/gojek/darkroom/pkg/processor"
	"image"
	"image/jpeg"
	"image/png"
)

const (
	KiloBytes = 1024
	MegaBytes = 1024 * KiloBytes
)

var compressionConfigs = []struct {
	sizeLevel int
	quality   int
}{
	{sizeLevel: 500 * KiloBytes, quality: 75},
	{sizeLevel: 1 * MegaBytes, quality: 75},
	{sizeLevel: 5 * MegaBytes, quality: 75},
	{sizeLevel: 10 * MegaBytes, quality: 75},
	{sizeLevel: 100 * MegaBytes, quality: 75},
	{sizeLevel: 1000 * MegaBytes, quality: 75},
}

// Encoder is an interface to Encode image and return the encoded byte array or error
type Encoder interface {
	Encode(img image.Image) ([]byte, error)
	EncodeWithSize(img image.Image, size int) ([]byte, error)
}

// JpegEncoder is an object to encode image to byte array with jpeg format
type JpegEncoder struct {
	Option *jpeg.Options
}

// PngEncoder is an object to encode image to byte array with png format
type PngEncoder struct {
	Encoder *png.Encoder
}

// WebPEncoder is an object to encode image to byte array with webp format
type WebPEncoder struct {
	Option *webp.Options
}

// NopEncoder is a no-op encoder object for unsupported format and will return error
type NopEncoder struct{}

func (e *PngEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := e.Encoder.Encode(buff, img)
	return buff.Bytes(), err
}

func (e *PngEncoder) EncodeWithSize(img image.Image, size int) ([]byte, error) {
	return e.Encode(img)
}

func (e *JpegEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := jpeg.Encode(buff, img, e.Option)
	return buff.Bytes(), err
}

func (e *JpegEncoder) EncodeWithSize(img image.Image, size int) ([]byte, error) {
	buff := &bytes.Buffer{}

	quality := e.Option.Quality
	for _, c := range compressionConfigs {
		if size <= c.sizeLevel {
			quality = c.quality
			break
		}
	}

	fmt.Println("Compressing with: ", quality)
	err := jpeg.Encode(buff, img, &jpeg.Options{Quality: quality})
	return buff.Bytes(), err
}

func (e *WebPEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := webp.Encode(buff, img, e.Option)
	return buff.Bytes(), err
}

func (e *WebPEncoder) EncodeWithSize(img image.Image, size int) ([]byte, error) {
	return e.Encode(img)
}

func (e *NopEncoder) Encode(img image.Image) ([]byte, error) {
	return nil, errors.New("unknown format: failed to encode image")
}

func (e *NopEncoder) EncodeWithSize(img image.Image, size int) ([]byte, error) {
	return nil, errors.New("unknown format: failed to encode image")
}

// Encoders is a struct to store all supported encoders so that we don't have to create new encoder every time
type Encoders struct {
	jpegEncoder *JpegEncoder
	pngEncoder  *PngEncoder
	noOpEncoder *NopEncoder
	webPEncoder *WebPEncoder
}

// EncodersOption represents builder function for Encoders
type EncodersOption func(*Encoders)

// GetEncoder takes an input of image and extension and return the appropriate Encoder for encoding the image
func (e *Encoders) GetEncoder(img image.Image, ext string) Encoder {
	switch ext {
	case processor.ExtensionJPG, processor.ExtensionJPEG:
		return e.jpegEncoder
	case processor.ExtensionPNG:
		if e.jpegEncoder.Option.Quality != 100 && isOpaque(img) {
			return e.jpegEncoder
		}
		return e.pngEncoder
	case processor.ExtensionWebP:
		return e.webPEncoder
	default:
		return e.noOpEncoder
	}
}

// WithJpegEncoder is a builder function for setting custom JpegEncoder
func WithJpegEncoder(jpegEncoder *JpegEncoder) EncodersOption {
	return func(e *Encoders) {
		e.jpegEncoder = jpegEncoder
	}
}

// WithPngEncoder is a builder function for setting custom PngEncoder
func WithPngEncoder(pngEncoder *PngEncoder) EncodersOption {
	return func(e *Encoders) {
		e.pngEncoder = pngEncoder
	}
}

// WithWebPEncoder is a builder function for setting custom WebPEncoder
func WithWebPEncoder(webPEncoder *WebPEncoder) EncodersOption {
	return func(e *Encoders) {
		e.webPEncoder = webPEncoder
	}
}

// NewEncoders creates a new Encoders, if called without parameter (builder), all encoders option will be default
func NewEncoders(opts ...EncodersOption) *Encoders {
	e := &Encoders{
		jpegEncoder: &JpegEncoder{Option: &jpeg.Options{Quality: jpeg.DefaultQuality}},
		pngEncoder: &PngEncoder{
			Encoder: &png.Encoder{CompressionLevel: png.BestCompression},
		},
		noOpEncoder: &NopEncoder{},
		webPEncoder: &WebPEncoder{},
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}
