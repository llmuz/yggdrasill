package ullimpl

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/llmuz/yggdrasill/ugorm"
	"github.com/llmuz/yggdrasill/ugorm/config"
	"github.com/llmuz/yggdrasill/ull"
)

type ConfigOption func(c *gorm.Config)

func WithName(strategy schema.Namer) ConfigOption {
	return func(c *gorm.Config) {
		c.NamingStrategy = strategy
	}
}

func WithLogger(l logger.Interface) ConfigOption {
	return func(c *gorm.Config) {
		c.Logger = l
	}
}

func NewLogger(log ull.Helper, conf *config.OrmLogConfig, optHooks ...ugorm.Hook) logger.Interface {
	writer := NewWriter(log, optHooks...)
	loggerConfig := buildLoggerConfig(conf)
	loggerInterface := NewLoggerInterface(writer, loggerConfig)
	return loggerInterface
}

func WithSingularTable() ConfigOption {
	return WithName(schema.NamingStrategy{SingularTable: true})
}

type dbEntry struct {
	db *gorm.DB
}

func (c *dbEntry) WithContext(ctx context.Context) (db *gorm.DB) {
	return c.db.WithContext(ctx)
}

func NewDBHelper(conf *config.EngineConfig, opts ...ConfigOption) (helper ugorm.DBHelper, err error) {
	var dialect ugorm.DialectFunc
	if dialect, err = ugorm.GetDialect(conf.GetDriver()); err != nil {
		return nil, err
	}

	var db *gorm.DB
	var dbConfig gorm.Config
	dbConfig.SkipDefaultTransaction = conf.GetSkipDefaultTransaction()
	dbConfig.DisableAutomaticPing = conf.GetDisableAutomaticPing()
	dbConfig.PrepareStmt = conf.GetPrepareStmt()
	if conf.CreateBatchSize != nil {
		dbConfig.CreateBatchSize = int(conf.GetCreateBatchSize())
	}

	for _, fn := range opts {
		fn(&dbConfig)
	}

	if db, err = gorm.Open(dialect(conf.GetDsn()), &dbConfig); err != nil {
		return nil, err
	}

	var sqlDB *sql.DB
	if sqlDB, err = db.DB(); err != nil {
		return nil, err
	}
	initDBConnPool(sqlDB, conf.GetConnPool())

	return &dbEntry{db: db}, err
}

func initDBConnPool(db *sql.DB, conf *config.ConnPool) {
	if conf == nil || db == nil {
		return
	}

	if duration, err := time.ParseDuration(conf.GetMaxIdleTime()); err != nil {
		err = nil

	} else {
		fmt.Println("max idle time ", duration)
		db.SetConnMaxIdleTime(duration)

	}

	if duration, err := time.ParseDuration(conf.GetMaxLifeTime()); err != nil {
		db.SetConnMaxLifetime(duration)
	} else {
		fmt.Println("")
	}

	if conf.GetMaxIdleConn() != 0 && conf.GetMaxOpenConn() > 0 {
		db.SetMaxIdleConns(int(conf.GetMaxIdleConn()))
	}

	if conf.GetMaxOpenConn() != 0 && conf.GetMaxOpenConn() > 0 {
		db.SetMaxOpenConns(int(conf.GetMaxOpenConn()))
	}
}

func buildLoggerConfig(conf *config.OrmLogConfig) logger.Config {
	c := logger.Config{}
	var err error
	if c.SlowThreshold, err = time.ParseDuration(conf.GetSlowThreshold()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "parser slow threshold %s", err)
		_, _ = fmt.Fprintf(os.Stderr, "use 1s as slowThreshold")
		c.SlowThreshold = time.Second
	}
	c.LogLevel = ParserLevel(conf.GetLogLevel())
	c.Colorful = conf.GetColorful()
	c.IgnoreRecordNotFoundError = conf.GetIgnoreRecordNotFoundError()
	return c
}
