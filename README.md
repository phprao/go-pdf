### 关于 unipdf
安装依赖
go get github.com/unidoc/unipdf/v3

示例代码
https://github.com/unidoc/unipdf-examples

获取KEY
登录 https://cloud.unidoc.io/ 注册账号，生成 KEY，但是需要收费。

### 关于 chromedp
使用Golang编写，主要功能是调用浏览器内核来渲染HTML页面，也可以用它来在页面上做一些操作，还有一个附加功能是将渲染后的页面保存为PDF文件。

### 关于pdfcpu
使用Golang编写，主要功能是操作PDF文件，功能比较齐全。

### 关于gofpdf
已经停止了维护，使用Golang编写，主要功能是操作PDF文件，功能比较少。

### 案例
1、将多个jpg文件合并到一个PDF文件中去，直接使用 pdfcpu。

2、将 epub 转换成一个pdf：先将epub解压，得到xhtml，然后使用chromedp将xhtml转换成pdf，最后调用pdfcpu将多个pdf合并成一个pdf文件。
