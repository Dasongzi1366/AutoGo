package files

import "C"
import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IsFile 返回路径path是否是文件。
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDir 返回路径path是否是文件夹。
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsEmptyDir 返回文件夹path是否为空文件夹。如果该路径并非文件夹，则直接返回false。
func IsEmptyDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(files) == 0
}

// Create 创建一个文件或文件夹并返回是否创建成功。如果文件已经存在，则直接返回false。
func Create(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// 文件或文件夹已经存在
		return false
	}
	if os.IsNotExist(err) {
		// 尝试创建文件夹
		if filepath.Ext(path) == "" {
			err = os.Mkdir(path, 0755)
		} else {
			// 尝试创建文件
			file, err := os.Create(path)
			if err != nil {
				return false
			}
			_ = file.Close()
		}
	}
	return err == nil
}

// CreateWithDirs 创建一个文件或文件夹并返回是否创建成功。如果文件所在文件夹不存在，则先创建它所在的一系列文件夹。如果文件已经存在，则直接返回false。
func CreateWithDirs(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// 文件或文件夹已经存在
		return false
	}
	if os.IsNotExist(err) {
		// 确保路径中的所有文件夹存在
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return false
		}
		// 尝试创建文件
		file, err := os.Create(path)
		if err != nil {
			return false
		}
		_ = file.Close()
	}
	return err == nil
}

// Exists 返回在路径path处的文件是否存在。
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsureDir 确保路径path所在的文件夹存在。如果该路径所在文件夹不存在，则创建该文件夹。
func EnsureDir(path string) bool {
	return os.MkdirAll(filepath.Dir(path), 0755) == nil
}

// Read 读取文本文件path的所有内容并返回。
func Read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// ReadBytes 读取文件path的所有内容并返回。
func ReadBytes(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

// Write 把text写入到文件path中。如果文件存在则覆盖，不存在则创建。
func Write(path, text string) {
	_ = os.WriteFile(path, []byte(text), 0644)
}

// WriteBytes 把bytes写入到文件path中。如果文件存在则覆盖，不存在则创建。
func WriteBytes(path string, bytes []byte) {
	Remove(path)
	_ = os.WriteFile(path, bytes, 0644)
}

// Append 把text追加到文件path的末尾。如果文件不存在则创建。
func Append(path string, text string) {
	// 以追加模式打开文件，如果文件不存在则创建
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	// 将文本写入文件末尾
	_, _ = file.WriteString(text)
}

// AppendBytes 把bytes追加到文件path的末尾。如果文件不存在则创建。
func AppendBytes(path string, bytes []byte) {
	// 以追加模式打开文件，如果文件不存在则创建
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	// 将字节数组写入文件末尾
	_, _ = file.Write(bytes)
}

// Copy 复制文件，返回是否复制成功。
func Copy(fromPath, toPath string) bool {
	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return false
	}
	defer sourceFile.Close()

	destFile, err := os.Create(toPath)
	if err != nil {
		return false
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err == nil
}

// Move 移动文件，返回是否移动成功。
func Move(fromPath, toPath string) bool {
	return os.Rename(fromPath, toPath) == nil
}

// Rename 重命名文件，并返回是否重命名成功。
func Rename(path, newName string) bool {
	return os.Rename(path, newName) == nil
}

// GetName 返回文件的文件名。例如files.getName("/sdcard/1.txt")返回"1.txt"。
func GetName(path string) string {
	return filepath.Base(path)
}

// GetNameWithoutExtension 返回不含拓展名的文件的文件名。例如files.getName("/sdcard/1.txt")返回"1"。
func GetNameWithoutExtension(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

// GetExtension 返回文件的拓展名。例如files.getExtension("/sdcard/1.txt")返回"txt"。
func GetExtension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 && ext[0] == '.' {
		return ext[1:]
	}
	return ext
}

// Remove 删除文件或文件夹，如果是文件夹会删除整个文件夹包含里面的所有文件，返回是否删除成功。
func Remove(path string) bool {
	return os.RemoveAll(path) == nil
}

// Path 返回相对路径对应的绝对路径。例如files.path("./1.png")，如果运行这个语句的脚本位于文件夹"/sdcard/脚本/"中，则返回"/sdcard/脚本/1.png"。
func Path(relativePath string) string {
	dir, _ := os.Getwd()
	return dir + "/" + GetName(relativePath)
}

// ListDir 列出文件夹path下的所有文件和文件夹
func ListDir(path string) []string {
	var entries []string
	dir, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil
	}
	for _, file := range files {
		entries = append(entries, filepath.Join(path, file.Name()))
	}
	return entries
}
