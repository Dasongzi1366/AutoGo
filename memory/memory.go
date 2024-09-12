package memory

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type Memory struct {
	pid       int
	memHandle *os.File
}

type region struct {
	Start uintptr
	End   uintptr
}

// New 需要root权限
func New(pid int) *Memory {
	memPath := fmt.Sprintf("/proc/%d/mem", pid)
	handle, err := os.Open(memPath)
	if err != nil {
		return nil
	}
	return &Memory{pid: pid, memHandle: handle}
}

func (m *Memory) SearchInt32(address []uintptr, value int32) []uintptr {
	var result []uintptr
	fd := int(m.memHandle.Fd())
	if address == nil {
		regions := readmaps(m.pid)
		var buff [1024]int32
		for _, region := range regions {
			c := (region.End - region.Start) / 4096
			for j := uintptr(0); j < c; j++ {
				buffBytes := (*[4096]byte)(unsafe.Pointer(&buff[0]))[:]
				_, err := syscall.Pread(fd, buffBytes, int64(region.Start+j*4096))
				if err != nil {
					continue
				}
				for i := 0; i < 1024; i++ {
					if buff[i] == value {
						addr := region.Start + j*4096 + uintptr(i)*4
						result = append(result, addr)
					}
				}
			}
		}
	} else {
		for _, addres := range address {
			var singleBuff int32
			singleBuffBytes := (*[4]byte)(unsafe.Pointer(&singleBuff))[:]
			_, err := syscall.Pread(fd, singleBuffBytes, int64(addres))
			if err != nil {
				continue
			}
			if singleBuff == value {
				result = append(result, addres)
			}
		}
	}
	return result
}

func (m *Memory) SearchInt64(address []uintptr, value int64) []uintptr {
	var result []uintptr
	fd := int(m.memHandle.Fd())
	if address == nil {
		regions := readmaps(m.pid)
		var buff [512]int64
		for _, region := range regions {
			c := (region.End - region.Start) / 4096
			for j := uintptr(0); j < c; j++ {
				buffBytes := (*[4096]byte)(unsafe.Pointer(&buff[0]))[:]
				_, err := syscall.Pread(fd, buffBytes, int64(region.Start+j*4096))
				if err != nil {
					continue
				}
				for i := 0; i < 512; i++ {
					if buff[i] == value {
						addr := region.Start + j*4096 + uintptr(i)*8
						result = append(result, addr)
					}
				}
			}
		}
	} else {
		for _, addres := range address {
			var singleBuff int64
			singleBuffBytes := (*[8]byte)(unsafe.Pointer(&singleBuff))[:]
			_, err := syscall.Pread(fd, singleBuffBytes, int64(addres))
			if err != nil {
				continue
			}
			if singleBuff == value {
				result = append(result, addres)
			}
		}
	}
	return result
}

func (m *Memory) ReadInt32(address uintptr) (int32, error) {
	var value int32
	valueBuff := (*[4]byte)(unsafe.Pointer(&value))[:]
	_, err := syscall.Pread(int(m.memHandle.Fd()), valueBuff, int64(address))
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (m *Memory) ReadInt64(address uintptr) (int64, error) {
	var value int64
	valueBuff := (*[8]byte)(unsafe.Pointer(&value))[:]
	_, err := syscall.Pread(int(m.memHandle.Fd()), valueBuff, int64(address))
	if err != nil {
		return 0, err
	}
	return value, nil
}

func readmaps(pid int) []region {
	mapsPath := fmt.Sprintf("/proc/%d/maps", pid)
	mapsFile, err := os.Open(mapsPath)
	if err != nil {
		return nil
	}
	defer mapsFile.Close()

	var regions []region
	scanner := bufio.NewScanner(mapsFile)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "rw") {
			fields := strings.Fields(line)
			addressRange := fields[0]
			addresses := strings.Split(addressRange, "-")
			start, err1 := strconv.ParseUint(addresses[0], 16, 64)
			end, err2 := strconv.ParseUint(addresses[1], 16, 64)
			if err1 != nil || err2 != nil {
				return nil
			}

			regions = append(regions, region{
				Start: uintptr(start),
				End:   uintptr(end),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil
	}

	return regions
}

func (m *Memory) Close() {
	_ = m.memHandle.Close()
}
