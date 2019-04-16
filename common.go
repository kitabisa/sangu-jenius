package jenius

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func (gateway *CoreGateway) generateBtpnSignature(httpMethod string, relativeUrl string, btpnTimestamp string, requestBody string) string {
	var stringToSignature string
	if requestBody != "" {
		stringToSignature = btpnCleanString(fmt.Sprintf("%v:%v:%v:%v:%v", httpMethod, relativeUrl, gateway.Client.JeniusApiKey, btpnTimestamp, requestBody))
	} else {
		stringToSignature = btpnCleanString(fmt.Sprintf("%v:%v:%v:%v", httpMethod, relativeUrl, gateway.Client.JeniusApiKey, btpnTimestamp))
	}

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(gateway.Client.JeniusApiSecret))

	// Write Data to it
	h.Write([]byte(stringToSignature))

	// Get result and encode as hexadecimal string
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return sha
}

func (gateway *CoreGateway) getBtpnTimestamp() string {
	// 'YYYY-MM-DDTHH24:MI:SS.ZZZ+07:00'
	jakarta := time.FixedZone("Asia/Jakarta", 7*60*60)
	t := time.Now().In(jakarta)
	return t.Format("2006-01-02T15:04:05.000-07:00")
}

func btpnCleanString(any string) string {
	s := strings.ReplaceAll(any, " ", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

func btpnConvertTimestamp(req uint) string {
	timeReq := time.Unix(int64(req), 0)
	return timeReq.Format("2006-01-02T15:04:05.000-07:00")
}