package esign

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

const fileChunk = 8192 // 8KB

// Base64Encode Base64编码
func Base64Encode(dataString string) string {
	encodeString := base64.StdEncoding.EncodeToString([]byte(dataString))
	return encodeString
}

// Base64Decode Base64解码
func Base64Decode(encodeString string) []byte {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		fmt.Println(err)
	}
	return decodeBytes
}

// Base64EncodeByFile 将文件进行Base64编码
func Base64EncodeByFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	} else {
		file.Close()
	}
	encodeString := base64.StdEncoding.EncodeToString(fileBytes)
	return encodeString
}

// SaveFileByBase64 保存文件
func SaveFileByBase64(base64String, outFilePath string) {
	// 将Base64字符串解码为[]byte
	var fileBytes = Base64Decode(base64String)

	saveFileErr := ioutil.WriteFile(outFilePath, fileBytes, 0666)

	if saveFileErr != nil {
		fmt.Println("文件保存失败:" + saveFileErr.Error())
		panic(saveFileErr)
	} else {
		fmt.Println("文件保存成功:" + outFilePath)
	}
}

// DoHashMd5 摘要md5
func DoHashMd5(body string) (md5Str string) {
	hash := md5.New()
	hash.Write([]byte(body))
	md5Data := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(md5Data)
}

// CountFileMd5 计算文件content-md5
func CountFileMd5(filePath string) (string, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		return err.Error(), 0
	}
	defer file.Close()

	info, _ := file.Stat()
	fileSize := info.Size()

	blocks := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		buf := make([]byte, blockSize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), fileSize
}

// DoSignatureBase64 sha256摘要签名
func DoSignatureBase64(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	buf := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(buf)
}

// AppendSignDataString 拼接请求参数
func AppendSignDataString(method string, accept string, contentMD5 string, contentType string, date string, headers string, url string) string {
	var buffer bytes.Buffer
	buffer.WriteString(method)
	buffer.WriteString("\n")
	buffer.WriteString(accept)
	buffer.WriteString("\n")
	buffer.WriteString(contentMD5)
	buffer.WriteString("\n")
	buffer.WriteString(contentType)
	buffer.WriteString("\n")
	buffer.WriteString(date)
	buffer.WriteString("\n")
	if len(headers) == 0 {
		buffer.WriteString(headers)
		buffer.WriteString(url)
	} else {
		buffer.WriteString(headers)
		buffer.WriteString("\n")
		buffer.WriteString(url)
	}
	return buffer.String()
}

// ByteToJson byte转json
func ByteToJson(initResult []byte) map[string]interface{} {
	var initResultJson interface{}
	_ = json.Unmarshal(initResult, &initResultJson)
	jsonMap, err := initResultJson.(map[string]interface{})
	if !err {
		fmt.Println("DO SOMETHING!")
		return nil
	}
	return jsonMap
}

// SendHttp 发送HTTP请求
func sendHttp(apiUrl string, data string, method string, headers map[string]string) ([]byte, int) {
	// API接口返回值
	var apiResult []byte
	url := apiUrl
	var jsonStr = []byte(data)
	var req *http.Request
	var err error
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	} else {
		var httpStatus = resp.StatusCode
		if httpStatus != http.StatusOK {
			return apiResult, httpStatus
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		apiResult = body
		return apiResult, httpStatus
	}
}

func SendCommHttp(apiUrl string, data interface{}, method string) (initResult []byte, httpStatus int) {
	dataJson, _ := json.Marshal(data)
	dataJsonStr := string(dataJson)
	httpUrl := Config().Host() + apiUrl
	var md5Str string
	md5Str = DoHashMd5(dataJsonStr)
	message := AppendSignDataString("POST", "*/*", md5Str, "application/json; charset=UTF-8", "", "", apiUrl)
	reqSignature := DoSignatureBase64(message, Config().ProjectSecret())
	// 初始化接口返回值
	initResult, httpStatus = sendHttp(httpUrl, dataJsonStr, method, buildCommHeader(md5Str, reqSignature))
	return initResult, httpStatus
}

func buildCommHeader(contentMD5 string, reqSignature string) (header map[string]string) {
	headers := map[string]string{}
	headers["X-Tsign-Open-App-Id"] = Config().ProjectId()
	headers["X-Tsign-Open-Ca-Timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	headers["Accept"] = "*/*"
	headers["X-Tsign-Open-Ca-Signature"] = reqSignature
	headers["Content-MD5"] = contentMD5
	headers["Content-Type"] = "application/json; charset=UTF-8"
	headers["X-Tsign-Open-Auth-Mode"] = "Signature"
	return headers
}
