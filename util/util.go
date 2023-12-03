package util

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	uuid "github.com/satori/go.uuid"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numLetterBytes = "0123456789"
const logRequestTime = 3
const DATE_FORMAT_MONTH = "2006-01"
const DATE_FORMAT_DAY = "2006-01-02"
const DATE_FORMAT_DAY_INT = "20060102"
const DATE_FORMAT_SECOND = "2006-01-02 15:04:05"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func Md5(str string) string {
	if str == "" {
		return ""
	}
	init := md5.New()
	init.Write([]byte(str))
	return fmt.Sprintf("%x", init.Sum(nil))
}

func Sha256(s string) string {
	if s == "" {
		return ""
	}
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Base64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Base64Decode(input string) string {
	str, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}

	return string(str)
}

func Base64SafeDecode(input string) string {
	input = strings.TrimRight(input, "=")
	input = strings.Replace(input, "+", "-", -1)
	input = strings.Replace(input, "/", "_", -1)
	str, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(str)
}

// 字符串转换整型
func IpString2Int(ipstring string) int64 {
	ipSegs := strings.Split(ipstring, ".")
	ipInt := 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return int64(ipInt)
}

// 整型转换成字符串
func IpInt2String(ipInt int) string {
	ipSegs := make([]string, 4)
	var length = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < length; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[length-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < length; i++ {
		buffer.WriteString(ipSegs[i])
		if i < length-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func randString(n int, LetterBytes string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(LetterBytes) {
			b[i] = LetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// CreateUuidStringNew Return a hex string, it's len is 2n.
func CreateUuidStringNew(n int) string {
	p := make([]byte, n)
	_, _ = crand.Read(p)
	return hex.EncodeToString(p)
}

func CreateUuidString() string {
	return uuid.NewV4().String()
}

func RandString(n int) string {
	return randString(n, letterBytes)
}

func RandNumString(n int) string {
	return randString(n, numLetterBytes)
}

// 生成以日期随机的id
func GenDateRandId() string {
	return time.Now().Format("20060102150405") + RandNumString(6)
}

var Rander = rand.New(src)

// 随机数生成
// @Param	min 	int	最小值
// @Param 	max		int	最大值
// @return  int		[min, max]
func RandInt(min int, max int) int {
	if min == max {
		return min
	}
	// Rander.Intn  --> [0, n)
	num := Rander.Intn(max-min+1) + min
	return num
}

// 获取服务器IP
func GetLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Get local IP addr failed!")
		return "127.0.0.1"
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func GetLocalIpPlus() string {
	ecsArray := map[string]string{
		"172.31.141.189": "opapi-001",
		"172.31.45.241":  "opapi-002",
		"172.31.45.238":  "opapi-003",
		"172.31.141.191": "opdb-001",
		"172.31.45.239":  "opdb-002",
		"172.31.45.240":  "opdb-003",
	}
	ip := GetLocalIp()
	if _, ok := ecsArray[ip]; ok {
		return ecsArray[ip] + " " + ip
	}
	return ip
}

// 是否是email
func IsEmail(email string) bool {
	if email == "" {
		return false
	}
	ok, _ := regexp.MatchString(`^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[0-9a-zA-Z]{2,3}$`, email)
	return ok
}

// --------------------------------------时间相关函数Start---------------------------------------

// 日期字符串 转 时间戳
// 2020-09-08 15:42:12  -->  1599550932
// 2020-09-08  -->  1599494400
func StrToTime(date string) int64 {
	f := DATE_FORMAT_SECOND
	if len(date) <= 10 {
		f = DATE_FORMAT_DAY
	}

	t, err := time.ParseInLocation(f, date, time.Local)
	if err != nil {
		return 0
	}

	return t.Unix()
}

// 返回某个时间戳的 零点时间戳
func GetDayTimeByUnix(u int64) int64 {
	dayStr := time.Unix(u, 0).Format(DATE_FORMAT_DAY)
	return StrToTime(dayStr)
}

// 返回当前日期 2021-06-08
func GetZeroDay() string {
	return time.Now().Format(DATE_FORMAT_DAY)
}

// 返回当前日期的时间戳 1623081600
func GetZeroUnix() int64 {
	return StrToTime(GetZeroDay())
}

// 获取某个月的第一天整型日期和最后一天整型日期
func FirstAndLastOfMonth(currentYear int, currentMonth int) (int, int, error) {
	if currentYear < 1970 || currentYear > 9999 {
		return 0, 0, errors.New("年份错误")
	}
	if currentMonth < 0 || currentMonth > 12 {
		return 0, 0, errors.New("月份错误")
	}
	// 设置时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, 0, errors.New("设置时区错误")
	}
	firstOfMonth := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, loc)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	f, err := strconv.Atoi(firstOfMonth.Format("20060102"))
	if err != nil {
		return 0, 0, errors.New("转换错误")
	}
	l, err := strconv.Atoi(lastOfMonth.Format("20060102"))
	if err != nil {
		return 0, 0, errors.New("转换错误")
	}
	return f, l, nil
}

// 活动时间函数转成前端需要的格式
func StartEndDateChange(startDate, endDate string) string {
	// 先将字符串转为时间戳
	startTime := StrToTime(startDate)
	endTime := StrToTime(endDate)

	// 格式化为前端需要的格式
	startDateStr := time.Unix(startTime, 0).Format("2006年1月2日")
	endDateStr := time.Unix(endTime, 0).Format("2006年1月2日")

	// 判断年份是否相同
	startYear := time.Unix(startTime, 0).Format("2006")
	endYear := time.Unix(endTime, 0).Format("2006")
	if startYear == endYear {
		endDateStr = time.Unix(endTime, 0).Format("1月2日")
	}
	return fmt.Sprintf("%s-%s", startDateStr, endDateStr)
}

// --------------------------------------时间相关函数End---------------------------------------

// 校验手机号
func CheckPhone(phone string) bool {
	if phone == "" {
		return false
	}
	exp := regexp.MustCompile(`^1[3456789]{1}\d{9}$`)
	res := exp.MatchString(phone)
	return res
}

type FetchResponseNew struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type FetchResponseOld struct {
	Status string      `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

func RemoveElemFromSlice(s []interface{}, elem interface{}) []interface{} {
	if len(s) == 0 {
		return s
	}
	for i, v := range s {
		if v == elem {
			s = append(s[:i], s[i+1:]...)
			return RemoveElemFromSlice(s, elem)
		}
	}
	return s
}

// 命令行模式下获取输入参数，flag 不适合
// ./voteapi vote2mqconsumer --action=run
func GetTerminalInput() (param map[string]string) {
	param = make(map[string]string)
	list := os.Args
	if len(list) <= 2 {
		return param
	}

	for _, v := range list {
		if strings.Contains(v, "--") {
			arr := strings.Split(v, "=")
			key := string(([]byte(arr[0]))[2:])
			param[key] = arr[1]
		}
	}

	return param
}

// 字符串反转
func StrReverse(str string) string {
	var result string
	strLen := len(str)
	for i := 0; i < strLen; i++ {
		result = result + fmt.Sprintf("%c", str[strLen-i-1])
	}
	return result
}

// in_array() string.
//
// Deprecated: Use InArray instead.
func InStringArray(item string, sli []string) (exist bool) {
	if len(sli) == 0 {
		return false
	}
	if item == "" {
		return false
	}
	for _, v := range sli {
		if v == item {
			return true
		}
	}
	return false
}

// in_array() int.
//
// Deprecated: Use InArray instead.
func InIntArray(item int, sli []int) (exist bool) {
	if len(sli) == 0 {
		return false
	}
	for _, v := range sli {
		if v == item {
			return true
		}
	}
	return false
}

func InArray[T comparable](item T, sli []T) (exist bool) {
	if len(sli) == 0 {
		return false
	}
	for _, v := range sli {
		if v == item {
			return true
		}
	}
	return false
}

func JoinSliceInt(nums []int) string {
	if len(nums) == 0 {
		return ""
	}

	s := make([]string, len(nums))
	for k, v := range nums {
		s[k] = strconv.Itoa(v)
	}

	return strings.Join(s, ",")
}

func Daemonize(args ...string) error {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	return cmd.Start()
}

// 计算签名 md5(按参数字典顺序排序的字符串值+密钥)
func GetSign(param map[string]string, appSecret string) string {
	var paramKey []string
	for k := range param {
		paramKey = append(paramKey, k)
	}
	sort.Strings(paramKey)
	var paramVal string
	for _, v := range paramKey {
		paramVal += param[v]
	}
	return Md5(paramVal + appSecret)
}

// 计算签名：https://pay.weixin.qq.com/wiki/doc/api/tools/miniprogram_hb.php?chapter=4_3
func GetWxSign(param map[string]string, appSecret string) string {
	var paramKey []string
	for k := range param {
		paramKey = append(paramKey, k)
	}
	sort.Strings(paramKey)
	var paramVal, dot string
	for _, v := range paramKey {
		paramVal += dot + v + "=" + param[v]
		dot = "&"
	}
	return Md5(paramVal + "&key=" + appSecret)
}

func PrettyPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return
	}

	fmt.Println(out.String())
}

/**
 * @desc    判断手机号是否是虚/物/卫号段：https://zhuanlan.zhihu.com/p/40878967
 * @date    2021-04-06 14:56:50
 */
func IsVirtualNumber(phone string) bool {
	ok, err := regexp.Match("^(1703|1705|1706|165|144|1704|1708|171|167|140|1700|1701|1702|162|1740|141)", []byte(phone))
	if err != nil {
		// 异常也标记为虚拟号
		return true
	}
	return ok
}

// 截取字符串，支持多字节字符
// start：起始下标，负数从从尾部开始，最后一个为-1
// length：截取长度，负数表示截取到末尾
func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}
	// 允许从尾部开始计算
	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}
	// 到末尾
	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}

	return
}

//func PassportKeyOp(str string) (result string) {
//	encrypt_key := "bookan@7788414"
//	encryptKey := Md5(encrypt_key)
//	init := 0
//	result = ""
//	l := len(str)
//	encryptKeyLen := len(encryptKey)
//	for i := 0; i < l; i++ {
//		if init == encryptKeyLen {
//			init = 0
//		}
//		result += string(str[i] ^ encryptKey[init])
//		init++
//	}
//
//	return result
//}

//func EncryptOp(str string) (result string) {
//	// 用于加密的key，使用随机数发生器产生 0~32000 的值并 MD5()
//	encryptKey := Md5(strconv.Itoa(RandInt(0, 32000)))
//	init := 0 // 初始化变量长度
//	result = ""
//	strLen := len(str)               // 待加密字符串的长度
//	encryptKeyLen := len(encryptKey) // 加密key的长度
//
//	for i := 0; i < strLen; i++ {
//		// 如果 $init = $encryptKey 的长度, 则 $init 清零
//		if init == encryptKeyLen {
//			init = 0
//		}
//
//		// $tmp 字串在末尾增加两位, 其第一位内容为 $encryptKey 的第 $init 位，
//		// 第二位内容为 $string 的第 $index 位与 $encryptKey 的 $init 位取异或。然后 $init = $init + 1
//		result += string(encryptKey[init]) + string(str[i]^encryptKey[init])
//		init++
//	}
//
//	// 返回结果，结果为 passportKeyOp() 函数返回值的 base64 编码结果
//	return Base64(PassportKeyOp(result))
//}

// WWcJbAQ5VDMPY1FhVjE= 对应 6292352
//func DecryptOp(str string) (result string) {
//	// 把空格替换为+号，+号在传递过程中变成了空格
//	str = strings.ReplaceAll(str, " ", "+")
//	str = PassportKeyOp(Base64Decode(str))
//	result = ""
//	l := len(str)
//	for index := 0; index < l; index++ {
//		if index+1 > l {
//			return ""
//		}
//
//		result += string(str[index] ^ str[index+1])
//		index++
//	}
//
//	return result
//}

// 随机数
func MtRand(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min+1) + min
}

// php加密函数迁移过来
func EncryptOp(str string) string {
	encryptKey := strconv.Itoa(int(MtRand(0, 32000)))
	init := 0
	tmp := ""
	strPlus := []rune(str)
	strLen := len(strPlus)
	encryptKeyPlus := []rune(encryptKey)
	encryptKeyLen := len(encryptKeyPlus)
	for index := 0; index < strLen; index++ {
		if init == encryptKeyLen {
			init = 0
		}
		strInt := int(strPlus[index])
		encryptKeyInt := int(encryptKeyPlus[init])
		tmp += string(encryptKeyPlus[init]) + string(rune(strInt^encryptKeyInt))

		init++
	}
	sign := passportKeyOp(tmp)
	sEnc := base64.StdEncoding.EncodeToString([]byte(sign))
	return sEnc
}

// php加密解密函数迁移过来
// WWcJbAQ5VDMPY1FhVjE= 对应 6292352
func passportKeyOp(str string) string {
	key := "bookan@7788414"
	encryptKey := Md5(key) // 加密的key
	init := 0
	result := ""
	strPlus := []rune(str)
	strLen := len(strPlus)

	encryptKeyPlus := []rune(encryptKey)
	encryptKeyLen := len(encryptKeyPlus)

	for index := 0; index < strLen; index++ {
		if init == encryptKeyLen {
			init = 0
		}
		result += string(rune(int(strPlus[index]) ^ int(encryptKeyPlus[init])))
		init++
	}
	return result
}

// php解密函数迁移过来
func DecryptOp(str string) string {
	// 把空格替换为+号，+号在传递过程中变成了空格
	str = strings.ReplaceAll(str, " ", "+")
	sDec, _ := base64.StdEncoding.DecodeString(str)
	str = passportKeyOp(string(sDec))
	result := ""
	strPlus := []rune(str)
	strLen := len(strPlus)
	for index := 0; index < strLen-1; index++ {
		result += string(rune(int(strPlus[index]) ^ int(strPlus[index+1])))
		index++
	}
	return result
}

// 判断是否为图片
func IsPhoto(url string) bool {
	if url == "" {
		return false
	}
	u := strings.Split(url, ".")
	ext := u[len(u)-1]
	e := []string{"jpg", "jpeg", "png"}
	for _, i := range e {
		if ext == i {
			return true
		}
	}

	return false
}

func IsMp4(url string) bool {
	if url == "" {
		return false
	}
	u := strings.Split(url, ".")
	ext := strings.ToLower(u[len(u)-1])
	e := []string{"mp4", "avi", "flv", "mkv", "wmv", "mov", "asf", "navi", "ra", "ram", "rm", "rmvb", "m4v", "f4v"}
	for _, i := range e {
		if ext == i {
			return true
		}
	}

	return false
}

func IsMp3(url string) bool {
	if url == "" {
		return false
	}
	u := strings.Split(url, ".")
	ext := u[len(u)-1]
	e := []string{"mp3"}
	for _, i := range e {
		if ext == i {
			return true
		}
	}

	return false
}

func GetFileExtName(url string) string {
	if url == "" {
		return ""
	}
	u := strings.Split(url, ".")
	ext := u[len(u)-1]

	return ext
}

// 获取程序绝对路径
func GetAbsolutePath() string {
	ePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	root := path.Dir(ePath)
	return root
}

func StringSliceUnique(in []string) []string {
	out := make([]string, 0)
	if len(in) == 0 {
		return out
	}

	for _, i := range in {
		tmp := true
		for _, j := range out {
			if i == j {
				tmp = false
				break
			}
		}
		if tmp {
			out = append(out, i)
		}
	}

	return out
}

func NickNameHide(str string) string {
	if str == "" {
		return ""
	}
	strLen := utf8.RuneCountInString(str)
	if strLen <= 1 {
		return str
	}
	if strLen >= 4 {
		strLen = 4
	}
	start := string([]rune(str)[0:1])
	end := strings.Repeat("*", strLen-1)
	return start + end
}

func IsNumber(str string) bool {
	if str == "" {
		return false
	}

	ok, err := regexp.MatchString(`^\d+$`, str)
	if err != nil {
		return false
	}

	return ok
}

func IsNumberAndChar(str string) bool {
	if str == "" {
		return false
	}
	ok, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, str)
	if err != nil {
		return false
	}

	return ok
}

// Round1 returns the rounded value of x to specified precision, use ROUND_HALF_UP mode.
// For example:
//
//		  Round1(0.6635, 3)   // 0.664
//	   Round1(0.363636, 2) // 0.36
//	   Round1(0.363636, 1) // 0.4
func RoundWithPrecision(x float64, precision int) float64 {
	if precision == 0 {
		return math.Round(x)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Round(x*p) * math.Pow10(-precision)
	}
	return math.Round(x*p) / p
}

// 压缩文件 .zip
// src 可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func CompressZip(src string, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	files := []*os.File{f}
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compressZip(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compressZip(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		// 增加对空目录的判断
		if len(fileInfos) <= 0 {
			header, err := zip.FileInfoHeader(info)
			header.Name = prefix
			if err != nil {
				fmt.Println("error is:" + err.Error())
				return err
			}
			_, err = zw.CreateHeader(header)
			if err != nil {
				fmt.Println("create error is:" + err.Error())
				return err
			}
			file.Close()
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compressZip(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 解压
func DeCompressZip(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

// .tar.gz
func CompressTarGz(filename string, sourceDir string) {
	// file write
	fw, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 打开文件夹
	dir, err := os.Open(sourceDir)
	if err != nil {
		panic(nil)
	}
	defer dir.Close()

	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// 遍历文件列表
	for _, fi := range fis {
		// 逃过文件夹, 我这里就不递归了
		if fi.IsDir() {
			continue
		}

		// 打印文件名称
		fmt.Println(fi.Name())

		// 打开文件
		fr, err := os.Open(dir.Name() + "/" + fi.Name())
		if err != nil {
			panic(err)
		}
		defer fr.Close()

		// 信息头
		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()

		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			panic(err)
		}

		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("tar.gz ok")
}

// .tar.gz
func DeCompressTarGz(filename string, dstDir string, mode int) (filenames []string, err error) {
	// file read
	fr, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	// tar read
	tr := tar.NewReader(gr)

	if !strings.HasSuffix(dstDir, "/") {
		dstDir = dstDir + "/"
	}

	filenames = make([]string, 0)

	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		var do = true

		if mode == 1 && !strings.Contains(h.Name, "_big") {
			// 只保留 _big 文件
			do = false
		}
		if mode == 2 && !strings.Contains(h.Name, ".epub") {
			// 只保留 epub 文件
			do = false
		}

		if do {
			filenames = append(filenames, h.Name[1:]) // 第一个字符为斜线 /，此处将其去掉

			newFile, err := os.OpenFile(dstDir+h.Name, os.O_CREATE|os.O_RDWR, 0777)
			if err != nil {
				panic(err)
			}
			defer newFile.Close()

			// 写文件
			_, err = io.Copy(newFile, tr)
			if err != nil {
				panic(err)
			}
		}
	}

	return
}

// .tar
func CompressTar() {

}

// .tar
func DeCompressTar() {

}

func FindTarGzFile(dir string) (res []string, err error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return res, err
	}

	res = make([]string, 0)

	for _, f := range fs {
		if f.IsDir() {
			r, _ := FindTarGzFile(dir + "/" + f.Name())
			if len(r) > 0 {
				res = append(res, r...)
			}
		} else {
			if strings.Contains(f.Name(), ".tar.gz") {
				res = append(res, dir+"/"+f.Name())
			}
		}
	}

	return
}

func UnTarEpubFile(epubFile string, dstDir string) (htmls []string, err error) {
	file, _ := os.Open(epubFile)
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return htmls, err
	}

	htmls = make([]string, 0)

	// 解压缩EPUB文件
	r, _ := zip.NewReader(file, info.Size())
	for _, f := range r.File {
		htmlFile, _ := f.Open()
		defer htmlFile.Close()

		// OEBPS/Text/Chapter_3_3.xhtml
		if strings.Contains(f.Name, "/") {
			os.MkdirAll(dstDir+"/"+f.Name[:strings.LastIndex(f.Name, "/")], 0777)
		}

		newFile, _ := os.OpenFile(dstDir+"/"+f.Name, os.O_CREATE|os.O_RDWR, 0777)
		defer newFile.Close()

		_, err = io.Copy(newFile, htmlFile)
		if err != nil {
			return htmls, err
		}
	}

	// 解析 directories.json 文件，里面包含了目录信息，对 xhtml 已经做了排序
	dictFile, _ := os.Open(dstDir + "/directories.json")
	defer dictFile.Close()
	dictBt, _ := io.ReadAll(dictFile)

	var dic Directory
	err = json.Unmarshal(dictBt, &dic)
	if err != nil {
		return htmls, err
	}
	for _, sp := range dic.Spine {
		// 有时候目录中的文件实际上不存在
		if _, err := os.Stat(dstDir + "/" + sp.Src); err == nil {
			htmls = append(htmls, sp.Src)
		}
	}

	return
}
