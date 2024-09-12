package imgui

/*
#include <stdlib.h>
void callJavaImguiStatusBarInit(int x1, int y1, int x2, int y2, int r, int g, int b,int textSize);
void callJavaImguiStatusPrint(const char* str);
void callJavaImguiStatusClose();
void callJavaImguiDrawRect(int x1, int y1, int x2, int y2, int r, int g, int b);
void callJavaImguiToast(const char* str);
void callJavaImguiClearRect();
const char* callJavaStringString(const char *name, const char *obj);
*/
import "C"
import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"unsafe"
)

// TextItem 表示颜色和文本的组合
type TextItem struct {
	TextColor color.Color //文字颜色 如果为 nil，则默认使用白色 (255,255,255)
	Text      string
}

func init() {
	downloadFile("imgui")
}

// StatusBarInit 转态条初始化
// x1, y1, x2, y2 是状态条的坐标
// bgColor 是状态条的背景颜色，如果为 nil，则默认使用灰色 (100, 100, 100)
// textSize 是状态条上的文字大小，如果小于等于 0，则默认使用 45
func StatusBarInit(x1, y1, x2, y2 int, bgColor color.Color, textSize int) {
	var rInt, gInt, bInt int
	if bgColor == nil {
		rInt, gInt, bInt = 100, 100, 100
	} else {
		r, g, b, _ := bgColor.RGBA()
		rInt = int(r >> 8)
		gInt = int(g >> 8)
		bInt = int(b >> 8)
	}
	if textSize <= 0 {
		textSize = 45
	}
	C.callJavaImguiStatusBarInit(C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(rInt), C.int(gInt), C.int(bInt), C.int(textSize))
}

// StatusBarPrint 状态条打印信息
func StatusBarPrint(items []TextItem) {
	str := ""
	for _, item := range items {
		if item.Text != "" {
			var rInt, gInt, bInt string
			if item.TextColor == nil {
				rInt, gInt, bInt = "255", "255", "255"
			} else {
				r, g, b, _ := item.TextColor.RGBA()
				rInt = i2s(int(r >> 8))
				gInt = i2s(int(g >> 8))
				bInt = i2s(int(b >> 8))
			}
			str = str + rInt + "|" + gInt + "|" + bInt + "|" + item.Text + "|"
		}
	}
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	C.callJavaImguiStatusPrint(cStr)
}

// StatusBarClose 状态条销毁
func StatusBarClose() {
	C.callJavaImguiStatusClose()
}

// DrawRect 绘制矩形
func DrawRect(x1, y1, x2, y2 int, c color.Color) {
	r, g, b, _ := c.RGBA()
	rInt := int(r >> 8)
	gInt := int(g >> 8)
	bInt := int(b >> 8)
	C.callJavaImguiDrawRect(C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(rInt), C.int(gInt), C.int(bInt))
}

// ClearRect 清除绘制的矩形
func ClearRect() {
	C.callJavaImguiClearRect()
}

// Toast 显示Toast提示信息
func Toast(message string) {
	cStr := C.CString(message)
	defer C.free(unsafe.Pointer(cStr))
	C.callJavaImguiToast(cStr)
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

func i2s(i int) string {
	return strconv.Itoa(i)
}
