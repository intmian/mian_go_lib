package misc

import (
	"path"
	"testing"
)

func TestFile(t *testing.T) {
	defer func() {
		_ = DeleteFileAndChild("test")
		_ = DeleteFileAndChild("test2")
	}()

	err := MakeEmptyFile("test", true)
	if err != nil {
		t.Fatal(err)
	}
	if !IsDir("test") {
		t.Fatal("test is not a dir")
	}
	err = MakeEmptyFile(path.Join("test", "test.txt"), false)
	if err != nil {
		t.Fatal(err)
	}
	if !PathExist(path.Join("test", "test.txt")) {
		t.Fatal("test.txt is not exist")
	}
	MakeEmptyFile("test2", true)

	f, err := GetFile(path.Join("test", "test.txt"))
	if err != nil {
		t.Fatal(err)
	}
	fSystem, err := GetFileTree("test")
	if err != nil {
		t.Fatal(err)
	}
	if fSystem.children[0].File.Equal(f) {
		t.Fatal("GetFileTree error")
	}

	err = f.Move("test2")
	if err != nil {
		t.Fatal(err)
	}
	if !PathExist(path.Join("test2", "test.txt")) {
		t.Fatal("test.txt is not exist")
	}
	err = f.Rename("test2.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = f.Delete()
	if err != nil {
		t.Fatal(err)
	}
	if PathExist(path.Join("test2", "test2.txt")) {
		t.Fatal("test2.txt is exist")
	}
}
