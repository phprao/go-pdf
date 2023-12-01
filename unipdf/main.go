package unipdf

/*

go get github.com/unidoc/unipdf/v3

*/

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/phprao/go-pdf/util"
)

var baseDir = "F:/jx/20231114_4359/periodical/resource/jpg"

// baseDir+"/jpage3/41287/41287-358053/41287_358053.tar.gz"

func Run() {
	dir := baseDir

	tarGzFiles, err := util.FindTarGzFile(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range tarGzFiles {
		if err := ProcessTarGzFile(v); err != nil {
			log.Panicln(v, ": ", err)
		} else {

		}
	}
}

func ProcessTarGzFile(filename string) error {
	// 不包含最后的斜线
	dstdir := filename[:strings.LastIndex(filename, "/")]

	// 临时文件都放到这里，方便后面删除
	dstdir = dstdir + "/source"
	if err := os.Mkdir(dstdir, 0755); err != nil {
		return err
	}

	// 解压
	filenames, err := util.DeCompressTarGz(filename, dstdir, 1)
	if err != nil {
		return err
	}

	// 排序，并加上前缀
	nodes := strings.Split(dstdir, "/")
	params := strings.Split(nodes[len(nodes)-2], "-")
	rid, _ := strconv.Atoi(params[0])
	iid, _ := strconv.Atoi(params[1])
	sortedFilenames := sortByName(filenames, int64(iid), int64(rid), dstdir)

	// 合成PDF
	outputPdfName := strings.Replace(filename, "tar.gz", "pdf", 1)
	if err := ImagesToPdf(sortedFilenames, outputPdfName); err != nil {
		return err
	}

	// 删除压缩包，删除图片文件
	// os.Remove(filename)
	os.RemoveAll(dstdir)

	return nil
}

func sortByName(filenames []string, iid int64, rid int64, dstDir string) (dst []string) {
	h := util.NewHash()
	ok, data := h.GetHash(iid, rid, 1, 1000)
	if ok == 0 {
		return
	}

	dst = make([]string, len(filenames))

	for _, h := range filenames {
		str := strings.Replace(h, "_big.jpg", "", 1)
		for k := range data {
			if str == data[k].Hash {
				dst[data[k].Page-1] = dstDir + "/" + h
			}
		}
	}

	return
}
