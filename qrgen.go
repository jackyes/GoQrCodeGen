package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

func generateQRCode(data string, size int) (image.Image, error) {
	qr, err := qrcode.New(data, qrcode.High)
	if err != nil {
		return nil, err
	}
	return qr.Image(size), nil
}

func decodeImage(file io.Reader) (image.Image, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	switch format {
	case "jpeg", "png":
		return img, nil
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
}

func overlayImageOnQRCode(qrCode image.Image, overlay image.Image, overlayPercent float64) (image.Image, error) {
	return overlayImageOnQRCodeWithOpacity(qrCode, overlay, overlayPercent, 1)
}

func overlayImageOnQRCodeWithOpacity(qrCode image.Image, overlay image.Image, overlayPercent, overlayOpacity float64) (image.Image, error) {
	qrBounds := qrCode.Bounds()
	qrWidth := qrBounds.Dx()
	qrHeight := qrBounds.Dy()

	overlayMaxWidth := int(float64(qrWidth) * overlayPercent)
	overlayMaxHeight := int(float64(qrHeight) * overlayPercent)

	overlay = resize.Thumbnail(uint(overlayMaxWidth), uint(overlayMaxHeight), overlay, resize.Lanczos3)

	offset := image.Pt((qrWidth-overlay.Bounds().Dx())/2, (qrHeight-overlay.Bounds().Dy())/2)

	b := qrBounds
	m := image.NewRGBA(b)

	draw.Draw(m, qrBounds, qrCode, image.Point{}, draw.Src)

	overlay = applyOpacity(overlay, overlayOpacity)

	draw.Draw(m, overlay.Bounds().Add(offset), overlay, image.Point{}, draw.Over)

	return m, nil
}

func applyOpacity(img image.Image, opacity float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			a = uint32(float64(a) * opacity)
			newImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	return newImg
}
