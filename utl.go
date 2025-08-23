package utl

import (
    "crypto/rand"
    "encoding/base64"
    "math"
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
