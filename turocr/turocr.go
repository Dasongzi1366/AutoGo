package turocr

/*
#include <stdlib.h>
int callJavaTurocrInit(const char* name,const char* libStr);
void callJavaLoadSo(const char* str);
const char* callJavaStringString(const char *name, const char *obj);
char* callJavaTurocrDetect(const char* name,int x1, int y1, int x2, int y2, int similar, int model, int interval, int minThreshold, int maxThreshold, int colorInverted, int slicRow, int slicColumn);
*/
import "C"
import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"unsafe"
)

type TurOcr struct {
	name string
}

func init() {
	loadSo("libtur.so")
}

// New 创建对象。
//
// 参数：
// - lib: 字库文件数据
// 返回值：
// - 成功返回对象指针失败返回nil
func New(lib []byte) *TurOcr {
	hash := md5.New()
	hash.Write(lib)
	md5Sum := hash.Sum(nil)
	name := hex.EncodeToString(md5Sum)

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cLibStr := C.CString(string(lib))
	defer C.free(unsafe.Pointer(cLibStr))

	i := int(C.callJavaTurocrInit(cName, cLibStr))
	if i > 0 {
		return &TurOcr{name: name}
	}
	return nil
}

// Detect 识别区域内文字。
//
// 参数：
// - x1, y1: 屏幕区域的左上角坐标
// - x2, y2: 屏幕区域的右下角坐标。为 0 时默认使用屏幕的最大宽高
// - similar: 匹配阈值。范围在 0 到 100 之间。值越大匹配越严格
// - model: 返回值的字符串格式。参考图灵。范围是0-6
// - interval: 间隔。不清楚就写20
// - minThreshold: 二值化起始位置
// - maxThreshold: 二值化结束位置
// - colorInverted: 颜色颠倒模式。不清楚就写-1
// - slicRow: 投影切割行。不清楚就写2
// - slicColumn: 投影切割列。不清楚就写1
// 返回值：
// - 返回符合条件的颜色像素数量，如果未找到符合条件的像素，则返回 0。
func (t *TurOcr) Detect(x1, y1, x2, y2, similar, model, interval, minThreshold, maxThreshold, colorInverted, slicRow, slicColumn int) string {
	cName := C.CString(t.name)
	defer C.free(unsafe.Pointer(cName))

	cStr := C.callJavaTurocrDetect(cName, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(similar), C.int(model), C.int(interval), C.int(minThreshold), C.int(maxThreshold), C.int(colorInverted), C.int(slicRow), C.int(slicColumn))
	if cStr != nil {
		defer C.free(unsafe.Pointer(cStr))
		return C.GoString(cStr)
	}
	return ""
}

func loadSo(name string) {
	downloadFile(name)
	dir, _ := os.Getwd()
	cStr := C.CString(dir + "/" + name)
	defer C.free(unsafe.Pointer(cStr))
	C.callJavaLoadSo(cStr)
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
