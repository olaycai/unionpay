package unionpay

import (
	"net/http"
	"fmt"
	"errors"
	"crypto"
	"crypto/sha256"
	"crypto/rsa"
	"crypto/rand"
	"net/url"
	"strings"
	"io/ioutil"
	"crypto/x509"
	"encoding/pem"
	"sort"
	"bytes"
	"encoding/base64"
)

type UnionPay struct {
	isProduction bool              // 是否会正式环境
	client       *http.Client      // 请求指针
	commonParam  map[string]string // 通用参数
	signCert     []byte            // 签名证书
	signPwd      string            // 签名证书密码
	signCertId   string            // 签名证书id
}

// New 构造函数
func New(merId string, privateKey []byte, signPwd string, signCertId string, isProduction bool) (client *UnionPay) {
	client = &UnionPay{}
	// 赋值
	client.isProduction = isProduction
	client.signCert = privateKey
	client.signPwd = signPwd
	client.signCertId = signCertId
	client.commonParam = map[string]string{
		"version":    DATA_VERSION, // 版本号
		"encoding":   ENCODEING,    // 编码方式
		"signMethod": SIGN_METHOD,  // 签名方法
		"merId":      merId,        // 商户代码
	}

	return client
}

// Sign 签名函数
func (this *UnionPay) Sign(param map[string]string) (signedParam url.Values, err error) {

	// 把map转换成url参数的形式
	var urlParam = url.Values{}
	for key, value := range param {
		urlParam.Add(key, value)
	}

	// 转换私钥格式
	var block *pem.Block
	block, _ = pem.Decode(this.signCert)
	if block == nil {
		return nil, errors.New("private key error")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	// 证书id
	param["certId"] = this.signCertId
	urlParam.Add("certId", this.signCertId)

	// 转换为string (此处不要使用Url.Values的Encode转换，因为银联签名时数据不能被urlencode)
	paramString := makeParams(param)
	// 使用完清空
	param = nil

	// 如果失败则直接返回空
	if err != nil {
		err = errors.New(fmt.Sprintf("get private key fail: %s", err))
		return
	}
	// sha256签名摘要(转2次)
	paramSha256 := sha256.Sum256([]byte (fmt.Sprintf("%x", sha256.Sum256([]byte (paramString)))))

	// 签名
	randReader := rand.Reader
	signer, err := rsa.SignPKCS1v15(randReader, privateKey.(*rsa.PrivateKey), crypto.SHA256, paramSha256[:])
	if err != nil {
		err = errors.New(fmt.Sprintf("sign fail: %s", err))
		return
	}

	urlParam.Add("signature", base64.StdEncoding.EncodeToString(signer))
	signedParam = urlParam
	return
}

// SendRequest 发起请求
func (this *UnionPay) SendRequest(method string, apiUrl string, params map[string]string) (result map[string]string, err error) {
	// 先签名
	urlParam, err := this.Sign(params)

	if err != nil {
		err = errors.New("sign fail")
		return
	}
	// 转换为io模式传输(发送前数据需要经过urlencode)
	ioReader := strings.NewReader(urlParam.Encode())
	request, err := http.NewRequest(method, apiUrl, ioReader)

	if err != nil {
		err = errors.New(fmt.Sprintf("create request fail: %s", err))
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	// 发起请求
	response, err := http.DefaultClient.Do(request)

	// 执行完后关闭流
	defer response.Body.Close()

	// 读取流数据
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("load response fail: %s", err))
		return
	}
	queryData, err := convertParam(string(data))
	// 结果转换为url参数模式方便读取
	if err != nil {
		return
	}
	result = queryData
	return
}

// 校验请求
func Validate(response map[string]string) (result bool, err error) {
	result = false
	// 获取接口返回中的公钥
	certContent := response["signPubKeyCert"]
	if certContent == "" {
		err = errors.New("cert empty")
		return
	}

	signature := response["signature"]
	delete(response, "signature")
	// 解析证书
	block, _ := pem.Decode([]byte(certContent))
	if block == nil {
		err = errors.New("block empty")
		return
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("fail")
		return
	}

	// 转换为字符串
	paramString := makeParams(response)
	// sha256签名摘要(转2次)
	paramSha256 := sha256.Sum256([]byte (fmt.Sprintf("%x", sha256.Sum256([]byte (paramString)))))
	decodeSignature, err := base64.StdEncoding.DecodeString(signature)

	// 验证签名
	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA256, paramSha256[:], decodeSignature)
	if err == nil {
		result = true
	}
	return
}

// 请求参数map转换为string模式
func makeParams(params map[string]string) (paramsString string) {
	var keys []string
	b := bytes.Buffer{}
	for k, _ := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, v := range keys {
		if b.Len() > 0 {
			b.WriteByte('&')
		}
		b.WriteString(v)
		b.WriteString("=")
		b.WriteString(params[v])
	}
	paramsString = b.String()
	return
}

// 转换param字符串为map
func convertParam(query string) (data map[string]string, err error) {
	queryData := map[string]string{}
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		queryData[key] = value
	}
	if len(queryData) <= 0 {
		err = errors.New("convert fail")
		return
	}
	data = queryData
	return
}