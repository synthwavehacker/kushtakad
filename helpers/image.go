package helpers

import (
	//stdlib

	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	return ValidateImage(buffer)
}

func ValidateImage(b []byte) (string, error) {

	m := http.DetectContentType(b)
	log.Debug("MIMETYPE: ", m)
	switch m {
	case "image/jpeg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "image/gif":
		return "gif", nil
	}
	return "", errors.New("This isn't a valid JPG, PNG, or GIF")
}

func ValidateMimeType(b []byte) string {

	m := http.DetectContentType(b)
	log.Debug("MIMETYPE: ", m)
	switch m {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	}
	return ""
}

func ValidateMimeTypeString(s string) string {
	mp4 := regexp.MustCompile(`\.mp4`)
	webm := regexp.MustCompile(`\.webm`)
	mov := regexp.MustCompile(`\.mov`)

	switch {
	case mp4.MatchString(s):
		return "mp4"
	case webm.MatchString(s):
		return "webm"
	case mov.MatchString(s):
		return "webm"
	}
	return ""
}

// this is broken for mp4
func ValidateVideoMimeType(b []byte) string {

	m := http.DetectContentType(b)
	switch m {
	case "video/mp4":
		return "mp4"
	}
	return ""
}

func WidthHeight(img image.Image, maxWidth uint, maxHeight uint) (uint, uint) {

	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	newWidth, newHeight := origWidth, origHeight

	log.Debug("NW:NH ", newWidth, newHeight)
	if newWidth > newHeight {
		diff := newHeight - maxHeight
		return newWidth - diff, maxHeight
	} else if newHeight > newWidth {
		diff := newWidth - maxWidth
		return maxWidth, newHeight - diff
	}
	return maxWidth, maxHeight
}

func ResizeImage(img image.Image, maxWidth uint, maxHeight uint) image.Image {

	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	newWidth, newHeight := origWidth, origHeight

	// Preserve aspect ratio
	if origWidth > maxWidth {
		newHeight = uint(origHeight * maxWidth / origWidth)
		if newHeight < 1 {
			newHeight = 1
		}
		newWidth = maxWidth
	}

	if newHeight > maxHeight {
		newWidth = uint(newWidth * maxHeight / newHeight)
		if newWidth < 1 {
			newWidth = 1
		}
		newHeight = maxHeight
	}

	//return resize.Resize(newWidth, newHeight, img, resize.Bicubic)
	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	//return resize.Resize(1000, 0, img, resize.Lanczos3)
	//return resize.Resize(newWidth, newHeight, img, resize.NearestNeighbor)
}

func ImageToPaletted(img image.Image, pal color.Palette) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, pal)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}

func CropImageFromCenter(img image.Image) (image.Image, error) {

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:   1,
		Height:  1,
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})

	if err != nil {
		return nil, err
	}
	return croppedImg, nil
}

func FormatJpg(width uint, height uint, src string, full_path string, crop bool) error {

	var croppedImg image.Image
	var m image.Image

	ior, err := os.Open(src)
	if err != nil {
		return err
	}
	defer ior.Close()

	img, err := jpeg.Decode(ior)
	if err != nil {
		return err
	}

	if crop == true {
		croppedImg, err = CropImageFromCenter(img)
		if err != nil {
			return err
		}
		m = ResizeImage(croppedImg, width, height)
	} else {
		m = ResizeImage(img, width, height)
	}

	dst, err := os.Create(full_path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// write new image to file
	jpeg.Encode(dst, m, nil)
	return nil
}

func FormatGif(width uint, height uint, src string, full_path string, crop bool) error {
	ior, err := os.Open(src)
	if err != nil {
		return err
	}
	defer ior.Close()

	t0 := time.Now()
	im, err := gif.DecodeAll(ior)
	if err != nil {
		return err
	}

	bounds := im.Image[0].Bounds()
	img := image.NewRGBA(bounds)
	pal := im.Image[0].Palette

	// Resize each frame.
	for index := range im.Image {
		err := FormatAndAssign(width, height, im.Image, pal, img, index, crop)
		if err != nil {
			return err
		}
	}
	dst, err := os.Create(full_path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// write new image to file
	gif.EncodeAll(dst, im)

	//optimize using gifsicle
	//out, err := OptimizeGif(dst)
	//if err != nil || len(out) > 0 { return err }

	t1 := time.Now()
	log.Debugf("The call took %v to run.\n", t1.Sub(t0))
	return nil
}

func FormatPng(width uint, height uint, src string, full_path string, crop bool) error {

	var croppedImg image.Image
	var m image.Image

	ior, err := os.Open(src)
	if err != nil {
		return err
	}
	defer ior.Close()

	img, err := png.Decode(ior)
	if err != nil {
		return err
	}

	if crop == true {
		croppedImg, err = CropImageFromCenter(img)
		if err != nil {
			return err
		}
		m = ResizeImage(croppedImg, width, height)
	} else {
		m = ResizeImage(img, width, height)
	}

	dst, err := os.Create(full_path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// write new image to file
	png.Encode(dst, m)
	return nil
}

func FormatAndAssign(width uint, height uint, im []*image.Paletted, pal color.Palette, img *image.RGBA, index int, crop bool) error {
	var err error
	im[index], err = FormatFrames(width, height, im[index], pal, img, crop)
	if err != nil {
		return err
	}
	return nil
}

func FormatFrames(width uint, height uint, frame *image.Paletted, pal color.Palette, img *image.RGBA, crop bool) (*image.Paletted, error) {

	var croppedImg image.Image
	var m image.Image
	var err error

	bounds := frame.Bounds()
	draw.Draw(img, bounds, frame, bounds.Min, draw.Over)

	if crop == true {
		croppedImg, err = CropImageFromCenter(img)
		if err != nil {
			return nil, err
		}
		m = ResizeImage(croppedImg, width, height)
	} else {
		m = ResizeImage(img, width, height)
	}

	return ImageToPaletted(m, pal), nil
}

func ImageContainer(rect image.Rectangle) draw.Image {
	b := image.Rect(0, 0, rect.Dx(), rect.Dy())
	//return image.NewRGBA(b)
	return image.NewAlpha16(b)
}
