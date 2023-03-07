package aliyun

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/url"
	"strconv"
	"time"
)

var (
	signMethodMap = map[string]func() hash.Hash{
		"HMAC-SHA1":   sha1.New,
		"HMAC-SHA256": sha256.New,
		"HMAC-MD5":    md5.New,
	}
)

func hmacSign(signMethod string, httpMethod string, appKeySecret string, params url.Values) (signature []byte) {
	key := []byte(appKeySecret + "&")

	var h hash.Hash
	if method, ok := signMethodMap[signMethod]; ok {
		h = hmac.New(method, key)
	} else {
		h = hmac.New(sha1.New, key)
	}
	makeDataToSign(h, httpMethod, params)
	return h.Sum(nil)
}

func hmacSignToB64(signMethod string, httpMethod string, appKeySecret string, params url.Values) (signature string) {
	return base64.StdEncoding.EncodeToString(hmacSign(signMethod, httpMethod, appKeySecret, params))
}

type strToEnc struct {
	s string
	e bool
}

func makeDataToSign(w io.Writer, httpMethod string, vals url.Values) {
	in := make(chan *strToEnc)
	go func() {
		in <- &strToEnc{s: httpMethod}
		in <- &strToEnc{s: "&"}
		in <- &strToEnc{s: "/", e: true}
		in <- &strToEnc{s: "&"}
		in <- &strToEnc{s: vals.Encode(), e: true}
		close(in)
	}()

	specialUrlEncode(in, w)
}

var (
	encTilde = "%7E"         // '~' -> "%7E"
	encBlank = []byte("%20") // ' ' -> "%20"
	tilde    = []byte("~")
)

func specialUrlEncode(in <-chan *strToEnc, w io.Writer) {
	for s := range in {
		if !s.e {
			io.WriteString(w, s.s)
			continue
		}

		l := len(s.s)
		for i := 0; i < l; {
			ch := s.s[i]

			switch ch {
			case '%':
				if encTilde == s.s[i:i+3] {
					w.Write(tilde)
					i += 3
					continue
				}
				fallthrough
			case '*', '/', '&', '=':
				fmt.Fprintf(w, "%%%02X", ch)
			case '+':
				w.Write(encBlank)
			default:
				fmt.Fprintf(w, "%c", ch)
			}

			i += 1
		}
	}
}

// aliyunSigner 签名
func aliyunSigner(accessKeyID, accessSecret string, params *url.Values) {
	// 公共参数
	params.Set("Version", "2015-01-09")
	params.Set("Format", "JSON")
	params.Set("AccessKeyId", accessKeyID)
	params.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("SignatureNonce", strconv.FormatInt(time.Now().UnixNano(), 10))
	params.Set("SignatureMethod", "HMAC-SHA1")
	params.Set("SignatureVersion", "1.0")
	//签名
	params.Set("Signature", hmacSignToB64("HMAC-SHA1", "GET", accessSecret, *params))
}
