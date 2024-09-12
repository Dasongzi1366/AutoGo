package media

/*
#include <stdlib.h>
const char* callJavaStringString(const char *name, const char *obj);
*/
import "C"
import "unsafe"

// ScanFile 扫描路径path的媒体文件，将它加入媒体库中
func ScanFile(path string) {
	shell("am broadcast -a android.intent.action.MEDIA_SCANNER_SCAN_FILE -d file://" + path)
}

func shell(cmd string) string {
	cStr1 := C.CString("shell")
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(cmd)
	defer C.free(unsafe.Pointer(cStr2))
	return C.GoString(C.callJavaStringString(cStr1, cStr2))
}
