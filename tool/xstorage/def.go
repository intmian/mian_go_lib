package xstorage

type BindFileType int

const (
	JSON BindFileType = iota
	TOML
)

type KeyValueProperty uint32

const (
	// RWLock 需加入读写锁
	RWLock KeyValueProperty = 1 << iota
	// NoCache 不缓存
	NoCache
	// JustCache 只缓存
	JustCache
)

type keyValueSaveType uint32

const (
	null keyValueSaveType = iota
	sqlite
	file
)

type KeyValueSetting struct {
	Property KeyValueProperty
	SaveType keyValueSaveType
}

type Setting interface {
	Load() bool
	Save() bool
}

type ValueType int

const (
	STRING ValueType = iota
	INT
	FLOAT
	BOOL
	SLICE_STRING
	SLICE_INT
	SLICE_FLOAT
	SLICE_BOOL
)
