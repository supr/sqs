package sqs

import (
	"crypto/hmac"
	"encoding/base64"
	"launchpad.net/goamz/aws"
	"sort"
	"strings"
)

var b64 = base64.StdEncoding

func sign(auth aws.Auth, method, path string, params map[string]string, host string) {
	params["AWSAccessKeyId"] = auth.AccessKey
	params["SignatureVersion"] = "2"
	params["SignatureMethod"] = "HmacSHA256"

	var sarray []string
	for k, v := range params {
		sarray = append(sarray, aws.Encode(k)+"="+aws.Encode(v))
	}
	sort.StringSlice(sarray).Sort()
	joined := strings.Join(sarray, "&")
	payload := method + "\n" + host + "\n" + path + "\n" + joined
	hash := hmac.NewSHA256([]byte(auth.SecretKey))
	hash.Write([]byte(payload))
	signature := make([]byte, b64.EncodedLen(hash.Size()))
	b64.Encode(signature, hash.Sum())

	params["Signature"] = string(signature)
}
