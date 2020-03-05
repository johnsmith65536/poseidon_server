package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	AccessKeyId     = "LTAI4Fg2xDMnaeK4mvBu8gwa"
	AccessKeySecret = "8FIMffgGhCVw81LUX1h6ffrWHAPLSn"
	DurationSeconds = 3600
	RoleArn         = "acs:ram::1387747617960990:role/poseidon-data-sts-rw"
)

type STSInfo struct {
	SecurityToken   string `json:"SecurityToken"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
}

func hMAC(content string, key string) []byte {
	//hmac ,use sha1
	keyByte := []byte(key)
	mac := hmac.New(sha1.New, keyByte)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func generateUUID() string {
	var characterSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var res string
	for i := 1; i <= 16; i++ {
		var x = rand.Intn(62)
		res += characterSet[x : x+1]
		if i%4 == 0 && i != 16 {
			res += "-"
		}
	}
	return res
}

func GetSTSInfo(userId int64) (*STSInfo, error) {
	type param struct {
		Key   string
		Value string
	}
	var realResp struct {
		Credentials STSInfo `json:"Credentials"`
	}

	sessionName := strconv.FormatInt(userId, 10)
	params := []param{{Key: "AccessKeyId", Value: AccessKeyId},
		{Key: "DurationSeconds", Value: strconv.FormatInt(DurationSeconds, 10)},
		{Key: "RoleArn", Value: RoleArn},
		{Key: "Action", Value: "AssumeRole"},
		{Key: "Format", Value: "JSON"},
		{Key: "RoleSessionName", Value: sessionName},
		{Key: "SignatureMethod", Value: "HMAC-SHA1"},
		{Key: "SignatureNonce", Value: generateUUID()},
		{Key: "SignatureVersion", Value: "1.0"},
		{Key: "Timestamp", Value: time.Now().UTC().Format("2006-01-02T15:04:05Z")},
		{Key: "Version", Value: "2015-04-01"},
	}
	sort.Slice(params, func(i, j int) bool {
		return params[i].Key < params[j].Key
	})
	var canonicalizedQueryString string
	for i := 0; i < len(params); i++ {
		canonicalizedQueryString += params[i].Key + "=" + url.QueryEscape(params[i].Value)
		if i != len(params)-1 {
			canonicalizedQueryString += "&"
		}
	}
	var stringToSign = "GET" + "&" + url.QueryEscape("/") + "&" + url.QueryEscape(canonicalizedQueryString)
	params = append(params, param{Key: "Signature", Value: base64.StdEncoding.EncodeToString(hMAC(stringToSign, AccessKeySecret+"&"))})
	httpParams := url.Values{}
	Url, err := url.Parse("https://sts.aliyuncs.com")
	if err != nil {
		return nil, err
	}
	for _, param := range params {
		httpParams.Set(param.Key, param.Value)
	}
	Url.RawQuery = httpParams.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &realResp)
	if err != nil {
		return nil, err
	}
	return &realResp.Credentials, nil
}
