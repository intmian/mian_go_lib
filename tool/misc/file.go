package misc

import (
	"io"
	"os"
	"path/filepath"
)

// ReadAllDir 返回root目录下所有的目录，当出现错误时返回nil，否则至少返回空
func ReadAllDir(root string) []string {
	rd, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	files := []string{} // 应当无视IDE给出的警告，遵守的话会导致返回nil
	for _, f := range rd {
		if f.IsDir() {
			files = append(files, f.Name())
		}
	}
	return files
}

// ReadAllFile 返回root目录下所有的目录，当出现错误时返回nil，否则至少返回空数组
func ReadAllFile(root string) []string {
	rd, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	files := []string{} // 应当无视IDE给出的警告，遵守的话会导致返回nil
	for _, f := range rd {
		if !f.IsDir() {
			files = append(files, f.Name())
		}
	}
	return files
}

// PathExist 返回文件或目录是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && err.Error() == os.ErrNotExist.Error() {
		return false
	}
	return true
}

// CopyFile 复制文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
		}
	}(src)
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
		}
	}(dst)
	return io.Copy(dst, src)
}

func GetNameFromPath(path string) string {
	return filepath.Base(path)
}

func DeleteFileAndChild(path string) error {
	return os.RemoveAll(path)
}
