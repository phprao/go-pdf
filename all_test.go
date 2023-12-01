package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/phprao/ColorOutput"
	"github.com/phprao/go-pdf/chromedp"
	"github.com/phprao/go-pdf/unipdf"
	"github.com/phprao/go-pdf/util"
)

func TestHash(t *testing.T) {
	h := util.NewHash()
	var iid int64 = 357242
	var rid int64 = 41061
	ok, data := h.GetHash(iid, rid, 1, 1000)
	if ok == 0 {
		return
	}

	hs := []string{
		"0b8e3474",
		"0b972606",
		"0bbefa12",
		"0c3b1e96",
		"0e4abc46",
		"01d3b323",
	}

	for _, h := range hs {
		for k := range data {
			if h == data[k].Hash {
				fmt.Println(h, data[k].Page)
			}
		}
	}
}

func TestImagesToPDF(t *testing.T) {
	outputPath := "unipdf/output.pdf"
	inputPaths := []string{
		"unipdf/images/0b8e3474_big.jpg",
		"unipdf/images/0b972606_big.jpg",
		"unipdf/images/0bbefa12_big.jpg",
		"unipdf/images/0c3b1e96_big.jpg",
		"unipdf/images/0e4abc46_big.jpg",
		"unipdf/images/01d3b323_big.jpg",
	}

	err := unipdf.ImagesToPdf(inputPaths, outputPath)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Complete, see output file: %s\n", outputPath)
}

func TestFile(t *testing.T) {
	filename := "F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/41287_358053.tar.gz"
	fmt.Println(filename[:strings.LastIndex(filename, "/")]) // F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053
}

func TestSort(t *testing.T) {
	var htmls = []string{
		"OEBPS/Text/Chapter_3_1.xhtml",
		"OEBPS/Text/Chapter_2_8.xhtml",
		"OEBPS/Text/Chapter_2_2.xhtml",
		"OEBPS/Text/Chapter_1_1.xhtml",
		"OEBPS/Text/Chapter_2_10.xhtml",
		"OEBPS/Text/Chapter_5_1.xhtml",
		"OEBPS/Text/Chapter_5_2.xhtml",
		"OEBPS/Text/Chapter_3_2.xhtml",
		"OEBPS/Text/Cover.xhtml",
		"OEBPS/Text/Chapter_2_5.xhtml",
	}
	sort.Stable(chromedp.SortStringSliceIncrement(htmls))

	fmt.Println(htmls)
}

func TestMerge(t *testing.T) {
	pdffile := "F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/165_310073349_2.pdf"
	var htmls = []string{
		"F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/source/Chapter_1_1.pdf",
		"F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/source/Cover.pdf",
		"F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/source/Chapter_2_5.pdf",
	}
	if err := unipdf.MergePdf(htmls, pdffile); err != nil {
		fmt.Println(err)
	}
}

func TestHtmlToPdf(t *testing.T) {
	sourceDir := "F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/source"
	html := "OEBPS/Text/Cover.xhtml"
	outputPdf := "F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/Cover2.pdf"
	chromedp.HtmlToPdf(sourceDir, html, outputPdf)
}

// go test -timeout 5m -v -run TestJPG
func TestJPG(t *testing.T) {
	unipdf.ProcessTarGzFile("F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/41287_358053.tar.gz")
	// unipdf.Run()
}

// go test -timeout 5m -v -run TestEPUB
func TestEPUB(t *testing.T) {
	chromedp.ProcessTarGzFile("F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/165_310073349.tar.gz")
	// chromedp.Run()
}

func TestColor(t *testing.T) {
	ColorOutput.Colorful.WithFrontColor("red").Println("test")
}

func TestFind(t *testing.T) {
	tarGzFiles, err := util.FindTarGzFile("F:/jx/20231114_4359/periodical/resource/epub")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range tarGzFiles {
		fmt.Println(v)
	}
}