package xstorage

type IDBCore interface {
	Get(key string) (bool, *ValueUnit, error)
	Set(key string, value *ValueUnit) error
	GetAll() (map[string]*ValueUnit, error)
	Delete(key string) error
}
