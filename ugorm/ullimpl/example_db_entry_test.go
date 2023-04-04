package ullimpl

import (
	"context"
	"testing"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/llmuz/yggdrasill/ugorm/config"
)

func TestNewDBHelper(t *testing.T) {
	type UserInfo struct {
		ID        int64
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt time.Time
	}
	var cnf = &config.EngineConfig{}
	var configToml = "" +
		"driver='mysql'\n" +
		"dsn='root:123456@tcp(127.0.0.1:5601)/crok?charset=utf8mb4&parseTime=True&loc=Local'\n" +
		""
	if _, err := toml.Decode(configToml, cnf); err != nil {
		t.Error(err)
		return
	}
	var opts = make([]ConfigOption, 0)

	opts = append(opts, WithSingularTable())
	var helper, err = NewDBHelper(cnf, opts...)

	if err != nil {
		t.Error("error ", err.Error())
		return
	}
	_ = helper.WithContext(context.TODO()).AutoMigrate(&UserInfo{})
	defer helper.WithContext(context.TODO()).Raw("drop table if exist ?", "user_info")

}
