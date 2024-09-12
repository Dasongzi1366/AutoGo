package ime

/*
#include <stdlib.h>
const char* callJavaGetClipText();
int callJavaSetClipText(const char* str);
*/
import "C"
import "unsafe"

// GetClipText 获取剪切板内容
func GetClipText() string {
	return C.GoString(C.callJavaGetClipText())
}

// SetClipText 设置剪切板内容
func SetClipText(text string) bool {
	cStr := C.CString(text)
	defer C.free(unsafe.Pointer(cStr))
	return int(C.callJavaSetClipText(cStr)) == 1
}
