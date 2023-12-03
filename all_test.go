package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/phprao/ColorOutput"
	"github.com/phprao/go-pdf/chromedp"
	"github.com/phprao/go-pdf/epub"
	"github.com/phprao/go-pdf/jpg"
	"github.com/phprao/go-pdf/pdfcpu"
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
	outputPath := "F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/output.pdf"
	inputPaths := []string{
		"F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/source/0b8e3474_big.jpg",
		"F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/source/0b972606_big.jpg",
		"F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/source/0bf817d2_big.jpg",
		"F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/source/0c3b1e96_big.jpg",
	}

	err := pdfcpu.ImagesToPdf(inputPaths, outputPath)
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

func TestMerge(t *testing.T) {
	pdffile := "F:/jx/20231114_4359/periodical/resource/epub/epub2/1158/1158-310040882/165_310073349_2.pdf"
	var htmls = []string{
		"F:/jx/20231114_4359/periodical/resource/epub/epub2/1158/1158-310040882/source/Chapter_1_1.pdf",
		"F:/jx/20231114_4359/periodical/resource/epub/epub2/1158/1158-310040882/source/Cover.pdf",
		"F:/jx/20231114_4359/periodical/resource/epub/epub2/1158/1158-310040882/source/Chapter_2_1.pdf",
	}
	if err := pdfcpu.MergePdf(htmls, pdffile); err != nil {
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
	jpg.ProcessTarGzFile("F:/jx/20231114_4359/periodical/resource/jpg/jpage3/41287/41287-358053/41287_358053.tar.gz")
	// unipdf.Run()
}

// go test -timeout 5m -v -run TestEPUB
func TestEPUB(t *testing.T) {
	epub.ProcessTarGzFile("F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/165_310073349.tar.gz")
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

func TestDiretory(t *testing.T) {
	str := `{"make-time":"2023-07-28 15:53:44","metadata":{"title":"\u57ce\u4e2d\u57ce\uff1a\u793e\u4f1a\u5b66\u5bb6\u7684\u8857\u5934\u53d1\u73b0","language":"zh","issueid":490521,"resourceType":3,"identifier":"urn:uuid:8745dd7b-ee41-45bf-894b-77e0cf545544","creator":"\uff08\u7f8e\uff09\u6587\u5361\u7279\u65af\uff08Venkatesh, S.\uff09","date":"2015-12-11","img_attr":1},"images":["OEBPS\/Images\/cover_small.jpg","OEBPS\/Images\/cover.jpg","OEBPS\/Images\/logo.jpg"],"css":"OEBPS\/Styles\/sgc-toc.css","js":"","catalog":[{"playOrder":1,"title":"\u5c01\u9762","src":"OEBPS\/Text\/Cover.xhtml","id":200000000,"lev":0,"is_cat":0,"titleId":-1},{"playOrder":2,"title":"\u5e8f","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0001.xhtml","id":1,"lev":0,"is_cat":0,"catalog_id":0,"titleId":1},{"playOrder":3,"title":"\u524d\u8a00","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0002.xhtml","id":2,"lev":0,"is_cat":0,"catalog_id":0,"titleId":2},{"playOrder":4,"title":"\u7b2c\u4e00\u7ae0 \u4f5c\u4e3a\u7a77\u56f0\u9ed1\u4eba\u7684\u611f\u89c9\u600e\u4e48\u6837\uff1f","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0003.xhtml","id":3,"lev":0,"is_cat":0,"catalog_id":0,"titleId":3},{"playOrder":5,"title":"\u7b2c\u4e8c\u7ae0 \u8054\u90a6\u8857\u7684\u6700\u521d\u65f6\u5149","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0004.xhtml","id":4,"lev":0,"is_cat":0,"catalog_id":0,"titleId":4},{"playOrder":6,"title":"\u7b2c\u4e09\u7ae0 \u8c01\u6765\u7f69\u7740\u6211","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0005.xhtml","id":5,"lev":0,"is_cat":0,"catalog_id":0,"titleId":5},{"playOrder":7,"title":"\u7b2c\u56db\u7ae0 \u9ed1\u5e2e\u8001\u5927\u7684\u4e00\u5929","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0006.xhtml","id":6,"lev":0,"is_cat":0,"catalog_id":0,"titleId":6},{"playOrder":8,"title":"\u7b2c\u4e94\u7ae0 \u8d1d\u5229\u5973\u58eb\u7684\u8857\u533a","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0007.xhtml","id":7,"lev":0,"is_cat":0,"catalog_id":0,"titleId":7},{"playOrder":9,"title":"\u7b2c\u516d\u7ae0 \u6df7\u6df7\u513f\u4e0e\u6df7\u8ff9","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0008.xhtml","id":8,"lev":0,"is_cat":0,"catalog_id":0,"titleId":8},{"playOrder":10,"title":"\u7b2c\u4e03\u7ae0 \u904d\u4f53\u9cde\u4f24","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0009.xhtml","id":9,"lev":0,"is_cat":0,"catalog_id":0,"titleId":9},{"playOrder":11,"title":"\u7b2c\u516b\u7ae0 \u56e2\u7ed3\u7684\u5e2e\u6d3e","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0010.xhtml","id":10,"lev":0,"is_cat":0,"catalog_id":0,"titleId":10},{"playOrder":12,"title":"\u4f5c\u8005\u58f0\u660e","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0011.xhtml","id":11,"lev":0,"is_cat":0,"catalog_id":0,"titleId":11},{"playOrder":13,"title":"\u81f4\u8c22","page":0,"author":"","length":0,"src":"OEBPS\/Text\/Section0001_0012.xhtml","id":12,"lev":0,"is_cat":0,"catalog_id":0,"titleId":12}],"spine":[{"id":-1,"length":0,"src":"OEBPS\/Text\/Cover.xhtml"},{"id":0,"length":0,"src":"OEBPS\/Text\/bq.xhtml"},{"id":0,"length":0,"src":"OEBPS\/Text\/TOC.xhtml"},{"id":0,"length":0,"src":"OEBPS\/Text\/Section0001.xhtml"},{"id":1,"length":2025,"src":"OEBPS\/Text\/Section0001_0001.xhtml"},{"id":2,"length":1573,"src":"OEBPS\/Text\/Section0001_0002.xhtml"},{"id":3,"length":25895,"src":"OEBPS\/Text\/Section0001_0003.xhtml"},{"id":4,"length":40659,"src":"OEBPS\/Text\/Section0001_0004.xhtml"},{"id":5,"length":48900,"src":"OEBPS\/Text\/Section0001_0005.xhtml"},{"id":6,"length":35072,"src":"OEBPS\/Text\/Section0001_0006.xhtml"},{"id":7,"length":44111,"src":"OEBPS\/Text\/Section0001_0007.xhtml"},{"id":8,"length":36218,"src":"OEBPS\/Text\/Section0001_0008.xhtml"},{"id":9,"length":30071,"src":"OEBPS\/Text\/Section0001_0009.xhtml"},{"id":10,"length":36533,"src":"OEBPS\/Text\/Section0001_0010.xhtml"},{"id":11,"length":834,"src":"OEBPS\/Text\/Section0001_0011.xhtml"},{"id":12,"length":4135,"src":"OEBPS\/Text\/Section0001_0012.xhtml"},{"id":0,"length":0,"src":"OEBPS\/Text\/Section0002.xhtml"}]}`

	var dic util.Directory
	json.Unmarshal([]byte(str), &dic)

	util.PrettyPrint(dic)
}
