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

// FileNode 文件树的节点，目前没有自我维护功能
type FileNode struct {
	File     FileCore
	Parent   *FileNode
	Children []FileNode
}

// FileCore 操作文件，获得文件的基础消息。
type FileCore struct {
	// 绝对路径
	Addr string
	// ./a/b/a_c.c.txt -> a_c.c.txt
	Name string
	// ./a/b/c.txt -> txt
	Extension string
	IsDir     bool
}

func (f *FileCore) GetSize() int64 {
	fileInfo, err := os.Stat(f.Addr)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

func (f *FileCore) Valid() bool {
	return PathExist(f.Addr)
}

func (f *FileCore) Equal(other FileCore) bool {
	if f.Addr != other.Addr {
		return false
	}
	if f.Name != other.Name {
		return false
	}

	if f.Extension != other.Extension {
		return false
	}
	if f.IsDir != other.IsDir {
		return false
	}
	return true
}

func (f *FileCore) Rename(newName string) error {
	err := os.Rename(f.Addr, filepath.Join(filepath.Dir(f.Addr), newName))
	if err != nil {
		return err
	}
	f.Name = newName
	f.Addr = filepath.Join(filepath.Dir(f.Addr), newName)
	f.Extension = filepath.Ext(newName)
	return nil
}

func (f *FileCore) Move(newDir string) error {
	if !IsDir(newDir) {
		return os.ErrNotExist
	}
	err := os.Rename(f.Addr, filepath.Join(newDir, f.Name))
	if err != nil {
		return err
	}
	f.Addr = filepath.Join(newDir, f.Name)
	return nil
}

func (f *FileCore) Delete() error {
	if f.IsDir {
		return os.RemoveAll(f.Addr)
	}
	return os.Remove(f.Addr)
}

func (f *FileCore) MakeChildEmptyFile(name string, isDir bool) error {
	err := MakeEmptyFile(filepath.Join(f.Addr, name), isDir)
	if err != nil {
		return err
	}
	return nil
}

func GetFile(addr string) (FileCore, error) {
	file := FileCore{
		Addr:      addr,
		Name:      filepath.Base(addr),
		Extension: filepath.Ext(addr),
		IsDir:     IsDir(addr),
	}
	return file, nil
}

// GetFileTree 返回root为根的文件树
func GetFileTree(root string) (FileNode, error) {
	node := FileNode{}
	if !PathExist(root) {
		return node, os.ErrNotExist
	}

	// 处理为绝对路径
	if !filepath.IsAbs(root) {
		root, _ = filepath.Abs(root)
	}

	// 获得文件信息
	Name := filepath.Base(root)
	Extension := filepath.Ext(root)
	isDir := IsDir(root)
	node.File = FileCore{
		Addr:      root,
		Name:      Name,
		Extension: Extension,
		IsDir:     isDir,
	}
	if !isDir {
		return node, nil
	}

	// 遍历子文件、子目录，补充信息
	children := []FileNode{}
	for _, f := range ReadAllFile(root) {
		file := FileCore{
			Addr:      filepath.Join(root, f),
			Name:      f,
			Extension: filepath.Ext(f),
			IsDir:     false,
		}
		child := FileNode{
			File:   file,
			Parent: &node,
		}
		children = append(children, child)
	}
	for _, f := range ReadAllDir(root) {
		child, err := GetFileTree(filepath.Join(root, f))
		if err != nil {
			return node, err
		}
		child.Parent = &node
		children = append(children, child)
	}
	node.Children = children

	return node, nil
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// PathExist 返回文件或目录是否存在
func PathExist(path string) bool {
	f, _ := os.Stat(path)
	if f == nil {
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

// MoveFile 移动文件
func MoveFile(srcPath, dstPath string) error {
	return os.Rename(srcPath, dstPath)
}

func MakeEmptyFile(path string, isDir bool) error {
	if isDir {
		return os.MkdirAll(path, os.ModePerm)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return file.Close()
}

func GetNameFromPath(path string) string {
	return filepath.Base(path)
}

func GetFileContent(path string) string {
	if !PathExist(path) {
		return ""
	}
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	s := ""
	buf := make([]byte, 100)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return ""
		}
		if n == 0 {
			break
		}
		s += string(buf[:n])
	}
	return s
}

func DeleteFileAndChild(path string) error {
	return os.RemoveAll(path)
}
