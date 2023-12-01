package chromedp

import (
	"context"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func HtmlToPdf(sourceDir string, html string, outputPdf string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, PrintToPDF(sourceDir+"/"+html, &buf)); err != nil {
		return err
	}

	if err := os.WriteFile(outputPdf, buf, 0777); err != nil {
		return err
	}

	return nil
}

// print a specific pdf page.
func PrintToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			pdf := page.PrintToPDF()
			pdf.MarginTop = 0.6
			pdf.MarginLeft = 0.8
			pdf.MarginRight = 0.8
			pdf.MarginBottom = 0.6
			pdf.PaperHeight = 11.2 // 默认值11会导致封面被分割成了两页，因此要调大一点
			buf, _, err := pdf.WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}