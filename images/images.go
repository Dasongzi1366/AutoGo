package images

/*
#include <stdlib.h>
void callJavaLoadSo(const char* str);
int callJavaPngToMat(const char* sign, const char *data, int dataLength);
void callJavaFindImage(const char* sign, double threshold, int x1, int y1, int x2, int y2, int *resultX, int *resultY);
void callJavaFindColor(int color, double threshold, int x1, int y1, int x2, int y2, int *resultX, int *resultY);
void callJavaFindMultiColors(int color, int *data, int dataLength, double threshold, int x1, int y1, int x2, int y2, int *resultX, int *resultY);
const char* callJavaGetPixel(int x, int y);
void callJavaGetBitmapPixels(int **pixels, int *length);
const char* callJavaStringString(const char *name, const char *obj);
int callJavaGetColorCountInRegion(int color, double threshold, int x1, int y1, int x2, int y2);
*/
import "C"
import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/bmp"
)

type colorsTemplate struct {
	color int
	data  []int32
}

var signMap = make(map[string]int) //值0代表无效图像 1代表已经加载
var colorMap = make(map[string]colorsTemplate)

func init() {
	loadSo("libc++_shared.so")
	loadSo("libopencv_java4.so")
}

// Read 读取在路径path的图片文件并返回一个Image对象。
// 如果文件不存在或者文件无法解码则返回nil。
//
// 参数:
// - path: 要读取的图片文件路径。
//
// 返回值:
// - 成功时返回指向image.NRGBA对象的指针，否则返回nil。
func Read(path string) *image.NRGBA {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

// Load 加载地址URL的网络图片并返回一个Image对象。
// 如果地址不存在或者图片无法解码则返回nil。
//
// 参数:
// - url: 要下载的图片的URL地址。
//
// 返回值:
// - 成功时返回指向image.NRGBA对象的指针，否则返回nil。
func Load(url string) *image.NRGBA {
	client := http.Client{
		Timeout: time.Duration(5000) * time.Millisecond,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

// Save 把图片image保存到path中。
// 如果文件不存在会被创建；文件存在会被覆盖。
//
// 参数:
// - img: 要保存的image.NRGBA对象。
// - path: 保存图片的文件路径。
// - quality: 保存图片的质量。
//
// 返回值:
// - 成功时返回true，否则返回false。
func Save(img *image.NRGBA, path string, quality int) bool {
	ext := strings.ToLower(filepath.Ext(path))
	ext = strings.ReplaceAll(ext, ".", "")
	data := ToBytes(img, ext, quality)
	if data != nil {
		return os.WriteFile(path, data, 0644) == nil
	}
	return false
}

// FromBase64 解码Base64数据并返回解码后的图片Image对象。
// 如果base64无法解码则返回nil。
//
// 参数:
// - base64Str: 要解码的Base64字符串。
//
// 返回值:
// - 成功时返回指向image.NRGBA对象的指针，否则返回nil。
func FromBase64(base64Str string) *image.NRGBA {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil
	}
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

// ToBase64 把Image对象编码为Base64数据并返回。
//
// 参数:
// - img: 要编码的image.NRGBA对象。
// - format: 编码的图片格式（如"png", "jpg"等）。
// - quality: 编码的图片质量。
//
// 返回值:
// - 编码后的Base64字符串。
func ToBase64(img *image.NRGBA, format string, quality int) string {
	data := ToBytes(img, format, quality)
	if data != nil {
		return base64.StdEncoding.EncodeToString(data)
	}
	return ""
}

// FromBytes 解码字节数组并返回解码后的图片Image对象。
// 如果字节数组无法解码则返回nil。
//
// 参数:
// - data: 要解码的字节数组。
//
// 返回值:
// - 成功时返回指向image.NRGBA对象的指针，否则返回nil。
func FromBytes(data []byte) *image.NRGBA {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

// ToBytes 把图片编码为字节数组并返回。
//
// 参数:
// - img: 要编码的image.NRGBA对象。
// - format: 编码的图片格式（如"png", "jpg"等）。
// - quality: 编码的图片质量。
//
// 返回值:
// - 编码后的字节数组。

func ToBytes(img *image.NRGBA, format string, quality int) []byte {
	var err error
	var buf bytes.Buffer
	format = strings.ToLower(format)
	switch format {
	case "png":
		err = png.Encode(&buf, img)
	case "jpg", "jpeg":
		options := &jpeg.Options{Quality: quality}
		err = jpeg.Encode(&buf, img, options)
	case "bmp":
		err = bmp.Encode(&buf, img)
	default:
		err = os.ErrInvalid
	}
	if err == nil {
		return buf.Bytes()
	}
	return nil
}

// Clip 从图片img的位置(x1, y1)处剪切至(x2, y2)区域，并返回该剪切区域的新图片。
//
// 参数:
// - img: 要剪切的image.NRGBA对象。
// - x1, y1: 剪切区域的左上角坐标。
// - x2, y2: 剪切区域的右下角坐标。为0时默认使用图片最大宽高
//
// 返回值:
// - 剪切后的image.NRGBA对象。

func Clip(img *image.NRGBA, x1, y1, x2, y2 int) *image.NRGBA {
	bounds := img.Bounds()
	if x1 < bounds.Min.X {
		x1 = bounds.Min.X
	}
	if y1 < bounds.Min.Y {
		y1 = bounds.Min.Y
	}
	if x2 > bounds.Max.X || x2 == 0 {
		x2 = bounds.Max.X
	}
	if y2 > bounds.Max.Y || y2 == 0 {
		y2 = bounds.Max.Y
	}
	width := x2 - x1
	height := y2 - y1
	newImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	var wg sync.WaitGroup
	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				srcIdx := (y1+y)*img.Stride + (x1+x)*4
				dstIdx := y*newImg.Stride + x*4
				copy(newImg.Pix[dstIdx:dstIdx+4], img.Pix[srcIdx:srcIdx+4])
			}
		}(y)
	}
	wg.Wait()
	return newImg
}

// Resize 调整图片大小，并返回调整后的图片。
//
// 参数:
// - img: 要调整大小的image.NRGBA对象。
// - width: 调整后的宽度。
// - height: 调整后的高度。
//
// 返回值:
// - 调整大小后的image.NRGBA对象。
func Resize(img *image.NRGBA, width, height int) *image.NRGBA {
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)
	resizedNRGBA := image.NewNRGBA(resizedImg.Bounds())
	copy(resizedNRGBA.Pix, resizedImg.Pix)
	return resizedNRGBA
}

// Rotate 将图片顺时针旋转degree度，返回旋转后的图片对象。
//
// 参数:
// - img: 要旋转的image.NRGBA对象。
// - degree: 顺时针旋转的角度。
//
// 返回值:
// - 旋转后的image.NRGBA对象。
func Rotate(img *image.NRGBA, degree int) *image.NRGBA {
	rotatedImg := imaging.Rotate(img, float64(-degree), color.Transparent)
	return rotatedImg
}

// Grayscale 将图片灰度化，并返回灰度化后的图片。
//
// 参数:
// - img: 要灰度化的image.NRGBA对象。
//
// 返回值:
// - 灰度化后的image.Gray对象。
func Grayscale(img *image.NRGBA) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				originalColor := img.At(x, y)
				grayColor := color.GrayModel.Convert(originalColor)
				grayImg.Set(x, y, grayColor)
			}
		}(y)
	}
	wg.Wait()
	return grayImg
}

// Threshold 将图片阈值化，并返回处理后的图像。
//
// 参数:
// - img: 要处理的image.NRGBA对象。
// - threshold: 阈值。
// - maxVal: 阈值化后的最大值。
// - typ: 阈值化类型（如"BINARY", "BINARY_INV", "TRUNC", "TOZERO", "TOZERO_INV"）。
//
// 返回值:
// - 阈值化处理后的image.Gray对象。
func Threshold(img *image.NRGBA, threshold, maxVal int, typ string) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			pixel := originalColor.Y
			var newPixel uint8
			switch typ {
			case "BINARY":
				if int(pixel) > threshold {
					newPixel = uint8(maxVal)
				} else {
					newPixel = 0
				}
			case "BINARY_INV":
				if int(pixel) > threshold {
					newPixel = 0
				} else {
					newPixel = uint8(maxVal)
				}
			case "TRUNC":
				if int(pixel) > threshold {
					newPixel = uint8(threshold)
				} else {
					newPixel = pixel
				}
			case "TOZERO":
				if int(pixel) > threshold {
					newPixel = pixel
				} else {
					newPixel = 0
				}
			case "TOZERO_INV":
				if int(pixel) > threshold {
					newPixel = 0
				} else {
					newPixel = pixel
				}
			default:
				newPixel = pixel
			}
			grayImg.Set(x, y, color.Gray{Y: newPixel})
		}
	}
	return grayImg
}

// AdaptiveThreshold 将图像进行自适应阈值化处理，并返回处理后的图像。
//
// 参数:
// - img: 要处理的image.NRGBA对象。
// - maxValue: 阈值化后的最大值。
// - adaptiveMethod: 自适应方法（如"MEAN_C", "GAUSSIAN_C"）。
// - thresholdType: 阈值化类型（如"BINARY", "BINARY_INV"）。
// - blockSize: 计算阈值的区域大小。
// - C: 常数值，用于调整计算出的阈值。
//
// 返回值:
// - 自适应阈值化处理后的image.Gray对象。
func AdaptiveThreshold(img *image.NRGBA, maxValue float64, adaptiveMethod string, thresholdType string, blockSize int, C float64) *image.Gray {
	bounds := img.Bounds()
	grayImg := Grayscale(img) // 先将图像灰度化
	dstImg := image.NewGray(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				var sum, sumWeight, mean float64
				for j := -blockSize / 2; j <= blockSize/2; j++ {
					for i := -blockSize / 2; i <= blockSize/2; i++ {
						xx := x + i
						yy := y + j
						if xx >= bounds.Min.X && xx < bounds.Max.X && yy >= bounds.Min.Y && yy < bounds.Max.Y {
							pixel := grayImg.GrayAt(xx, yy).Y
							if adaptiveMethod == "GAUSSIAN_C" {
								weight := math.Exp(-(float64(i*i + j*j)) / (2.0 * float64(blockSize*blockSize)))
								sum += float64(pixel) * weight
								sumWeight += weight
							} else {
								sum += float64(pixel)
								sumWeight += 1.0
							}
						}
					}
				}

				if adaptiveMethod == "GAUSSIAN_C" {
					mean = sum / sumWeight
				} else {
					mean = sum / (float64(blockSize) * float64(blockSize))
				}

				threshold := mean - C
				pixel := grayImg.GrayAt(x, y).Y

				var newPixel uint8
				switch thresholdType {
				case "BINARY":
					if float64(pixel) > threshold {
						newPixel = uint8(maxValue)
					} else {
						newPixel = 0
					}
				case "BINARY_INV":
					if float64(pixel) > threshold {
						newPixel = 0
					} else {
						newPixel = uint8(maxValue)
					}
				}

				dstImg.SetGray(x, y, color.Gray{Y: newPixel})
			}
		}(y)
	}
	wg.Wait()
	return dstImg
}

// Binarize 将图像进行二值化处理，颜色值大于threshold的变成255，否则变成0。
//
// 参数:
// - img: 要处理的image.NRGBA对象。
// - threshold: 阈值。
//
// 返回值:
// - 二值化处理后的image.Gray对象。
func Binarize(img *image.NRGBA, threshold int) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				originalColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
				pixel := originalColor.Y
				var newPixel uint8
				if int(pixel) > threshold {
					newPixel = uint8(255)
				} else {
					newPixel = 0
				}
				grayImg.SetGray(x, y, color.Gray{Y: newPixel})
			}
		}(y)
	}
	wg.Wait()
	return grayImg
}

// ImageToNRGBA 将通用的Image对象转换为NRGBA对象。
//
// 参数:
// - img: 要转换的image.Image对象。
//
// 返回值:
// - 转换后的image.NRGBA对象。
func ImageToNRGBA(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	nrgbaImg := image.NewNRGBA(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				srcColor := img.At(x, y)
				r, g, b, a := srcColor.RGBA()
				i := (y-bounds.Min.Y)*nrgbaImg.Stride + (x-bounds.Min.X)*4
				nrgbaImg.Pix[i] = uint8(r >> 8)
				nrgbaImg.Pix[i+1] = uint8(g >> 8)
				nrgbaImg.Pix[i+2] = uint8(b >> 8)
				nrgbaImg.Pix[i+3] = uint8(a >> 8)
			}
		}(y)
	}
	wg.Wait()
	return nrgbaImg
}

// FindImage 在屏幕区域内查找指定模板图片，并返回找到的位置。
//
// 参数：
// - template: png模板图片的字节数组。
// - threshold: 匹配阈值，范围在0到1之间。
// - x1, y1: 屏幕区域的左上角坐标。
// - x2, y2: 屏幕区域的右下角坐标。为0时默认使用屏幕最大宽高
//
// 返回值：
// - 返回找到的模板图片的位置，如果没有找到则返回(-1, -1)。
func FindImage(template *[]byte, threshold float64, x1, y1, x2, y2 int) (int, int) {
	sign := i2s(int(uintptr(unsafe.Pointer(&(*template)[0]))))
	if _, ok := signMap[sign]; !ok {
		if pngToMat(sign, *template) {
			signMap[sign] = 1
		} else {
			signMap[sign] = 0
		}
	}
	if signMap[sign] == 0 {
		return -1, -1
	}
	resultX := C.int(-1)
	resultY := C.int(-1)
	cSign := C.CString(sign)
	defer C.free(unsafe.Pointer(cSign))
	C.callJavaFindImage(cSign, C.double(threshold), C.int(x1), C.int(y1), C.int(x2), C.int(y2), &resultX, &resultY)
	return int(resultX), int(resultY)
}

// FindMultiColors 在屏幕区域内查找多个颜色点，并返回找到的位置。
//
// 参数：
// - colors: 颜色模板字符串，例如: "#e22021",[[1,0,"#a86165-505050"],[292,-193,"#00af9c"],[368,-197,"#00af9c"]] 可以没有#"[]这些符号
// - threshold: 匹配阈值，范围在0到1之间。
// - x1, y1: 屏幕区域的左上角坐标。
// - x2, y2: 屏幕区域的右下角坐标。为0时默认使用屏幕最大宽高
//
// 返回值：
// - 返回找到的颜色位置，如果没有找到则返回(-1, -1)。
func FindMultiColors(colors string, threshold float64, x1, y1, x2, y2 int) (int, int) {
	if _, ok := colorMap[colors]; !ok {
		colorMap[colors] = parseColorTemplate(colors)
	}
	dataLength := len(colorMap[colors].data)
	cData := (*C.int)(unsafe.Pointer(&colorMap[colors].data[0]))
	var resultX, resultY C.int
	C.callJavaFindMultiColors(C.int(colorMap[colors].color), cData, C.int(dataLength), C.double(threshold), C.int(x1), C.int(y1), C.int(x2), C.int(y2), &resultX, &resultY)
	return int(resultX), int(resultY)
}

// FindColor 在屏幕区域内查找指定颜色，并返回找到的位置。
//
// 参数：
// - color: 要查找的颜色字符串，格式为#RGB，例如: #706B73
// - threshold: 匹配阈值，范围在0到1之间。
// - x1, y1: 屏幕区域的左上角坐标。
// - x2, y2: 屏幕区域的右下角坐标。为0时默认使用屏幕最大宽高
//
// 返回值：
// - 返回找到的颜色位置，如果没有找到则返回(-1, -1)。
func FindColor(color string, threshold float64, x1, y1, x2, y2 int) (int, int) {
	resultX := C.int(-1)
	resultY := C.int(-1)
	C.callJavaFindColor(C.int(int(hexColorToInt(color))), C.double(threshold), C.int(x1), C.int(y1), C.int(x2), C.int(y2), &resultX, &resultY)
	return int(resultX), int(resultY)
}

// GetColorCountInRegion 计算指定屏幕区域内符合颜色条件的像素数量。
//
// 参数：
// - color: 要查找的颜色字符串，格式为 #RGB，例如: #706B73。
// - threshold: 匹配阈值，范围在 0 到 1 之间。值越大匹配越严格。
// - x1, y1: 屏幕区域的左上角坐标。
// - x2, y2: 屏幕区域的右下角坐标。为 0 时默认使用屏幕的最大宽高。
//
// 返回值：
// - 返回符合条件的颜色像素数量，如果未找到符合条件的像素，则返回 0。
func GetColorCountInRegion(color string, threshold float64, x1, y1, x2, y2 int) int {
	return int(C.callJavaGetColorCountInRegion(C.int(int(hexColorToInt(color))), C.double(threshold), C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}

// Pixel 获取指定坐标点的颜色值。
//
// 参数：
// - x, y: 要获取颜色值的坐标点。
//
// 返回值：
// - 返回颜色值的字符串表示。
func Pixel(x, y int) string {
	return C.GoString(C.callJavaGetPixel(C.int(x), C.int(y)))
}

// CaptureScreen 截取当前屏幕并返回一个Image对象。
//
// 返回值：
// - 成功时返回截取的image.NRGBA对象，否则返回nil。
func CaptureScreen() *image.NRGBA {
	var pixels *C.int
	var length C.int
	C.callJavaGetBitmapPixels(&pixels, &length)
	if length < 2 || pixels == nil {
		return nil
	}
	goPixels := unsafe.Slice(pixels, int(length))
	width := int(goPixels[len(goPixels)-1])
	height := (len(goPixels) - 1) / width
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	var wg sync.WaitGroup
	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				i := y*width + x
				pixel := goPixels[i]
				r := uint8((pixel >> 16) & 0xFF)
				g := uint8((pixel >> 8) & 0xFF)
				b := uint8(pixel & 0xFF)
				a := uint8((pixel >> 24) & 0xFF)
				offset := y*img.Stride + x*4
				img.Pix[offset] = r
				img.Pix[offset+1] = g
				img.Pix[offset+2] = b
				img.Pix[offset+3] = a
			}
		}(y)
	}
	wg.Wait()
	C.free(unsafe.Pointer(pixels))
	return img
}

// PNG图片转换为Mat对象供后续找图使用
func pngToMat(sign string, pngData []byte) bool {
	cSign := C.CString(sign)
	defer C.free(unsafe.Pointer(cSign))
	return C.callJavaPngToMat(cSign, (*C.char)(unsafe.Pointer(&pngData[0])), C.int(len(pngData))) == 1
}

func parseColorTemplate(input string) colorsTemplate {
	input = strings.ReplaceAll(input, "[", "")
	input = strings.ReplaceAll(input, "]", "")
	input = strings.ReplaceAll(input, "#", "")
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, `"`, "")

	// 使用逗号分隔字符串
	parts := strings.Split(input, ",")

	// 将颜色字符串转换为整数
	colorToInt := hexColorToInt(parts[0])

	// 处理data部分
	var data []int32
	for i := 1; i < len(parts); i++ {
		// 如果是颜色字符串，转换为整数
		if len(parts[i]) == 6 || len(parts[i]) == 13 {
			if len(parts[i]) == 6 {
				value := hexColorToInt(parts[i])
				data = append(data, value)
				data = append(data, 0)
			} else {
				arr := strings.Split(parts[i], "-")
				data = append(data, hexColorToInt(arr[0]))
				data = append(data, hexColorToInt(arr[1]))
			}
		} else {
			// 否则是坐标
			value, _ := strconv.ParseInt(parts[i], 10, 32)
			data = append(data, int32(value))
		}
	}
	return colorsTemplate{
		color: int(colorToInt),
		data:  data,
	}
}

func hexColorToInt(hexColor string) int32 {
	colorToInt, err := strconv.ParseInt(strings.ReplaceAll(hexColor, "#", ""), 16, 32)
	if err != nil {
		fmt.Println("Error parsing color:", err)
		return 0
	}
	return int32(0xFF000000 | colorToInt)
}

func i2s(i int) string {
	return strconv.Itoa(i)
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
