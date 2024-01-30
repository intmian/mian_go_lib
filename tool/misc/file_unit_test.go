package misc

import (
	"os"
	"testing"
)

func TestTFileUnit(t *testing.T) {
	jsonFileAddr := "./test.json"
	tomlFileAddr := "./test.toml"
	t.Run("jsonWR", func(t *testing.T) {
		jsonStruct := &struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{
			Name: "name",
			Age:  19,
		}
		err := GJsonTool.Write(jsonFileAddr, jsonStruct)
		if err != nil {
			t.Error(err)
		}
		jsonStruct.Name = ""
		jsonStruct.Age = 0
		err = GJsonTool.Read(jsonFileAddr, jsonStruct)
		if err != nil {
			t.Error(err)
		}
		if jsonStruct.Name != "name" || jsonStruct.Age != 19 {
			t.Error("json write or read error")
		}
	})
	t.Run("tomlWR", func(t *testing.T) {
		tomlStruct := &struct {
			Name string `toml:"name"`
			Age  int    `toml:"age"`
		}{
			Name: "name",
			Age:  19,
		}
		err := GTomlTool.Write(tomlFileAddr, tomlStruct)
		if err != nil {
			t.Error(err)
		}
		tomlStruct.Name = ""
		tomlStruct.Age = 0
		err = GTomlTool.Read(tomlFileAddr, tomlStruct)
		if err != nil {
			t.Error(err)
		}
		if tomlStruct.Name != "name" || tomlStruct.Age != 19 {
			t.Error("toml write or read error")
		}
	})
	t.Run("jsonFileUnit", func(t *testing.T) {
		type json struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		j := NewFileUnit[json](FileUnitJson, jsonFileAddr)
		err := j.Load()
		if err != nil {
			t.Error(err)
		}
		j.SaveUseData(func(a json) {
			if a.Name != "name" || a.Age != 19 {
				t.Error("json write or read error")
			}
		}, false)
		err = j.Save()
		if err != nil {
			t.Error(err)
		}
		err = j.SaveOther(jsonFileAddr)
		if err != nil {
			t.Error(err)
		}
		err = j.Load()
		if err != nil {
			t.Error(err)
		}
		j.SaveUseData(func(a json) {
			if a.Name != "name" || a.Age != 19 {
				t.Error("json write or read error")
			}
		}, true)
	})
	t.Run("test file delete", func(t *testing.T) {
		err1 := os.Remove(jsonFileAddr)
		if err1 != nil {
			t.Error(err1)
		}
		err2 := os.Remove(tomlFileAddr)
		if err2 != nil {
			t.Error(err2)
		}
	})
}
