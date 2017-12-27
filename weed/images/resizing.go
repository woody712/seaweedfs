package images

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/disintegration/imaging"
    "github.com/lazywei/go-opencv/opencv"
)

func Resized(ext string, data []byte, width, height int, mode string) (resized []byte, w int, h int) {
	if width == 0 && height == 0 {
		return data, 0, 0
	}
	srcImage, _, err := image.Decode(bytes.NewReader(data))
	if err == nil {
		bounds := srcImage.Bounds()
		var dstImage *image.NRGBA
		if bounds.Dx() > width && width != 0 || bounds.Dy() > height && height != 0 {
			switch mode {
			case "fit":
				dstImage = imaging.Fit(srcImage, width, height, imaging.Lanczos)
			case "fill":
				dstImage = imaging.Fill(srcImage, width, height, imaging.Center, imaging.Lanczos)
			default:
				if width == height && bounds.Dx() != bounds.Dy() {
					dstImage = imaging.Thumbnail(srcImage, width, height, imaging.Lanczos)
					w, h = width, height
				} else {
					dstImage = imaging.Resize(srcImage, width, height, imaging.Lanczos)
				}
			}
		} else {
			return data, bounds.Dx(), bounds.Dy()
		}
		var buf bytes.Buffer
		switch ext {
		case ".png":
			png.Encode(&buf, dstImage)
		case ".jpg", ".jpeg":
			jpeg.Encode(&buf, dstImage, nil)
		case ".gif":
			gif.Encode(&buf, dstImage, nil)
		}
		return buf.Bytes(), dstImage.Bounds().Dx(), dstImage.Bounds().Dy()
	} else {
		glog.Error(err)
	}
	return data, 0, 0
}

func ResizedByCv(ext string, data []byte, width, height int, mode string) (resized []byte, w int, h int) {
	if width == 0 && height == 0 {
		return data, 0, 0
	}
	srcImage := opencv.DecodeImageMem(data)

	if srcImage.Width() > width && width != 0 || srcImage.Height() > height && height != 0 {
		switch mode {
		case "fill":
			w, h := width, height
		case "fit":
		default:
			if srcImage.Width() / width > srcImage.Height() / height {
				w, h := width, srcImage.Height() / srcImage.Width() * width
			} else {
				w, h := srcImage.Width() / srcImage.Height() * height, height
			}
		}
		newImg := opencv.Resize(srcImage, width, height, 0)
	} else {
		newImg := srcImage
	}
	params := []int{opencv.CV_IMWRITE_JPEG_QUALITY, 95, 0}
	return opencv.EncodeImage(ext, unsafe.Pointer(newImg), params).GetData(), 0, 0
}
