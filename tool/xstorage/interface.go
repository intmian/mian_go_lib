package xstorage

type IMgr interface {
	Get(key string) (*ValueUnit, error)
	Set(key string, value *ValueUnit) error
}
