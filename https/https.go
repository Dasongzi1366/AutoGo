package https

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Get 对地址url进行一次HTTP GET 请求。超时时间单位为毫秒
func Get(url string, timeout int) (code int, data []byte) {
	if timeout == 0 {
		timeout = 5000
	}
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	resp, err := client.Get(url)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

// PostMultipart 向目标地址发起类型为multipart/form-data的请求（通常用于文件上传等)
func PostMultipart(url string, fileName string, fileData []byte) (code int, data []byte) {
	// 创建一个缓冲区用于存储 multipart 数据
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// 创建一个 form 文件字段
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return 0, nil
	}

	// 将文件数据写入 part
	if _, err = io.Copy(part, bytes.NewReader(fileData)); err != nil {
		return 0, nil
	}

	// 关闭 multipart writer，设置结束标志
	_ = writer.Close()

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return 0, nil
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil
	}

	return resp.StatusCode, body
}
