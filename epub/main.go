package epub

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/phprao/go-pdf/chromedp"
	"github.com/phprao/go-pdf/pdfcpu"
	"github.com/phprao/go-pdf/util"
)

// 完整路径：F:/jx/20231114_4359/periodical/resource/epub/epub2/165/165-310073349/165_310073349/310073349_d81add70.epub
// root：F:/jx/20231114_4359/periodical/resource
func Run(root string, ch chan util.Msg) {
	dir := root + "/epub"

	tarGzFiles, err := util.FindTarGzFile(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range tarGzFiles {
		if err := ProcessTarGzFile(v); err != nil {
			ch <- util.Msg{
				SourceFile: v,
				DstFile:    strings.Replace(v, "tar.gz", "pdf", 1),
				Status:     "failed: " + err.Error(),
			}
		} else {
			ch <- util.Msg{
				SourceFile: v,
				DstFile:    strings.Replace(v, "tar.gz", "pdf", 1),
				Status:     "success",
			}
		}
	}
}

func ProcessTarGzFile(filename string) error {
	// 不包含最后的斜线
	dstdir := filename[:strings.LastIndex(filename, "/")]

	// 临时文件都放到这里，方便后面删除
	dstdir = dstdir + "/source"
	os.RemoveAll(dstdir)
	if err := os.Mkdir(dstdir, 0755); err != nil {
		return err
	}
	// 解压tar.gz，得到一个epub文件和一个jpg封面文件
	filenames, err := util.DeCompressTarGz(filename, dstdir, 2)
	if err != nil {
		return err
	}

	// 解压epub
	htmls, err := util.UnTarEpubFile(dstdir+"/"+filenames[0], dstdir)
	if err != nil {
		return err
	}

	// create pdf
	for k, html := range htmls {
		// OEBPS/Text/Chapter_3_3.xhtml

		outputPdf := dstdir + html[strings.LastIndex(html, "/"):strings.LastIndex(html, ".")] + ".pdf"

		if err := chromedp.HtmlToPdf(dstdir, html, outputPdf); err != nil {
			return err
		}

		// 转换成绝对地址
		htmls[k] = outputPdf
	}

	// 排序
	sort.Stable(SortStringSliceIncrement(htmls))

	// 合并pdf
	outputPdfName := strings.Replace(filename, "tar.gz", "pdf", 1)
	if err := pdfcpu.MergePdf(htmls, outputPdfName); err != nil {
		fmt.Println(err)
		return err
	}

	os.Remove(filename)
	os.RemoveAll(dstdir)

	return nil
}

type SortStringSliceIncrement []string

func (s SortStringSliceIncrement) Len() int { return len(s) }

func (s SortStringSliceIncrement) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SortStringSliceIncrement) Less(i, j int) bool {
	// OEBPS/Text/Chapter_3_3.xhtml
	// OEBPS/Text/Caver.xhtml
	//
	// OEBPS/Text/220125091.xhtml
	// OEBPS/Text/220125117.xhtml
	// OEBPS/Text/Caver.xhtml

	// Cover 排在第一位
	stri := s[i][strings.LastIndex(s[i], "/")+1 : strings.LastIndex(s[i], ".")]
	if strings.Contains(stri, "Cover") {
		return true
	}

	strj := s[j][strings.LastIndex(s[j], "/")+1 : strings.LastIndex(s[j], ".")]
	if strings.Contains(strj, "Cover") {
		return false
	}

	if strings.Contains(stri, "_") {
		arri := strings.Split(stri, "_")
		arrj := strings.Split(strj, "_")

		arri1, _ := strconv.Atoi(arri[1])
		arri2, _ := strconv.Atoi(arri[2])

		arrj1, _ := strconv.Atoi(arrj[1])
		arrj2, _ := strconv.Atoi(arrj[2])

		if arri1 < arrj1 {
			return true
		}
		if arri1 > arrj1 {
			return false
		}
		if arri1 == arrj1 {
			if arri2 < arrj2 {
				return true
			}
			if arri2 > arrj2 {
				return false
			}
		}
	} else {
		if stri < strj {
			return true
		} else {
			return false
		}
	}

	return false
}