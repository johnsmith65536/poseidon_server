package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"poseidon/entity"
	"poseidon/infra/mysql"
	"poseidon/utils"
	"sort"
	"strconv"
	"sync"
	"time"
)

var initOne sync.Once
var accessKey *entity.AccessKey

const (
	DurationSeconds = 3600
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

func GetSTSInfo(userId int64) (*STSInfo, error) {
	initOne.Do(func() {
		var err error
		accessKey, err = mysql.LoadSecretKey()
		if err != nil {
			panic(err)
		}
	})
	type param struct {
		Key   string
		Value string
	}
	var realResp struct {
		Credentials STSInfo `json:"Credentials"`
	}

	sessionName := strconv.FormatInt(userId, 10)
	params := []param{{Key: "AccessKeyId", Value: accessKey.AccessKeyId},
		{Key: "DurationSeconds", Value: strconv.FormatInt(DurationSeconds, 10)},
		{Key: "RoleArn", Value: accessKey.RoleArn},
		{Key: "Action", Value: "AssumeRole"},
		{Key: "Format", Value: "JSON"},
		{Key: "RoleSessionName", Value: sessionName},
		{Key: "SignatureMethod", Value: "HMAC-SHA1"},
		{Key: "SignatureNonce", Value: utils.GenerateUUID()},
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
	params = append(params, param{Key: "Signature", Value: base64.StdEncoding.EncodeToString(hMAC(stringToSign, accessKey.AccessKeySecret+"&"))})
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
