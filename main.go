package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"path/filepath"

	"github.com/phprao/ColorOutput"
	"github.com/phprao/go-pdf/epub"
	"github.com/phprao/go-pdf/jpg"
	"github.com/phprao/go-pdf/util"
)

var wg sync.WaitGroup

func main() {
	args := os.Args

	if len(args) < 2 || args[1] == "" || !filepath.IsAbs(args[1]) || !strings.Contains(args[1], "resource") {
		log.Fatal("请指定resource目录的绝对路径, 比如 F:/jx/20231114_4359/periodical/resource")
	}

	root := strings.TrimRight(strings.ReplaceAll(args[1], "\\", "/"), "/")

	ch := make(chan util.Msg, 0)
	wg.Add(2)

	go func() {
		epub.Run(root, ch)
		wg.Done()
	}()

	go func() {
		jpg.Run(root, ch)
		wg.Done()
	}()

	go func() {
		PrintMsg(root, ch)
	}()

	wg.Wait()

	close(ch)

	time.Sleep(3 * time.Second)

	log.Println("done.")
}

func PrintMsg(root string, ch chan util.Msg) {
	logfile := root + "/logs.log"

	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Println("日志文件创建失败：" + logfile)
		return
	}
	defer f.Close()

	for msg := range ch {
		str := fmt.Sprintf("[%s] %s ---> %s ---> %s", time.Now().Format(util.DATE_FORMAT_SECOND), msg.SourceFile, msg.DstFile, msg.Status)
		io.WriteString(f, str+"\n")

		if msg.Status != "success" {
			ColorOutput.Colorful.WithFrontColor("red").Println(str)
		} else {
			fmt.Println(str)
		}
	}
}
