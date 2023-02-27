package ugorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	_ = Register("mysql", func(dsn string) gorm.Dialector {
		return mysql.Open(dsn)
	})
}
