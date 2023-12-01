package pdfcpu

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

/*

go get github.com/pdfcpu/pdfcpu/...

常用的包
github.com/pdfcpu/pdfcpu/pkg/api
github.com/pdfcpu/pdfcpu/pkg/pdfcpu

在 pkg/api/example_test.go 找示例

*/

func ImagesToPdf(inputPaths []string, outputPath string) error {
	api.ImportImagesFile(inputPaths, outputPath, nil, nil)
	return nil
}

func MergePdf(inputPaths []string, outputPath string) error {
	api.MergeCreateFile(inputPaths, outputPath, nil)
	return nil
}
