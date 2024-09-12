package app

/*
#include <stdlib.h>
const char* callJavaStringString(const char *name, const char *obj);
*/
import "C"
import (
	"mime"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"
)

type IntentOptions struct {
	Action      string
	Type        string
	Data        string
	Category    []string
	PackageName string
	ClassName   string
	Extras      map[string]string
	Flags       []string
}

// GetPackageAndClass 获取当前页面应用包名和类名
func GetPackageAndClass() (string, string) {
	re := regexp.MustCompile(`mCurrentFocus=Window\{[^}]+\s([^\s/]+)/([^\s}]+)`)
	output := shell("dumpsys window | grep mCurrentFocus")
	match := re.FindStringSubmatch(output)
	if len(match) > 2 {
		return match[1], match[2]
	}
	return "", ""
}

// Launch 通过应用包名启动应用。如果该包名对应的应用不存在，则返回false；否则返回true。
func Launch(packageName string) bool {
	return strings.Contains(shell("am start -n $(cmd package resolve-activity --brief "+packageName+" android.intent.action.MAIN | grep "+packageName+"/)"), "Starting")
}

// OpenAppSetting 打开应用的详情页(设置页)。如果找不到该应用，返回false; 否则返回true。
func OpenAppSetting(packageName string) bool {
	return !strings.Contains(shell(`am start -a android.settings.APPLICATION_DETAILS_SETTINGS -d "package:"+packageName`), "Error:")
}

// ViewFile 用其他应用查看文件。文件不存在的情况由查看文件的应用处理。
func ViewFile(path string) {
	StartActivity(IntentOptions{
		Action: "VIEW",
		Type:   getMimeType(path),
		Data:   "file://" + path,
	})
}

// EditFile 用其他应用编辑文件。文件不存在的情况由编辑文件的应用处理
func EditFile(path string) {
	StartActivity(IntentOptions{
		Action: "EDIT",
		Type:   getMimeType(path),
		Data:   "file://" + path,
	})
}

// Uninstall 卸载应用
func Uninstall(packageName string) {
	shell("pm uninstall " + packageName)
}

// Install 安装应用
func Install(path string) {
	shell("cp -rf " + path + " /data/local/tmp/1.apk;pm install /data/local/tmp/1.apk")
}

// IsInstalled 判断是否已经安装某个应用
func IsInstalled(packageName string) bool {
	return shell("pm list packages | grep "+packageName) != ""
}

// Clear 清除应用数据
func Clear(packageName string) {
	shell("pm clear " + packageName)
}

// ForceStop 强制停止应用
func ForceStop(packageName string) {
	shell("am force-stop " + packageName)
}

// Disable 禁用应用
func Disable(packageName string) {
	shell("pm disable-user " + packageName)
}

// Enable 启用应用
func Enable(packageName string) {
	shell("pm enable " + packageName)
}

// OpenUrl 用浏览器打开网站url。
func OpenUrl(url string) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	StartActivity(IntentOptions{
		Action: "VIEW",
		Data:   url,
	})
}

// StartActivity 根据选项构造一个Intent，并启动该Activity。
func StartActivity(options IntentOptions) {
	shell(buildIntentCommand(options, "start"))
}

// SendBroadcast 根据选项构造一个Intent，并发送该广播。
func SendBroadcast(options IntentOptions) {
	shell(buildIntentCommand(options, "broadcast"))
}

// StartService 根据选项构造一个Intent，并启动该服务。
func StartService(options IntentOptions) {
	shell(buildIntentCommand(options, "startservice"))
}

func buildIntentCommand(options IntentOptions, commandType string) string {
	var commandBuilder strings.Builder

	commandBuilder.WriteString("am " + commandType)

	if options.Action != "" {
		if strings.HasPrefix(options.Action, "android.intent.action.") {
			commandBuilder.WriteString(" -a " + options.Action)
		} else {
			commandBuilder.WriteString(" -a android.intent.action." + options.Action)
		}
	}

	if options.Type != "" {
		commandBuilder.WriteString(" -t " + options.Type)
	}

	if options.Data != "" {
		commandBuilder.WriteString(" -d " + options.Data)
	}

	for _, category := range options.Category {
		commandBuilder.WriteString(" -c " + category)
	}

	if options.PackageName != "" {
		commandBuilder.WriteString(" -n " + options.PackageName)
		if options.ClassName != "" {
			commandBuilder.WriteString("/" + options.ClassName)
		}
	}

	for key, value := range options.Extras {
		commandBuilder.WriteString(" --es " + key + " \"" + value + "\"")
	}

	for _, flag := range options.Flags {
		commandBuilder.WriteString(" --ez " + flag)
	}

	return commandBuilder.String()
}

func getMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".apk":
		return "application/vnd.android.package-archive"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".flac":
		return "audio/flac"
	case ".aac":
		return "audio/aac"
	case ".m4a":
		return "audio/mp4"
	default:
		// Try to detect mime type based on file extension
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			return "*/*"
		}
		return mimeType
	}
}

func shell(cmd string) string {
	cStr1 := C.CString("shell")
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(cmd)
	defer C.free(unsafe.Pointer(cStr2))
	return C.GoString(C.callJavaStringString(cStr1, cStr2))
}
