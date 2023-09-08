package xstorage

type BindFileType int

const (
	JSON BindFileType = iota
	TOML
)

type KeyValueProperty uint32

const (
	RWLocal KeyValueProperty = 1 << iota
)

type keyValueSaveType uint32

const (
	null keyValueSaveType = iota
	sqlite
)

type KeyValueSetting struct {
	Property KeyValueProperty
	SaveType keyValueSaveType
}

type Setting interface {
	Load() bool
	Save() bool
}
