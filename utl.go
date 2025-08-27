package utl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

func RndStr(l int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Rows(path string) ([]string, error) {
	var client http.Client
	resp, err := client.Get("https://sarnet.ru/" + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var body string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		body = string(bodyBytes)
	}
	rows := strings.Split(body, "\n")
	return rows, nil
}

func Cols(row string) []string {
	space := regexp.MustCompile(`\s+`)
	row = space.ReplaceAllString(row, " ")
	return strings.Split(row, " ")
}

func IsNum(str string) bool {
	if len(str) == 0 {
		return false
	}
	if len(str) > 1 && (str[0] == '-' || str[0] == '+') {
		str = str[1:]
	}
	for _, s := range str {
		if s != '.' && (s < '0' || s > '9') {
			return false
		}
	}
	return true
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func IsIP(str string) bool {
	parts := strings.Split(str, ".")
	if len(parts) != 4 {
		return false
	}
	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func Localip() string {
	client := &http.Client{}
	resp, err := client.Get("https://sarnet.ru/ip/")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(body))
}

func Localid() int {
	client := &http.Client{}
	resp, err := client.Get("https://sarnet.ru/id/")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	ret, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		return 0
	}
	return int(ret)
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func Remove[T comparable](slice []T, i int) []T {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func Shuffle[S ~[]T, T any](items S) {
	if len(items) <= 1 {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	swap := func(i, j int) {
		items[i], items[j] = items[j], items[i]
	}
	r.Shuffle(len(items), swap)
}

func Exec(com string) ([]byte, error) {
	arr := strings.Split(com, " ")
	cmd := exec.Command(arr[0], arr[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}

func CutString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n%2 == 0 {
		return s[:n/2-2] + "..." + s[len(s)-n/2-1:]
	}
	return s[:n/2-1] + "..." + s[len(s)-n/2-1:]
}

func PrettyPrint(v any) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func PrettyString(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return ""
}

func LA() (la1, la5, la15 float64) {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return
	}
	defer file.Close()
	fmt.Fscanf(file, "%f %f %f", &la1, &la5, &la15)
	return
}

func UP() (up float64) {
	b, _ := os.ReadFile("/proc/uptime")
	str := string(b)
	fmt.Sscanf(str, "%f", &up)
	return
}

func FD() int {
	b, _ := os.ReadFile("/proc/sys/fs/file-nr")
	str := string(b)
	var fd int
	fmt.Sscanf(str, "%d", &fd)
	return fd
}

func SOC() int {
	b, _ := os.ReadFile("/proc/net/sockstat")
	str := string(b)
	var soc int
	fmt.Sscanf(str, "sockets: used %d", &soc)
	return soc
}

func CPU() (idle, total uint64) {
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.SplitSeq(string(contents), "\n")
	for line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					val = 0
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return
}

func MEM() (memtotal, memfree, swaptotal, swapfree, tmp uint64) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if fields[0] == "MemTotal:" {
			memtotal, _ = strconv.ParseUint(fields[1], 10, 64)
			continue
		}
		if fields[0] == "MemAvailable:" {
			memfree, _ = strconv.ParseUint(fields[1], 10, 64)
			continue
		}
		if fields[0] == "SwapTotal:" {
			swaptotal, _ = strconv.ParseUint(fields[1], 10, 64)
			continue
		}
		if fields[0] == "SwapFree:" {
			swapfree, _ = strconv.ParseUint(fields[1], 10, 64)
			continue
		}
		if fields[0] == "Shmem:" {
			tmp, _ = strconv.ParseUint(fields[1], 10, 64)
			continue
		}
	}
	return
}

func SSD() (rr uint64, ww uint64, dt uint64, df uint64, du uint64, ok bool) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var u64 uint64
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		if fields[0] != "8" {
			continue
		}
		if regexp.MustCompile(`\d`).MatchString(fields[2]) {
			continue
		}
		u64, err = strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			continue
		}
		rr += u64 * 512
		u64, err = strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			continue
		}
		ww += u64 * 512
	}
	ssd, err := disk.Usage("/")
	if runtime.GOOS == "windows" {
		ssd, err = disk.Usage(`C:\`)
	}

	if err != nil {
		ok = false
		return
	}

	dt = ssd.Total
	du = ssd.Used
	df = ssd.Free
	ok = true

	/*
		if runtime.GOOS != "windows" {
			fs := syscall.Statfs_t{}
			syscall.Statfs("/", &fs)
			dt = fs.Blocks * uint64(fs.Bsize)
			df = fs.Bfree * uint64(fs.Bsize)
			if dt > 0 {
				du = (float64(dt - df)) * 100 / float64(dt)
			}
			mount, err := os.ReadFile("/proc/mounts")
			if err != nil {
				return
			}
			lines := strings.SplitSeq(string(mount), "\n")
			for line := range lines {
				fields := strings.Fields(line)
				if fields[1] == "/" && fields[3][:2] == "rw" {
					ok = true
					return
				}
			}
			ok = false
		}
	*/
	return
}

func NET() (ii uint64, oo uint64) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var u64 uint64
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 17 {
			continue
		}
		if fields[0] == "lo:" {
			continue
		}
		if strings.HasPrefix(fields[0], "docker") {
			continue
		}
		u64, err = strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		ii += u64
		u64, err = strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			continue
		}
		oo += u64
	}
	return
}

func Processes(n ...string) map[int]string {
	d, err := os.Open("/proc")
	if err != nil {
		return nil
	}
	defer d.Close()
	var res = make(map[int]string)
	for {
		names, err := d.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		for _, name := range names {
			if name[0] < '0' || name[0] > '9' {
				continue
			}
			pid, err := strconv.ParseInt(name, 10, 64)
			if err != nil {
				continue
			}
			exe, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
			if err != nil {
				continue
			}
			e := strings.TrimSpace(string(exe))
			if len(n) > 0 && e != n[0] {
				continue
			}
			res[int(pid)] = e
		}
	}
	return res
}

func Nginx() string {
	out, _ := exec.Command("systemctl", "is-active", "nginx").Output()
	return string(out)
}

func Redis() string {
	out, _ := exec.Command("systemctl", "is-active", "redis").Output()
	return string(out)
}

func Lang(c string) string {
	lo := Locale(c)
	arr := strings.Split(lo, "-")
	if len(arr) == 0 || arr[0] == "en" {
		return "en-US,en"
	} else {
		return arr[0] + ",en-US,en"
	}
}

func Acceptlang(c string) string {
	lo := Locale(c)
	arr := strings.Split(lo, "-")
	if len(arr) == 0 || arr[0] == "en" {
		return "en-US,en;q=0.9"
	} else {
		return arr[0] + ",en-US;q=0.9,en;q=0.8"
	}
}

func Locale(c string) string {
	if len(c) == 0 {
		return "en-US"
	}
	lang := make(map[string]string)
	lang["IN"] = "hi"
	lang["CH"] = "de"
	lang["NG"] = "en"
	lang["BA"] = "bs"
	lang["LU"] = "de"
	lang["BE"] = "nl"
	lang["MY"] = "en"
	lang["SG"] = "en"
	lang["ZA"] = "en"
	lang["NZ"] = "en"
	lang["CA"] = "en"
	lang["IE"] = "en"
	lang["US"] = "en"
	lang["GB"] = "en"
	lang["CN"] = "zh"
	lang["ET"] = "am"
	lang["EG"] = "ar"
	lang["DZ"] = "ar"
	lang["BH"] = "ar"
	lang["AE"] = "ar"
	lang["IQ"] = "ar"
	lang["JO"] = "ar"
	lang["KW"] = "ar"
	lang["LB"] = "ar"
	lang["LY"] = "ar"
	lang["MA"] = "ar"
	lang["OM"] = "ar"
	lang["QA"] = "ar"
	lang["SA"] = "ar"
	lang["SY"] = "ar"
	lang["TN"] = "ar"
	lang["YE"] = "ar"
	lang["CZ"] = "cs"
	lang["DK"] = "da"
	lang["LI"] = "de"
	lang["MV"] = "dv"
	lang["GR"] = "el"
	lang["AU"] = "en"
	lang["BZ"] = "en"
	lang["JM"] = "en"
	lang["PH"] = "en"
	lang["TT"] = "en"
	lang["ZW"] = "en"
	lang["AT"] = "de"
	lang["BY"] = "be"
	lang["BD"] = "bn"
	lang["AR"] = "es"
	lang["BO"] = "es"
	lang["CL"] = "es"
	lang["CO"] = "es"
	lang["CR"] = "es"
	lang["DO"] = "es"
	lang["EC"] = "es"
	lang["GT"] = "es"
	lang["HN"] = "es"
	lang["MX"] = "es"
	lang["NI"] = "es"
	lang["PA"] = "es"
	lang["PE"] = "es"
	lang["PR"] = "es"
	lang["PY"] = "es"
	lang["SV"] = "es"
	lang["UY"] = "es"
	lang["VE"] = "es"
	lang["GE"] = "ka"
	lang["EE"] = "et"
	lang["IR"] = "fa"
	lang["MC"] = "fr"
	lang["IL"] = "he"
	lang["AM"] = "hy"
	lang["JP"] = "ja"
	lang["KZ"] = "kk"
	lang["GL"] = "kl"
	lang["KH"] = "km"
	lang["KR"] = "ko"
	lang["KG"] = "ky"
	lang["LU"] = "lb"
	lang["LA"] = "lo"
	lang["BN"] = "ms"
	lang["NP"] = "ne"
	lang["AF"] = "ps"
	lang["BR"] = "pt"
	lang["LK"] = "si"
	lang["SI"] = "sl"
	lang["AL"] = "sq"
	lang["CS"] = "sr"
	lang["ME"] = "sr"
	lang["RS"] = "sr"
	lang["KE"] = "sw"
	lang["TJ"] = "tg"
	lang["TM"] = "tk"
	lang["UA"] = "uk"
	lang["PK"] = "ur"
	lang["VN"] = "vi"
	lang["SN"] = "wo"
	lang["HK"] = "zh"
	lang["MO"] = "zh"
	lang["TW"] = "zh"
	if val, ok := lang[c]; ok {
		return val + "-" + strings.ToUpper(c)
	}
	return strings.ToLower(c) + "-" + strings.ToUpper(c)
}
