package frpc

/*
#include <stdlib.h>
const char* callJavaStringString(const char *name, const char *obj);
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

var config, mServerIp, mServerApiPort, mServerFrpConnPort string

func init() {
	downloadFile("frpc")
	shell("pkill -f frpc -c")
}

func Init(serverIp string, serverApiPort, serverFrpConnPort int) {
	mServerIp = serverIp
	mServerApiPort = strconv.Itoa(serverApiPort)
	mServerFrpConnPort = strconv.Itoa(serverFrpConnPort)
	config = fmt.Sprintf(`serverAddr = "%s"
serverPort = %s

`, mServerIp, mServerFrpConnPort)
}

// AddForward name建议使用 设备码:localPort 保证唯一性,转发成功返回转发后的服务器端口号
func AddForward(name string, localPort int) (int, error) {
	code, data := httpGet("http://"+mServerIp+":"+mServerApiPort+"/api/addport?name="+name, 0)
	if code != 200 {
		return 0, fmt.Errorf("连接FRPS服务器失败")
	}

	var apiResp ApiResponse
	err := json.Unmarshal(data, &apiResp)
	if err != nil {
		return 0, err
	}

	if apiResp.Code != 200 {
		return 0, fmt.Errorf(apiResp.Message)
	}

	config = config + fmt.Sprintf(`[[proxies]]
name = "%s"
type = "tcp"
localIP = "127.0.0.1"
localPort = %d
remotePort = %s

`, name, localPort, apiResp.Data)
	return s2i(apiResp.Data), nil
}

// Connect 阻塞方法,直到连接断开
func Connect() string {
	dir, _ := os.Getwd()
	_ = os.WriteFile(dir+"/frpc.toml", []byte(config), 0644)
	go func() {
		sleep(1000)
		_ = os.RemoveAll(dir + "/frpc.toml")
	}()
	return shell(dir + "/frpc -c ./frpc.toml")
}

// httpGet 对地址url进行一次HTTP GET 请求。超时时间单位为毫秒
func httpGet(url string, timeout int) (code int, data []byte) {
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

func downloadFile(name string) {
	dir, _ := os.Getwd()
	path := dir + "/" + name
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	fmt.Println("download " + name + " ..")
	for i := 0; i < 5; i++ {
		resp, err := http.Get("https://yunkong-1257133387.cos.ap-shanghai.myqcloud.com/AutoGo/" + runtime.GOARCH + "/" + name)
		if err != nil {
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			_ = resp.Body.Close()
			continue
		}
		_ = os.WriteFile(path, body, 0644)
		shell("chmod 755 " + path)
		fmt.Println("download " + name + " success")
		return
	}
	fmt.Println("download " + name + " failed")
	os.Exit(1)
}

func shell(cmd string) string {
	cStr1 := C.CString("shell")
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(cmd)
	defer C.free(unsafe.Pointer(cStr2))
	return C.GoString(C.callJavaStringString(cStr1, cStr2))
}

func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}
