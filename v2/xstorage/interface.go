package xstorage

type IDBCore interface {
	Get(key string, rec *ValueUnit) (bool, error)
	Set(key string, value *ValueUnit) error
	GetAll() (map[string]*ValueUnit, error)
	Delete(key string) error
}

type IFileCore interface {
	GetAll() (map[string]*ValueUnit, error)
	SaveAll(data map[string]*ValueUnit) error
}
