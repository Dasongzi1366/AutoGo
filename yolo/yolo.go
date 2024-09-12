package yolo

/*
#include <stdlib.h>
void callJavaLoadSo(const char* str);
long callJavaYoloV5Init(int threadNum, const char* labelPath, const char* paramPath, const char* binPath);
char* callJavaYoloV5Detect(long pointer,int size, int x1, int y1, int x2, int y2);
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
	"strings"
	"unsafe"
)

type Yolo struct {
	pointer int64
}

type tempObj struct {
	X      int     `json:"X"`
	Y      int     `json:"Y"`
	Width  int     `json:"宽"`
	Label  string  `json:"标签"`
	Prob   float32 `json:"精度"`
	Height int     `json:"高"`
}

type Obj struct {
	Label   string
	Prob    float32
	Left    int
	Right   int
	Top     int
	Bottom  int
	CenterX int
	CenterY int
	Width   int
	Height  int
}

func init() {
	loadSo("libc++_shared.so")
	loadSo("libyolov5.so")
	loadSo("libyolov5_native.so")
}

func New(cpuThreadNum int, labelPath, paramPath, binPath string) *Yolo {

	cLabelPath := C.CString(labelPath)
	defer C.free(unsafe.Pointer(cLabelPath))

	cParamPath := C.CString(paramPath)
	defer C.free(unsafe.Pointer(cParamPath))

	cBinPath := C.CString(binPath)
	defer C.free(unsafe.Pointer(cBinPath))

	i := int64(C.callJavaYoloV5Init(C.int(cpuThreadNum), cLabelPath, cParamPath, cBinPath))
	return &Yolo{pointer: i}
}

// Detect x2 y2为0时默认使用设备的宽高
func (y *Yolo) Detect(size, x1, y1, x2, y2 int) []*Obj {
	cStr := C.callJavaYoloV5Detect(C.long(y.pointer), C.int(size), C.int(x1), C.int(y1), C.int(x2), C.int(y2))
	if cStr != nil {
		defer C.free(unsafe.Pointer(cStr))
		str := C.GoString(cStr)
		var tempObjs []tempObj
		err := json.Unmarshal([]byte(str), &tempObjs)
		if err != nil {
			return nil
		}

		var objs []*Obj
		for _, temp := range tempObjs {
			obj := Obj{
				Label:   temp.Label,
				Prob:    temp.Prob,
				Left:    x1 + temp.X,
				Right:   x1 + temp.X + temp.Width,
				Top:     y1 + temp.Y,
				Bottom:  y1 + temp.Y + temp.Height,
				CenterX: x1 + temp.X + temp.Width/2,
				CenterY: y1 + temp.Y + temp.Height/2,
				Width:   temp.Width,
				Height:  temp.Height,
			}
			objs = append(objs, &obj)
		}
		return objs
	}
	return nil
}

func ToString(objs []*Obj) string {
	str := ""
	for _, obj := range objs {
		str = str + fmt.Sprintf("%+v", obj) + "\n"
	}
	if strings.HasSuffix(str, "\n") {
		str = str[:len(str)-1]
	}
	return str
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
