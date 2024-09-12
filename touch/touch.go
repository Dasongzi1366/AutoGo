package touch

/*
void callJavaTouchDown(int clientX, int clientY, int fingerID);
void callJavaTouchMove(int clientX, int clientY, int fingerID);
void callJavaTouchUp(int clientX, int clientY, int fingerID);
void callJavaSwipe(int startX, int startY, int endX, int endY, int duration);
void callJavaSwipeWithBezier(int startX, int startY, int endX, int endY, int duration);
*/
import "C"
import (
	"math/rand"
	"time"
)

// Down 手指按下 手指ID范围1-10
func Down(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	C.callJavaTouchDown(C.int(x), C.int(y), C.int(fingerID))
}

// Move 手指移动 手指ID范围1-10
func Move(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	C.callJavaTouchMove(C.int(x), C.int(y), C.int(fingerID))
}

// Up 手指弹起 手指ID范围1-10
func Up(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	C.callJavaTouchUp(C.int(x), C.int(y), C.int(fingerID))
}

// Click 手指点击
func Click(x, y, fingerID int) {
	Down(x, y, fingerID)
	sleep(random(10, 20))
	Up(x, y, fingerID)
	sleep(10)
}

// LongClick 手指长按 按住时长单位毫秒
func LongClick(x, y, duration int) {
	Down(x, y, 1)
	sleep(duration + random(1, 20))
	Up(x, y, 1)
	sleep(10)
}

// Swipe 模拟从坐标(x1, y1)滑动到坐标(x2, y2) 滑动耗时单位毫秒
func Swipe(x1, y1, x2, y2, duration int) {
	C.callJavaSwipe(C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(duration))
}

// Swipe2 使用贝塞尔曲线方式进行滑动
func Swipe2(x1, y1, x2, y2, duration int) {
	C.callJavaSwipeWithBezier(C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(duration))
}

func random(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}
