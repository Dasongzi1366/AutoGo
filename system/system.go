package system

/*
#include <stdlib.h>
void logI(const char* label,const char* message);
void logE(const char* label, const char* message);
const char* callJavaStringString(const char *name, const char *obj);
void callJavaLoadSo(const char* str);
*/
import "C"
import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

// GetPid 返回某个进程名的 PID，如果进程名为空则返回自身进程的 PID
// 如果未找到指定进程，返回 -1
func GetPid(processName string) int {
	if processName == "" {
		return os.Getpid()
	}

	cmd := fmt.Sprintf("ps -ef | grep '%s' | grep -v 'grep'", processName)
	output := Shell(cmd)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			pid, err := strconv.Atoi(fields[1])
			if err == nil {
				return pid
			}
		}
	}

	return -1
}

// GetMemoryUsage 返回指定进程的内存占用（以KB为单位）
// 如果pid为0，则返回当前进程的内存占用
// 查询失败返回-1
func GetMemoryUsage(pid int) int {
	if pid == 0 {
		pid = os.Getpid()
	}

	cmd := fmt.Sprintf("cat /proc/%d/status | grep -e VmRSS", pid)
	output := Shell(cmd)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "VmRSS") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memory, err := strconv.Atoi(fields[1])
				if err == nil {
					return memory
				}
			}
		}
	}

	return -1
}

// GetCpuUsage 返回指定进程的 CPU 使用率
// 如果pid为0，则返回当前进程的 CPU 使用率
// 查询失败返回0.0
func GetCpuUsage(pid int) float64 {
	if pid == 0 {
		pid = os.Getpid()
	}

	cmd := fmt.Sprintf("top -b -n 1 | grep '^ *%d '", pid)
	output := Shell(cmd)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 8 {
			cpuUsage, err := strconv.ParseFloat(fields[8], 64)
			if err == nil {
				return cpuUsage
			}
		}
	}

	return 0.0 // 返回 0.0 表示查询失败
}

func LogI(label, message string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.logI(cLabel, cMessage)
}

func LogE(label, message string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.logE(cLabel, cMessage)
}

// RestartSelf 重启自身进程
func RestartSelf() {
	os.Exit(123)
}

func Shell(cmd string) string {
	cStr1 := C.CString("shell")
	defer C.free(unsafe.Pointer(cStr1))
	cStr2 := C.CString(cmd)
	defer C.free(unsafe.Pointer(cStr2))
	return C.GoString(C.callJavaStringString(cStr1, cStr2))
}

func loadSo(name string, data []byte) {
	dir, _ := os.Getwd()
	path := dir + "/"
	_, err := os.Stat(path + name)
	if err != nil {
		_ = os.WriteFile(path+name, data, 0644)
	}
	cStr := C.CString(path + name)
	defer C.free(unsafe.Pointer(cStr))
	C.callJavaLoadSo(cStr)
}
