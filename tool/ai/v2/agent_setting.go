package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type SettingFieldDoc struct {
	Name     string `json:"name"`
	GoType   string `json:"goType"`
	JSONType string `json:"jsonType"`
	IsArray  bool   `json:"isArray"`
}

type IAgentSetting interface {
	ExportJSON() ([]byte, error)
	ImportJSON(data []byte) error
	ExportJSONDoc() ([]SettingFieldDoc, error)
}

func ExportSettingJSON(setting any) ([]byte, error) {
	value, err := settingStructValue(setting)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]any)
	if err := rangeSettingFields(value, func(field reflect.StructField, fieldValue reflect.Value, info settingFieldInfo) error {
		if isZeroSettingValue(fieldValue) {
			return nil
		}
		ret[info.name] = fieldValue.Interface()
		return nil
	}); err != nil {
		return nil, err
	}
	return json.Marshal(ret)
}

func ImportSettingJSON[T any](base T, data []byte) (T, error) {
	value := reflect.ValueOf(&base).Elem()
	if err := ensureSettingStruct(value); err != nil {
		return base, err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return base, err
	}

	err := rangeSettingFields(value, func(field reflect.StructField, fieldValue reflect.Value, info settingFieldInfo) error {
		rawValue, ok := raw[info.name]
		if !ok {
			return nil
		}
		next := reflect.New(fieldValue.Type()).Elem()
		if err := json.Unmarshal(rawValue, next.Addr().Interface()); err != nil {
			return fmt.Errorf("%s: %w", info.name, err)
		}
		if isZeroSettingValue(next) {
			return nil
		}
		fieldValue.Set(next)
		return nil
	})
	return base, err
}

func ExportSettingDoc(setting any) ([]SettingFieldDoc, error) {
	value, err := settingStructValue(setting)
	if err != nil {
		return nil, err
	}

	docs := make([]SettingFieldDoc, 0, value.NumField())
	if err := rangeSettingFields(value, func(field reflect.StructField, fieldValue reflect.Value, info settingFieldInfo) error {
		docs = append(docs, SettingFieldDoc{
			Name:     info.name,
			GoType:   field.Type.String(),
			JSONType: info.jsonType,
			IsArray:  info.isArray,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return docs, nil
}

type settingFieldInfo struct {
	name     string
	jsonType string
	isArray  bool
}

func settingStructValue(setting any) (reflect.Value, error) {
	if setting == nil {
		return reflect.Value{}, errors.New("setting is nil")
	}
	value := reflect.ValueOf(setting)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, errors.New("setting is nil")
		}
		value = value.Elem()
	}
	if err := ensureSettingStruct(value); err != nil {
		return reflect.Value{}, err
	}
	return value, nil
}

func ensureSettingStruct(value reflect.Value) error {
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return errors.New("setting must be a struct")
	}
	return nil
}

func rangeSettingFields(value reflect.Value, visit func(reflect.StructField, reflect.Value, settingFieldInfo) error) error {
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldValue := value.Field(i)
		info, err := parseSettingField(field)
		if err != nil {
			return err
		}
		if info.name == "-" {
			continue
		}
		if err := visit(field, fieldValue, info); err != nil {
			return err
		}
	}
	return nil
}

func parseSettingField(field reflect.StructField) (settingFieldInfo, error) {
	name := field.Name
	if tag := field.Tag.Get("json"); tag != "" {
		parts := strings.Split(tag, ",")
		if parts[0] != "" {
			name = parts[0]
		}
	}
	if name == "-" {
		return settingFieldInfo{name: "-"}, nil
	}

	jsonType, isArray, err := settingJSONType(field.Type)
	if err != nil {
		return settingFieldInfo{}, fmt.Errorf("%s: %w", field.Name, err)
	}
	return settingFieldInfo{
		name:     name,
		jsonType: jsonType,
		isArray:  isArray,
	}, nil
}

func settingJSONType(t reflect.Type) (string, bool, error) {
	isArray := false
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		isArray = true
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string", isArray, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number", isArray, nil
	case reflect.Float32, reflect.Float64:
		return "number", isArray, nil
	case reflect.Bool:
		return "boolean", isArray, nil
	default:
		return "", false, errors.New("unsupported setting field type")
	}
}

func isZeroSettingValue(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		if value.Len() == 0 {
			return true
		}
		for i := 0; i < value.Len(); i++ {
			if !isZeroSettingValue(value.Index(i)) {
				return false
			}
		}
		return true
	}
	return value.IsZero()
}
