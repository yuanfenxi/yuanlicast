package yfxcast

import (
	"strings"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"crypto/md5"
	"io"
)


func VerifySignature(auth_key string,auth_timestamp string,auth_version string,body_md5 string,post_method string,uri string,auth_secret string) string {
	urls := strings.Split(uri,"?")
	string_to_sign := fmt.Sprintf("%s\n%s\nauth_key=%s&auth_timestamp=%s&auth_version=%s&body_md5=%s",
		post_method,
		urls[0],
		auth_key,
		auth_timestamp, auth_version,
		body_md5)
	//log.Println(fmt.Sprintf("string to hash is:%s,seret is :%s",string_to_sign,auth_secret))
	sig := hmac.New(sha256.New, []byte(auth_secret))
	sig.Write([]byte(string_to_sign))
	return hex.EncodeToString(sig.Sum(nil))
}


func getQueryString(auth_key string, auth_timestamp string, auth_version string, body_md5 string) string {
	return fmt.Sprintf("auth_key=%s&auth_timestamp=%s&auth_version=%s&body_md5=%s",
		auth_key,
		auth_timestamp,
		auth_version,
		body_md5)
}

func GetMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}
