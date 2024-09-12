package accessibility

/*
#include <stdlib.h>
int callJavaNewAccessibility();
const char* callJavaFindOnce666(int obj, const char *selector);
const char* callJavaFind(int obj, const char *selector);
int callJavaUiObjectClick(const char *obj);
*/
import "C"
import (
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Accessibility struct {
	index    int
	selector string
}

// New 创建一个新的 Accessibility 对象
func New() *Accessibility {
	return &Accessibility{index: int(C.callJavaNewAccessibility())}
}

// Text 设置选择器的 text 属性
func (a *Accessibility) Text(value string) *Accessibility {
	a.selector = a.selector + "text||" + value + "\n"
	return a
}

// TextContains 设置选择器的 textContains 属性，用于匹配包含指定文本的控件
func (a *Accessibility) TextContains(value string) *Accessibility {
	a.selector = a.selector + "textContains||" + value + "\n"
	return a
}

// TextStartsWith 设置选择器的 textStartsWith 属性，用于匹配以指定文本开头的控件
func (a *Accessibility) TextStartsWith(value string) *Accessibility {
	a.selector = a.selector + "textStartsWith||" + value + "\n"
	return a
}

// TextEndsWith 设置选择器的 textEndsWith 属性，用于匹配以指定文本结尾的控件
func (a *Accessibility) TextEndsWith(value string) *Accessibility {
	a.selector = a.selector + "textEndsWith||" + value + "\n"
	return a
}

// TextMatches 设置选择器的 textMatches 属性，用于匹配符合指定正则表达式的控件
func (a *Accessibility) TextMatches(value string) *Accessibility {
	a.selector = a.selector + "textMatches||" + value + "\n"
	return a
}

// Desc 设置选择器的 desc 属性，用于匹配描述等于指定文本的控件
func (a *Accessibility) Desc(value string) *Accessibility {
	a.selector = a.selector + "desc||" + value + "\n"
	return a
}

// DescContains 设置选择器的 descContains 属性，用于匹配描述包含指定文本的控件
func (a *Accessibility) DescContains(value string) *Accessibility {
	a.selector = a.selector + "descContains||" + value + "\n"
	return a
}

// DescStartsWith 设置选择器的 descStartsWith 属性，用于匹配描述以指定文本开头的控件
func (a *Accessibility) DescStartsWith(value string) *Accessibility {
	a.selector = a.selector + "descStartsWith||" + value + "\n"
	return a
}

// DescEndsWith 设置选择器的 descEndsWith 属性，用于匹配描述以指定文本结尾的控件
func (a *Accessibility) DescEndsWith(value string) *Accessibility {
	a.selector = a.selector + "descEndsWith||" + value + "\n"
	return a
}

// DescMatches 设置选择器的 descMatches 属性，用于匹配描述符合指定正则表达式的控件
func (a *Accessibility) DescMatches(value string) *Accessibility {
	a.selector = a.selector + "descMatches||" + value + "\n"
	return a
}

// Id 设置选择器的 id 属性，用于匹配ID等于指定值的控件
func (a *Accessibility) Id(value string) *Accessibility {
	a.selector = a.selector + "id||" + value + "\n"
	return a
}

// IdContains 设置选择器的 idContains 属性，用于匹配ID包含指定值的控件
func (a *Accessibility) IdContains(value string) *Accessibility {
	a.selector = a.selector + "idContains||" + value + "\n"
	return a
}

// IdStartsWith 设置选择器的 idStartsWith 属性，用于匹配ID以指定值开头的控件
func (a *Accessibility) IdStartsWith(value string) *Accessibility {
	a.selector = a.selector + "idStartsWith||" + value + "\n"
	return a
}

// IdEndsWith 设置选择器的 idEndsWith 属性，用于匹配ID以指定值结尾的控件
func (a *Accessibility) IdEndsWith(value string) *Accessibility {
	a.selector = a.selector + "idEndsWith||" + value + "\n"
	return a
}

// IdMatches 设置选择器的 idMatches 属性，用于匹配ID符合指定正则表达式的控件
func (a *Accessibility) IdMatches(value string) *Accessibility {
	a.selector = a.selector + "idMatches||" + value + "\n"
	return a
}

// ClassName 设置选择器的 className 属性，用于匹配类名等于指定值的控件
func (a *Accessibility) ClassName(value string) *Accessibility {
	a.selector = a.selector + "className||" + value + "\n"
	return a
}

// ClassNameContains 设置选择器的 classNameContains 属性，用于匹配类名包含指定值的控件
func (a *Accessibility) ClassNameContains(value string) *Accessibility {
	a.selector = a.selector + "classNameContains||" + value + "\n"
	return a
}

// ClassNameStartsWith 设置选择器的 classNameStartsWith 属性，用于匹配类名以指定值开头的控件
func (a *Accessibility) ClassNameStartsWith(value string) *Accessibility {
	a.selector = a.selector + "classNameStartsWith||" + value + "\n"
	return a
}

// ClassNameEndsWith 设置选择器的 classNameEndsWith 属性，用于匹配类名以指定值结尾的控件
func (a *Accessibility) ClassNameEndsWith(value string) *Accessibility {
	a.selector = a.selector + "classNameEndsWith||" + value + "\n"
	return a
}

// ClassNameMatches 设置选择器的 classNameMatches 属性，用于匹配类名符合指定正则表达式的控件
func (a *Accessibility) ClassNameMatches(value string) *Accessibility {
	a.selector = a.selector + "classNameMatches||" + value + "\n"
	return a
}

// PackageName 设置选择器的 packageName 属性，用于匹配包名等于指定值的控件
func (a *Accessibility) PackageName(value string) *Accessibility {
	a.selector = a.selector + "packageName||" + value + "\n"
	return a
}

// PackageNameContains 设置选择器的 packageNameContains 属性，用于匹配包名包含指定值的控件
func (a *Accessibility) PackageNameContains(value string) *Accessibility {
	a.selector = a.selector + "packageNameContains||" + value + "\n"
	return a
}

// PackageNameStartsWith 设置选择器的 packageNameStartsWith 属性，用于匹配包名以指定值开头的控件
func (a *Accessibility) PackageNameStartsWith(value string) *Accessibility {
	a.selector = a.selector + "packageNameStartsWith||" + value + "\n"
	return a
}

// PackageNameEndsWith 设置选择器的 packageNameEndsWith 属性，用于匹配包名以指定值结尾的控件
func (a *Accessibility) PackageNameEndsWith(value string) *Accessibility {
	a.selector = a.selector + "packageNameEndsWith||" + value + "\n"
	return a
}

// PackageNameMatches 设置选择器的 packageNameMatches 属性，用于匹配包名符合指定正则表达式的控件
func (a *Accessibility) PackageNameMatches(value string) *Accessibility {
	a.selector = a.selector + "packageNameMatches||" + value + "\n"
	return a
}

// Bounds 设置选择器的 bounds 属性，用于匹配控件在屏幕上的范围
func (a *Accessibility) Bounds(left, top, right, bottom int) *Accessibility {
	a.selector = a.selector + "bounds||" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "\n"
	return a
}

// BoundsInside 设置选择器的 boundsInside 属性，用于匹配控件在屏幕内的范围
func (a *Accessibility) BoundsInside(left, top, right, bottom int) *Accessibility {
	a.selector = a.selector + "boundsInside||" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "\n"
	return a
}

// BoundsContains 设置选择器的 boundsContains 属性，用于匹配控件包含在指定范围内
func (a *Accessibility) BoundsContains(left, top, right, bottom int) *Accessibility {
	a.selector = a.selector + "boundsContains||" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "\n"
	return a
}

// DrawingOrder 设置选择器的 drawingOrder 属性，用于匹配控件在父控件中的绘制顺序
func (a *Accessibility) DrawingOrder(value int) *Accessibility {
	a.selector = a.selector + "drawingOrder||" + i2s(value) + "\n"
	return a
}

// ClickAble 设置选择器的 clickAble 属性，用于匹配控件是否可点击
func (a *Accessibility) ClickAble(value bool) *Accessibility {
	a.selector = a.selector + "clickAble||" + b2s(value) + "\n"
	return a
}

// LongClickAble 设置选择器的 longClickAble 属性，用于匹配控件是否可长按
func (a *Accessibility) LongClickAble(value bool) *Accessibility {
	a.selector = a.selector + "longClickAble||" + b2s(value) + "\n"
	return a
}

// CheckAble 设置选择器的 checkAble 属性，用于匹配控件是否可选中
func (a *Accessibility) CheckAble(value bool) *Accessibility {
	a.selector = a.selector + "checkAble||" + b2s(value) + "\n"
	return a
}

// Selected 设置选择器的 selected 属性，用于匹配控件是否被选中
func (a *Accessibility) Selected(value bool) *Accessibility {
	a.selector = a.selector + "selected||" + b2s(value) + "\n"
	return a
}

// Enabled 设置选择器的 enabled 属性，用于匹配控件是否启用
func (a *Accessibility) Enabled(value bool) *Accessibility {
	a.selector = a.selector + "enabled||" + b2s(value) + "\n"
	return a
}

// ScrollAble 设置选择器的 scrollAble 属性，用于匹配控件是否可滚动
func (a *Accessibility) ScrollAble(value bool) *Accessibility {
	a.selector = a.selector + "scrollAble||" + b2s(value) + "\n"
	return a
}

// Editable 设置选择器的 editable 属性，用于匹配控件是否可编辑
func (a *Accessibility) Editable(value bool) *Accessibility {
	a.selector = a.selector + "editable||" + b2s(value) + "\n"
	return a
}

// MultiLine 设置选择器的 multiLine 属性，用于匹配控件是否多行
func (a *Accessibility) MultiLine(value bool) *Accessibility {
	a.selector = a.selector + "multiLine||" + b2s(value) + "\n"
	return a
}

// Checked 设置选择器的 checked 属性，用于匹配控件是否被勾选
func (a *Accessibility) Checked(value bool) *Accessibility {
	a.selector = a.selector + "checked||" + b2s(value) + "\n"
	return a
}

// Focusable 设置选择器的 focusable 属性，用于匹配控件是否可聚焦
func (a *Accessibility) Focusable(value bool) *Accessibility {
	a.selector = a.selector + "focusable||" + b2s(value) + "\n"
	return a
}

// Dismissable 设置选择器的 dismissable 属性，用于匹配控件是否可解散
func (a *Accessibility) Dismissable(value bool) *Accessibility {
	a.selector = a.selector + "dismissable||" + b2s(value) + "\n"
	return a
}

// AccessibilityFocused 设置选择器的 accessibilityFocused 属性，用于匹配控件是否是辅助功能焦点
func (a *Accessibility) AccessibilityFocused(value bool) *Accessibility {
	a.selector = a.selector + "accessibilityFocused||" + b2s(value) + "\n"
	return a
}

// ContextClickable 设置选择器的 contextClickable 属性，用于匹配控件是否是上下文点击
func (a *Accessibility) ContextClickable(value bool) *Accessibility {
	a.selector = a.selector + "contextClickable||" + b2s(value) + "\n"
	return a
}

// IndexInParent 设置选择器的 indexInParent 属性，用于匹配控件在父控件中的索引
func (a *Accessibility) IndexInParent(value int) *Accessibility {
	a.selector = a.selector + "indexInParent||" + i2s(value) + "\n"
	return a
}

// Click 点击屏幕上的文本
func (a *Accessibility) Click(text string) bool {
	obj := a.Text(text).FindOnce()
	if obj != nil {
		return obj.Click() || obj.GetParent().Click()
	} else {
		obj := a.Desc(text).FindOnce()
		if obj != nil {
			return obj.Click() || obj.GetParent().Click()
		}
	}
	return false
}

// WaitFor 等待控件出现并返回 UiObject 对象 超时单位为毫秒,写0代表无限等待,超时返回nil
func (a *Accessibility) WaitFor(timeout int) *UiObject {
	cSelector := C.CString(a.selector)
	defer C.free(unsafe.Pointer(cSelector))
	var str string
	startTime := time.Now()
	for {
		str = C.GoString(C.callJavaFindOnce(C.int(a.index), cSelector))
		if str != "" {
			break
		}
		if timeout > 0 && time.Since(startTime).Milliseconds() >= int64(timeout) {
			return nil
		}
		sleep(100)
	}
	return &UiObject{objStr: str}
}

// FindOnce 查找单个控件并返回 UiObject 对象
func (a *Accessibility) FindOnce() *UiObject {
	cSelector := C.CString(a.selector)
	defer C.free(unsafe.Pointer(cSelector))
	a.selector = ""
	str := C.GoString(C.callJavaFindOnce(C.int(a.index), cSelector))
	if str == "" {
		return nil
	}
	return &UiObject{objStr: str}
}

// Find 查找所有符合条件的控件并返回 UiObject 对象数组
func (a *Accessibility) Find() []*UiObject {
	cSelector := C.CString(a.selector)
	defer C.free(unsafe.Pointer(cSelector))
	a.selector = ""
	str := C.GoString(C.callJavaFind(C.int(a.index), cSelector))
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

func b2s(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func i2s(i int) string {
	return strconv.Itoa(i)
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}
