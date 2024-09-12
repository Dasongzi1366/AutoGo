package keys

/*
#include <stdlib.h>
void callJavaKeyCode(int code);
const char* callJavaStringString(const char *name, const char *obj);
*/
import "C"
import "unsafe"

// Home 模拟按下Home键。
func Home() {
	keyCode(KEYCODE_HOME)
}

// Back 模拟按下返回键。
func Back() {
	keyCode(KEYCODE_BACK)
}

// Recents 显示最近任务。
func Recents() {
	keyCode(KEYCODE_APP_SWITCH)
}

// PowerDialog 弹出电源键菜单。
func PowerDialog() {
	shell("input keyevent --longpress KEYCODE_POWER")
}

// Notifications 拉出通知栏。
func Notifications() {
	keyCode(KEYCODE_NOTIFICATION)
}

// QuickSettings 显示快速设置(下拉通知栏到底)。
func QuickSettings() {
	shell("cmd statusbar expand-settings")
}

// VolumeUp 按下音量上键。
func VolumeUp() {
	keyCode(KEYCODE_VOLUME_UP)
}

// VolumeDown 按下音量下键。
func VolumeDown() {
	keyCode(KEYCODE_VOLUME_DOWN)
}

// Camera 模拟按下照相键。
func Camera() {
	keyCode(KEYCODE_CAMERA)
}

// Up 模拟按下物理按键上。
func Up() {
	keyCode(KEYCODE_DPAD_UP)
}

// Down 模拟按下物理按键下。
func Down() {
	keyCode(KEYCODE_DPAD_DOWN)
}

// Left 模拟按下物理按键左。
func Left() {
	keyCode(KEYCODE_DPAD_LEFT)
}

// Right 模拟按下物理按键右。
func Right() {
	keyCode(KEYCODE_DPAD_RIGHT)
}

// OK 模拟按下物理按键确定。
func OK() {
	keyCode(KEYCODE_DPAD_CENTER)
}

// Text 输入文字 只能为英文或英文符号。
func Text(txt string) {
	shell(`input text "` + txt + `"`)
}

// KeyCode 要按下的按键代码 具体参考KEYCODE_开头常量。
func KeyCode(code int) {
	keyCode(code)
}

func keyCode(k int) {
	C.callJavaKeyCode(C.int(k))
}

func shell(cmd string) string {
	cStr1 := C.CString("shell")
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(cmd)
	defer C.free(unsafe.Pointer(cStr2))
	return C.GoString(C.callJavaStringString(cStr1, cStr2))
}
