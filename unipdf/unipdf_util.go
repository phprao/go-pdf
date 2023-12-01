package unipdf

/*

go get github.com/unidoc/unipdf/v3

*/

import (
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

var UNIDOC_LICENSE_API_KEY = "e1c7b4aa760dfe992d9a396a358acca0ecd0f7e278fd0ab536b6791a162c5856"

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(UNIDOC_LICENSE_API_KEY)
	if err != nil {
		panic(err)
	}
}

// 将一个文件夹中的图片合并成一个PDF文件
// Images to PDF.
func ImagesToPdf(inputPaths []string, outputPath string) error {
	c := creator.New()

	for _, imgPath := range inputPaths {
		if imgPath == "" {
			continue
		}
		common.Log.Debug("Image: %s", imgPath)

		img, err := c.NewImageFromFile(imgPath)
		if err != nil {
			common.Log.Debug("Error loading image: %v", err)
			return err
		}
		img.ScaleToWidth(612.0)

		// Use page width of 612 points, and calculate the height proportionally based on the image.
		// Standard PPI is 72 points per inch, thus a width of 8.5"
		height := 612.0 * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{612, height})
		c.NewPage()
		img.SetPos(0, 0)
		_ = c.Draw(img)
	}

	err := c.WriteToFile(outputPath)
	return err
}

func MergePdf(inputPaths []string, outputPath string) error {
	pdfWriter := model.NewPdfWriter()

	for _, inputPath := range inputPaths {
		pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
		if err != nil {
			return err
		}
		defer f.Close()

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
