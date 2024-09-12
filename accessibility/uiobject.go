package accessibility

/*
#include <stdlib.h>
int callJavaStringBool(const char *name, const char *obj);
int callJavaStringInt(const char *name, const char *obj);
const char* callJavaStringString(const char *name, const char *obj);
const char* callJavaGetParent(const char *obj);
const char* callJavaGetChild(const char *obj,int index);
const char* callJavaGetChildren(const char *obj);
int callJavaSetSelection(const char *obj,int start, int end);
int callJavaSetText(const char *obj, const char *txt);
*/
import "C"
import (
	"strconv"
	"strings"
	"unsafe"
)

type UiObject struct {
	objStr string
}

type Rect struct {
	Left    int
	Right   int
	Top     int
	Bottom  int
	CenterX int
	CenterY int
	Width   int
	Height  int
}

// Click 点击该控件，并返回是否点击成功
func (u *UiObject) Click() bool {
	return callJavaStringBool("uiObjectClick", u.objStr)
}

// ClickCenter 使用坐标点击该控件的中点，相当于click(uiObj.bounds().centerX(), uiObject.bounds().centerY())
func (u *UiObject) ClickCenter() bool {
	return callJavaStringBool("uiObjectClickCenter", u.objStr)
}

// ClickLongClick 长按该控件，并返回是否点击成功
func (u *UiObject) ClickLongClick() bool {
	return callJavaStringBool("uiObjectLongClick", u.objStr)
}

// Copy 对输入框文本的选中内容进行复制，并返回是否操作成功
func (u *UiObject) Copy() bool {
	return callJavaStringBool("uiObjectCopy", u.objStr)
}

// Cut 对输入框文本的选中内容进行剪切，并返回是否操作成功
func (u *UiObject) Cut() bool {
	return callJavaStringBool("uiObjectCut", u.objStr)
}

// Paste 对输入框控件进行粘贴操作，把剪贴板内容粘贴到输入框中，并返回是否操作成功
func (u *UiObject) Paste() bool {
	return callJavaStringBool("uiObjectPaste", u.objStr)
}

// ScrollForward 对控件执行向前滑动的操作，并返回是否操作成功
func (u *UiObject) ScrollForward() bool {
	return callJavaStringBool("uiObjectScrollForward", u.objStr)
}

// ScrollBackward 对控件执行向后滑动的操作，并返回是否操作成功
func (u *UiObject) ScrollBackward() bool {
	return callJavaStringBool("uiObjectScrollBackward", u.objStr)
}

// Collapse 对控件执行折叠操作，并返回是否操作成功
func (u *UiObject) Collapse() bool {
	return callJavaStringBool("uiObjectCollapse", u.objStr)
}

// Expand 对控件执行展开操作，并返回是否操作成功
func (u *UiObject) Expand() bool {
	return callJavaStringBool("uiObjectExpand", u.objStr)
}

// Show 执行显示操作，并返回是否操作成功
func (u *UiObject) Show() bool {
	return callJavaStringBool("uiObjectShow", u.objStr)
}

// Select 对控件执行"选中"操作，并返回是否操作成功
func (u *UiObject) Select() bool {
	return callJavaStringBool("uiObjectSelect", u.objStr)
}

// ClearSelect 清除控件的选中状态，并返回是否操作成功
func (u *UiObject) ClearSelect() bool {
	return callJavaStringBool("uiObjectClearSelect", u.objStr)
}

// SetSelection 对输入框控件设置选中的文字内容，并返回是否操作成功
func (u *UiObject) SetSelection(start, end int) bool {
	cStr := C.CString(u.objStr)
	defer C.free(unsafe.Pointer(cStr))
	return int(C.callJavaSetSelection(cStr, C.int(start), C.int(end))) == 1
}

// SetText 设置输入框控件的文本内容，并返回是否设置成功
func (u *UiObject) SetText(str string) bool {
	cStr1 := C.CString(u.objStr)
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(str)
	defer C.free(unsafe.Pointer(cStr2))
	return int(C.callJavaSetText(cStr1, cStr2)) == 1
}

// GetClickAble 获取控件的 clickable 属性
func (u *UiObject) GetClickAble() bool {
	return callJavaStringBool("uiObjectGetClickAble", u.objStr)
}

// GetLongClickAble 获取控件的 longClickable 属性
func (u *UiObject) GetLongClickAble() bool {
	return callJavaStringBool("uiObjectGetLongClickAble", u.objStr)
}

// GetCheckable 获取控件的 checkable 属性
func (u *UiObject) GetCheckable() bool {
	return callJavaStringBool("uiObjectGetCheckable", u.objStr)
}

// GetSelected 获取控件的 selected 属性
func (u *UiObject) GetSelected() bool {
	return callJavaStringBool("uiObjectGetSelected", u.objStr)
}

// GetEnabled 获取控件的 enabled 属性
func (u *UiObject) GetEnabled() bool {
	return callJavaStringBool("uiObjectGetEnabled", u.objStr)
}

// GetScrollAble 获取控件的 scrollable 属性
func (u *UiObject) GetScrollAble() bool {
	return callJavaStringBool("uiObjectGetScrollAble", u.objStr)
}

// GetEditable 获取控件的 editable 属性
func (u *UiObject) GetEditable() bool {
	return callJavaStringBool("uiObjectGetEditable", u.objStr)
}

// GetMultiLine 获取控件的 multiLine 属性
func (u *UiObject) GetMultiLine() bool {
	return callJavaStringBool("uiObjectGetMultiLine", u.objStr)
}

// GetChecked 获取控件的 checked 属性
func (u *UiObject) GetChecked() bool {
	return callJavaStringBool("uiObjectGetChecked", u.objStr)
}

// GetFocusable 获取控件的 focusable 属性
func (u *UiObject) GetFocusable() bool {
	return callJavaStringBool("uiObjectGetFocusable", u.objStr)
}

// GetDismissable 获取控件的 dismissable 属性
func (u *UiObject) GetDismissable() bool {
	return callJavaStringBool("uiObjectGetDismissable", u.objStr)
}

// GetContextClickable 获取控件的 contextClickable 属性
func (u *UiObject) GetContextClickable() bool {
	return callJavaStringBool("uiObjectGetContextClickable", u.objStr)
}

// GetAccessibilityFocused 获取控件的 accessibilityFocused 属性
func (u *UiObject) GetAccessibilityFocused() bool {
	return callJavaStringBool("uiObjectGetAccessibilityFocused", u.objStr)
}

// GetChildCount 获取控件的子控件数目
func (u *UiObject) GetChildCount() int {
	return callJavaStringInt("uiObjectGetChildCount", u.objStr)
}

// GetDrawingOrder 获取控件在父控件中的绘制次序
func (u *UiObject) GetDrawingOrder() int {
	return callJavaStringInt("uiObjectGetDrawingOrder", u.objStr)
}

// GetIndexInParent 获取控件在父控件中的索引
func (u *UiObject) GetIndexInParent() Rect {
	bounds := callJavaStringString("uiObjectGetIndexInParent", u.objStr)
	arr := strings.Split(bounds, ",")
	if len(arr) == 4 {
		return Rect{Left: s2i(arr[0]), Top: s2i(arr[1]), Right: s2i(arr[2]), Bottom: s2i(arr[3]), Width: s2i(arr[2]) - s2i(arr[0]), Height: s2i(arr[3]) - s2i(arr[1]), CenterX: s2i(arr[0]) + ((s2i(arr[2]) - s2i(arr[0])) / 2), CenterY: s2i(arr[1]) + ((s2i(arr[3]) - s2i(arr[1])) / 2)}
	}
	return Rect{}
}

// GetBounds 获取控件在屏幕上的范围
func (u *UiObject) GetBounds() Rect {
	bounds := callJavaStringString("uiObjectGetBounds", u.objStr)
	arr := strings.Split(bounds, ",")
	if len(arr) == 4 {
		return Rect{Left: s2i(arr[0]), Top: s2i(arr[1]), Right: s2i(arr[2]), Bottom: s2i(arr[3]), Width: s2i(arr[2]) - s2i(arr[0]), Height: s2i(arr[3]) - s2i(arr[1]), CenterX: s2i(arr[0]) + ((s2i(arr[2]) - s2i(arr[0])) / 2), CenterY: s2i(arr[1]) + ((s2i(arr[3]) - s2i(arr[1])) / 2)}
	}
	return Rect{}
}

// GetBoundsInParent 获取控件在父控件中的范围
func (u *UiObject) GetBoundsInParent() string {
	return callJavaStringString("uiObjectGetBoundsInParent", u.objStr)
}

// GetId 获取控件的ID
func (u *UiObject) GetId() string {
	return callJavaStringString("uiObjectGetId", u.objStr)
}

// GetText 获取控件的文本内容
func (u *UiObject) GetText() string {
	return callJavaStringString("uiObjectGetText", u.objStr)
}

// GetDesc 获取控件的描述内容
func (u *UiObject) GetDesc() string {
	return callJavaStringString("uiObjectGetDesc", u.objStr)
}

// GetPackageName 获取控件的包名
func (u *UiObject) GetPackageName() string {
	return callJavaStringString("uiObjectGetPackageName", u.objStr)
}

// GetClassName 获取控件的类名
func (u *UiObject) GetClassName() string {
	return callJavaStringString("uiObjectGetClassName", u.objStr)
}

// GetParent 获取控件的父控件
func (u *UiObject) GetParent() *UiObject {
	cStr := C.CString(u.objStr)
	defer C.free(unsafe.Pointer(cStr))
	str := C.GoString(C.callJavaGetParent(cStr))
	if str == "" {
		return nil
	}
	return &UiObject{objStr: str}
}

// GetChild 获取控件的指定索引的子控件
func (u *UiObject) GetChild(index int) *UiObject {
	cStr := C.CString(u.objStr)
	defer C.free(unsafe.Pointer(cStr))
	str := C.GoString(C.callJavaGetChild(cStr, C.int(index)))
	if str == "" {
		return nil
	}
	return &UiObject{objStr: str}
}

// GetChildren 获取控件的所有子控件
func (u *UiObject) GetChildren() []*UiObject {
	cStr := C.CString(u.objStr)
	defer C.free(unsafe.Pointer(cStr))
	str := C.GoString(C.callJavaGetChildren(cStr))
	if str == "" {
		return nil
	}
	arr := strings.Split(str, "\n")
	var uiObjectArray []*UiObject
	for _, s := range arr {
		if s != "" {
			uiObjectArray = append(uiObjectArray, &UiObject{objStr: s})
		}
	}
	return uiObjectArray
}

// ToString 节点对象转文本
func (u *UiObject) ToString() string {
	return callJavaStringString("uiObjectToString", u.objStr)
}

func callJavaStringBool(name, value string) bool {
	cStr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(value)
	defer C.free(unsafe.Pointer(cStr2))
	return int(C.callJavaStringBool(cStr1, cStr2)) == 1
}

func callJavaStringInt(name, value string) int {
	cStr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(value)
	defer C.free(unsafe.Pointer(cStr2))
	return int(C.callJavaStringInt(cStr1, cStr2))
}

func callJavaStringString(name, value string) string {
	cStr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(value)
	defer C.free(unsafe.Pointer(cStr2))
	return C.GoString(C.callJavaStringString(cStr1, cStr2))
}

func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}
