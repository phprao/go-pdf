package gofpdf

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
)

type chapter struct {
	FileName string `xml:"file-name,attr"`
	Content  string `xml:",innerxml"`
}

func Run() {
	// 读取EPUB文件
	file, _ := os.Open("F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/165_310073349/310073349_d81add70.epub")
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// 解压缩EPUB文件
	r, _ := zip.NewReader(file, info.Size())
	for _, f := range r.File {
		// 检查文件类型
		if f.Name[len(f.Name)-6:] == ".xhtml" {
			// 读取HTML文件中的内容
			htmlFile, _ := f.Open()
			defer htmlFile.Close()
			htmlContent, _ := io.ReadAll(htmlFile)

			// 解析HTML内容
			var c chapter
			xml.Unmarshal(htmlContent, &c)

			// 将HTML内容转换为PDF格式
			pdf := gofpdf.New("P", "mm", "A4", "")
			pdf.AddPage()
			pdf.Write(5, c.Content)
			pdf.OutputFileAndClose(fmt.Sprintf("%s.pdf", c.FileName))
		}
	}
}

func Download() {
	// 读取EPUB文件
	file, _ := os.Open("F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/165_310073349/310073349_d81add70.epub")
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// 解压缩EPUB文件
	r, _ := zip.NewReader(file, info.Size())
	for _, f := range r.File {

		// fmt.Println(f.Name)

		// if f.Name == "OEBPS/Styles/Style.css" {
		htmlFile, _ := f.Open()
		defer htmlFile.Close()

		newFile, _ := os.OpenFile("resource/310073349_d81add70/"+f.Name, os.O_CREATE|os.O_RDWR, 0777)
		defer newFile.Close()
		_, err = io.Copy(newFile, htmlFile)
		if err != nil {
			panic(err)
		}
		// }

	}
}
