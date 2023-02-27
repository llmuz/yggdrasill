package ugorm

import (
	"fmt"

	"gorm.io/gorm"
)

type DialectFunc func(dsn string) gorm.Dialector

var (
	plugins = make(map[string]DialectFunc)
)

func GetDialect(driver string) (dialect DialectFunc, err error) {
	var ok bool
	if dialect, ok = plugins[driver]; !ok {
		return nil, fmt.Errorf("drive:%s not register", driver)
	}
	return dialect, nil
}

func Register(driver string, dialect DialectFunc) (err error) {
	if _, ok := plugins[driver]; ok {
		return fmt.Errorf("dialect:%s has register", driver)
	}

	plugins[driver] = dialect

	return nil
}
