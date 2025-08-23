package utl

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"math"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

func RndStr(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l]
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

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
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
