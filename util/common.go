package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net"

	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-micro/constant"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/nfnt/resize"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v2"
)

var (
	Loc, _ = time.LoadLocation("Asia/Shanghai")
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
)

func Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func EncodeString(data interface{}) string {
	if binData, err := json.Marshal(data); err == nil {
		return string(binData)
	}
	return ""
}

func Decode(str string, data interface{}) error {
	return json.UnmarshalFromString(str, data)
}

func DecodeByte(b []byte, data interface{}) error {
	return json.Unmarshal(b, data)
}

func StringUUID() string {

	return strings.ReplaceAll(uuid.New().String(), "-", "")

}

func IntUUID() uint32 {
	return uuid.New().ID()

}

func InArray(val string, array []string) bool {
	for i := 0; i < len(array); i++ {
		// if strings.HasPrefix(val, array[i]) {
		if val == array[i] {
			return true
		}
	}
	return false
}

func PrefixInArray(val string, array []string) bool {
	for i := 0; i < len(array); i++ {
		if strings.HasPrefix(val, array[i]) {
			return true
		}
	}
	return false
}

func FindStrIndex(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func FindIntIndex(slice []int64, val int64) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

// RandomStr 获得随机字符串
func RandomStr(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// func FindIndex(slice []string, val string) int {
// 	for i, item := range slice {
// 		if item == val {
// 			return i
// 		}
// 	}
// 	return -1
// }

func checkExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
func GetExt(fileName string) string {
	return strings.TrimLeft(path.Ext(fileName), ".")
}
func CheckImageExt(fileName string, exts []string) (bool, string) {
	ext := GetExt(fileName)
	for _, allow := range exts {
		if strings.EqualFold(allow, ext) {
			return true, ext
		}
	}
	return false, ext
}

func GetSize(f multipart.File) (int, error) {
	if content, err := io.ReadAll(f); err != nil {
		return 0, err
	} else {
		return len(content), nil
	}
}

// 创建目录（可以是多级目录，任何一级不存在则创建）
func TouchDir(dir string) error {
	if !checkExists(dir) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func CreateFile(name, sep string) *os.File {
	if !checkExists(name) {
		r := []rune(name)
		if n := strings.LastIndex(name, sep); n != -1 {
			if err := os.MkdirAll(string(r[:n]), os.ModePerm); err != nil {
				log.Fatalln(`Mkdir`, err)
			}
		}
	}

	f, err := os.Create(name) //创建文件
	if err != nil {
		log.Fatalln(`file create`, err)
	}
	return f
}

func OpenFile(name, sep string) *os.File {
	if !checkExists(name) {
		r := []rune(name)
		if n := strings.LastIndex(name, sep); n != -1 {
			if err := os.MkdirAll(string(r[:n]), os.ModePerm); err != nil {
				log.Fatalln(`Mkdir`, err)
			}
		}
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(`file create`, err)
	}
	return f
}

func ToCamel(str string) (result string) {
	rn := []rune(str)

	if rn[0] >= 97 && rn[0] <= 122 {
		result = string(rn[0] - 32)
	}

	l := len(rn)
	for i := 1; i < l; i++ {
		if rn[i] == 95 && rn[i+1] >= 97 && rn[i+1] <= 122 { //过滤下划线
			rn[i+1] -= 32
		} else {
			result += string(rn[i])
		}
	}
	return
}

// 驼峰式变下划线式
func ToSnake(str string) (result string) {
	rn := []rune(str)

	if rn[0] >= 65 && rn[0] <= 90 {
		result = string(rn[0] + 32)
	}

	n := len(rn)
	for i := 1; i < n; i++ {
		if rn[i] >= 65 && rn[i] <= 90 { //大写变小写
			rn[i] += 32
			result += "_"
		}
		result += string(rn[i])
	}
	return
}

// start>0时从串首开始截取n位，start<0时从串尾倒数|start|开始截取n位
func Substr(s string, start int, n int) string {
	rs := []rune(s)
	sn := len(rs)

	if n > sn || n < 0 {
		return ""
	}

	if start < 0 {
		if sn+start < 0 {
			return ""
		}
		start = sn + start
	} else {
		if start > n {
			return ""
		}
	}

	if n > sn-start {
		n = sn - start
	}

	return string(rs[start : start+n])
}

func ASCII(r rune) rune {
	switch {
	case 97 <= r && r <= 122:
		return r - 32
	case 65 <= r && r <= 90:
		return r + 32
	default:
		return r
	}
}

func IndexString(str string, sep rune, level int) string {
	rs := []rune(str)
	var buffer bytes.Buffer
	var n int
	for i, ln := 0, len(rs); i < ln; i++ {
		if rs[i] == sep {
			n += 1
		}
		if n == level {
			break
		}
		buffer.WriteRune(rs[i])
	}
	return buffer.String()
}

func LastIndexString(src, spec string) string {
	s := strings.Split(src, spec)
	if n := len(s); n > 1 {
		return s[n-2]
	}
	return ""
}

func IsEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func MillisecondToDateString(ms int64) string {
	return time.Unix(ms, 0).Format("2006-01-02")
}

func MillisecondToDateHMS(ms int64) string {
	return time.Unix(ms, 0).Format(constant.FORMAT_ISO8601_DATE_TIME)
}

// func ListMap(rows *orm.Rows, call func(map[string]interface{}) (string, string)) (result []map[string]interface{}) {
// 	for rows.Next() {
// 		tmp := make(map[string]interface{})
// 		rows.MapScan(tmp)
// 		for k, encoded := range tmp {
// 			switch encoded.(type) {
// 			case []byte:
// 				tmp[k] = string(encoded.([]byte))
// 			}
// 		}
// 		if call != nil {
// 			key, res := call(tmp)
// 			tmp[key] = res
// 		}
// 		result = append(result, tmp)
// 	}
// 	return
// }

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(t string) (int64, error) {
	if tm, err := time.ParseInLocation(constant.FORMAT_ISO8601_DATE_TIME, t, Loc); err != nil {
		return 0, nil
	} else {
		tm = tm.AddDate(0, 0, -tm.Day()+1)
		return GetZeroTime(tm).Unix(), nil
	}
}

// 获取传入的时间所在月份的第一天凌晨与最后一天23点59分59秒
func GetFirstAndLastDateOfMonth(t string) (int64, int64, error) {
	if tm, err := time.ParseInLocation(constant.FORMAT_ISO8601_DATE_TIME, t, Loc); err != nil {
		return 0, 0, nil
	} else {
		tm = tm.AddDate(0, 0, -tm.Day()+1)
		return GetZeroTime(tm).Unix(), constant.DAY_TS_DIFF + GetZeroTime(tm).AddDate(0, 1, -1).Unix(), nil
	}
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, Loc)
}

func GetZeroTSByTS(ms int64) int64 {
	dt := time.Unix(ms, 0)
	y, m, d := dt.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, Loc).Unix()
}

func GetMonthEndTSByTS(ms int64) int64 {
	dt := time.Unix(ms, 0)
	return dt.AddDate(0, 1, -1).Unix() + constant.DAY_TS_DIFF
}

func GetZeroTS() int64 {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, Loc).Unix()
}

func GetTS(t string) (int64, int64) {
	ts, _ := time.ParseInLocation(constant.FORMAT_ISO8601_DATE_TIME, t+" 00:00:00", Loc)
	min := ts.Unix()
	return min, min + constant.DAY_TS_DIFF
}

func GetDiffTS(diff int) (int64, int64) {
	t := time.Now()
	dt := t.AddDate(0, 0, diff)
	y, m, d := dt.Date()
	min := time.Date(y, m, d, 0, 0, 0, 0, Loc).Unix()
	return min, t.Unix()
}

func Contain(p, c string) bool {
	return strings.Contains(","+p+",", ","+c+",")
}

// s1是否包含s2; 其中s1为逗号隔开的字符串，s2为可能出现在s1中，逗号分隔的某个子串;
// any 是否只包含其中一个即可，还是必须包含所有
func SubContain(s1 string, any bool, s2 ...string) bool {
	if any {
		for _, e := range s2 {
			if Contain(s1, e) {
				return true
			}
			continue
		}
		return false
	} else {
		for _, e := range s2 {
			if !Contain(s1, e) {
				return false
			}
			continue
		}
		return true
	}
}

func SliceRemove(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}

func String2Int8(str string) int8 {
	istr, err := strconv.ParseInt(str, 10, 8)
	if err == nil {
		return int8(istr)
	}
	return 0
}

func String2Int32(str string) int32 {
	istr, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		return int32(istr)
	}
	return 0
}

func String2Int64(str string) int8 {
	istr, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return int8(istr)
	}
	return 0
}

// OrdinalWeekPeriod 获取某年某周的起始日期
func OrdinalWeekPeriod(year, week int) (start, end time.Time) {
	if year > 0 && week > 0 {
		start = time.Date(year, 1, 0, 0, 0, 0, 0, Loc)
		// 第一天是周几
		firstWeekday := int(start.AddDate(0, 0, 1).Weekday())
		// 当年第一周有几天
		firstWeekDays := 1
		if firstWeekday != 0 {
			firstWeekDays = 7 - firstWeekday + 1
		}
		if week == 1 {
			end = start.AddDate(0, 0, firstWeekDays)
		} else {
			end = start.AddDate(0, 0, firstWeekDays+(week-1)*7)
			start = end.AddDate(0, 0, -7)
		}
	}
	return
}

// OrdinalWeek 当前日期为该年第几周周几，星期天为0
func OrdinalWeek(t time.Time) (year, week, weekday int) {
	year = t.Year()
	weekday = int(t.Weekday())
	yearDays := t.YearDay()
	firstDay := t.AddDate(0, 0, -yearDays+1)
	firstWeekday := int(firstDay.Weekday())
	// 当年第一周有几天
	firstWeekDays := 1
	if firstWeekday != 0 {
		firstWeekDays = 7 - firstWeekday + 1
	}
	if yearDays <= firstWeekDays {
		week = 1
	} else {
		week = (yearDays-firstWeekDays)/7 + 2
	}
	return
}

func HJ212RegExtract(text []byte) map[string]string {
	cpReg := regexp.MustCompile(`CP=&&(.+)&&`)
	infoReg := regexp.MustCompile(`([^=]+)=([^;]+);?`)

	var tmp []byte
	if tmp = cpReg.ReplaceAll(text, []byte("")); len(tmp) < 6 {
		return nil
	}

	info := infoReg.FindAllSubmatch(tmp[6:len(tmp)-7], -1)
	res := map[string]string{}
	if regData := cpReg.FindSubmatch(text)[1]; len(regData) != 0 {
		cp := regexp.MustCompile(`([^=]+)=([^;,]+)[;,]?`).FindAllSubmatch(regData, -1)
		for _, v := range cp {
			res[string(v[1])] = string(v[2])
		}
	}

	for _, v := range info {
		res[string(v[1])] = string(v[2])
	}
	return res
}

// GetDiffTen 获取某一天00:00:00至23:59:59，以10min为单位的时间戳起止序列
func GetDiffTen(cur time.Time) (int64, int64) {
	t1 := GetZeroTime(cur).Unix() / 600
	return t1, t1 + 143
}

func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func Number2String(n interface{}) string {
	// fmt.Println(reflect.TypeOf(n))
	switch n := n.(type) {
	case int:
		return strconv.Itoa(n)
	case int32:
		return strconv.FormatInt(int64(n), 10)
	case int64:
		return strconv.FormatInt(n, 10)
	case float32:
		//return fmt.Sprint(n.(float32))
		return strconv.FormatFloat(float64(n), 'f', -1, 32)
	case float64:
		//return fmt.Sprint(n.(float64))
		return strconv.FormatFloat(n, 'f', -1, 64)
	default:
		return ""
	}
}

func IfNull(str interface{}, sep string) string {
	if str != nil && str.(string) != "" {
		// return sep + str + sep
		return str.(string)
	}
	return sep
}

func SortRange(m map[string]interface{}, f func(int, string)) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		f(i, k)
	}
}

func HasField(st reflect.Value, name string) bool {

	if s := st.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) }); s.IsValid() {
		return true
	}
	return false
}

func FieldByName(st reflect.Value, name string) reflect.Value {
	if s := st.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) }); s.IsValid() {
		return s
	}
	return reflect.Value{}
}

func SetFieldByName(st reflect.Value, name string, val interface{}) bool {

	// if st.Kind() == reflect.Ptr {
	// 	st = st.Elem()
	// 	st = reflect.Indirect(st)
	// 	if s := st.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) }); s.IsValid() {
	// 		if stype := s.Type(); stype == reflect.TypeOf(val) {
	// 			s.Set(reflect.ValueOf(val))
	// 		} else {
	// 			s.Set(reflect.ValueOf(val).Convert(stype))
	// 		}
	// 		return true
	// 	}
	// }
	if s := st.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) }); s.IsValid() {
		if stype := s.Type(); stype == reflect.TypeOf(val) {
			s.Set(reflect.ValueOf(val))
		} else {
			s.Set(reflect.ValueOf(val).Convert(stype))
		}
		return true
	}
	return false
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		if vm, ok := v.(map[string]interface{}); ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// key: tag中的键值，delimiter: tag中的分隔符，structs：要生成表的结构体
func StructToTable(key, delimiter string, toCamel bool, structs ...interface{}) {
	var sb strings.Builder
	var fieldName string
	var attrs []string
	for _, sct := range structs {
		s := reflect.TypeOf(sct)
		sb.WriteString(`CREATE TABLE `)
		sb.WriteString(s.Name())
		sb.WriteString(" (\n")
		fn := s.NumField()
		for i := 0; i < fn; i++ {
			sb.WriteString("    ")
			attrs = nil
			if tag := string(s.Field(i).Tag.Get(key)); tag == "" {
				if toCamel {
					fieldName = ToCamel(s.Field(i).Name)
				} else {
					fieldName = s.Field(i).Name
				}
			} else {
				attrs = strings.Split(tag, delimiter)
				fieldName = attrs[0]
			}
			sb.WriteString(fieldName)
			sb.WriteString(" ")
			switch s.Field(i).Type.Name() {
			case "int8":
				sb.WriteString("TINYINT")
			case "int", "int16", "int32":
				sb.WriteString("INT")
			case "int64":
				sb.WriteString("BIGINT")
			case "string":
				fallthrough
			default:
				sb.WriteString("VARCHAR(50)")
			}

			if len(attrs) > 1 {
				sb.WriteString(" ")
				sb.WriteString(strings.Join(attrs[1:], " "))
			}

			if i+1 != fn {
				sb.WriteString(",")
			}
			sb.WriteString("\n")
		}
		sb.WriteString(");\n\n")
	}
	fmt.Println(sb.String())
}

func CompressJPEG(img image.Image, width uint, filePath string) error {
	// defer file.Close()
	// img, err := jpeg.Decode(file)
	m := resize.Resize(width, 0, img, resize.Bilinear)

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	return jpeg.Encode(out, m, nil)
}

// IsIPv4 check if the string is an IP version 4.
func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ".")
}

// IsIPv6 check if the string is an IP version 6.
func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

func Min(first int, args ...int) int {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func AnyToBytes(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

// 返回 interface{} 该为某个结构体的指针类型
// func bytesToAny(val []byte, data interface{}) (err error) {
// 	buf := new(bytes.Buffer)
// 	dec := gob.NewDecoder(buf)
// 	err = dec.Decode(data)
// 	return
// }

func BytesToAny[T any](val []byte) (data T, err error) {
	err = msgpack.Unmarshal(val, &data)
	return
}

func Loadyaml(path string, conf interface{}) {
	c, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.UnmarshalStrict(c, conf)
	if err != nil {
		log.Fatalln(err)
	}
}
