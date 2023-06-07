package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/llmuz/yggdrasill/ull"
	"github.com/llmuz/yggdrasill/ull/config"
	"github.com/llmuz/yggdrasill/ull/zapimpl"
)

type H struct {
	Key string
}

func (c *H) Levels() (lvs []ull.Level) {
	return []ull.Level{ull.DebugLevel, ull.InfoLevel, ull.ErrorLevel, ull.WarnLevel, ull.PanicLevel}
}

func (c *H) Fire(e ull.Entry) (err error) {
	err = e.AppendField(ull.Any(c.Key, e.Context().Value("trace_id")))
	return err
}

var (
	pf          = flag.String("conf", "ull/examples/zaplog/config.toml", "配置文件")
	cfg         config.LogConfig
	complexData = make(map[string]interface{})
)

func main() {
	flag.Parse()
	if _, err := toml.DecodeFile(*pf, &cfg); err != nil {
		panic(err)
	}

	logger, err := zapimpl.NewZapLogger(&cfg)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	// 移除日志文件
	//defer os.RemoveAll(filepath.Dir(cfg.GetFileName()))
	n := time.Now()
	ctx := context.WithValue(context.TODO(), "trace_id", md5.New().Sum([]byte(time.Now().String())))
	helper := zapimpl.NewHelper(logger, zapimpl.AddHook(&H{Key: "hello"}), zapimpl.AddHook(&H{Key: "900x"}), zapimpl.AddHook(&H{Key: "trace_id"}))

	fmt.Println("start ", n)
	helper.WithContext(ctx).Errorf("hello %s %s", "你", "好")
	for i := 0; i < 100; i++ {
		helper.WithContext(ctx).Debugf("hello %#v %#v", ull.Any("hello", complexData), ull.Any("now", time.Now()))
		helper.WithContext(ctx).Infof("hello %#v %#v", ull.Any("hello", complexData), ull.Any("now", time.Now()))
		helper.WithContext(ctx).Warnf("hello %#v %#v", ull.Any("hello", complexData), ull.Any("now", time.Now()))
		helper.WithContext(ctx).Errorf("hello %#v %#v", ull.Any("hello", complexData), ull.Any("now", time.Now()))
	}
	fmt.Println("end ", time.Now().Sub(n))
}
