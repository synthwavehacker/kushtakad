package helpers

import (
	"os"
)

// exists returns whether the given file or directory exists or not
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}

func CreateImages(prefix string, path string, ext string, src string, crop bool) error {

	// create logos directory if not found
	res, err := FileExists(path)
	if err != nil || res == false {
		err := os.MkdirAll(path, 0744)
		if err != nil {
			return err
		}
	}

	var smallWidth, smallHeight, mediumWidth, mediumHeight, largeWidth, largeHeight, origWidth, origHeight uint
	smallWidth, smallHeight = 100, 100
	mediumWidth, mediumHeight = 300, 300
	largeWidth, largeHeight = 500, 500
	origWidth, origHeight = 1024, 1024

	small_name := "/" + prefix + "_small." + ext
	small_path := path + small_name

	medium_name := "/" + prefix + "_medium." + ext
	medium_path := path + medium_name

	large_name := "/" + prefix + "_large." + ext
	large_path := path + large_name

	orig_name := "/" + prefix + "_orig." + ext
	orig_path := path + orig_name

	switch ext {
	case "jpg":

		// small
		err := FormatJpg(smallWidth, smallHeight, src, small_path, crop)
		if err != nil {
			return err
		}

		// medium
		err = FormatJpg(mediumWidth, mediumHeight, src, medium_path, crop)
		if err != nil {
			return err
		}

		// large
		err = FormatJpg(largeWidth, largeHeight, src, large_path, crop)
		if err != nil {
			return err
		}

		// orig
		err = FormatJpg(origWidth, origHeight, src, orig_path, crop)
		if err != nil {
			return err
		}

	case "png":
		// small
		err := FormatPng(smallWidth, smallHeight, src, small_path, crop)
		if err != nil {
			return err
		}

		// medium
		err = FormatPng(mediumWidth, mediumHeight, src, medium_path, crop)
		if err != nil {
			return err
		}

		// large
		err = FormatPng(largeWidth, largeHeight, src, large_path, crop)
		if err != nil {
			return err
		}

		// orig
		err = FormatPng(origWidth, origHeight, src, orig_path, crop)
		if err != nil {
			return err
		}

	case "gif":
		// small
		err := FormatGif(smallWidth, smallHeight, src, small_path, crop)
		if err != nil {
			return err
		}

		// medium
		err = FormatGif(mediumWidth, mediumHeight, src, medium_path, crop)
		if err != nil {
			return err
		}

		//large
		err = FormatGif(largeWidth, largeHeight, src, large_path, crop)
		if err != nil {
			return err
		}

		// orig
		err = FormatGif(origWidth, origHeight, src, orig_path, crop)
		if err != nil {
			return err
		}
	}
	return nil
}
