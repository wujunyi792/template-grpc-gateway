package mysql

import (
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"pinnacle-primary-be/core/logx"
	"pinnacle-primary-be/core/middleware/gormv2"
)

type (
	Orm struct {
		Host     string
		Port     string
		User     string
		Pass     string
		Database string
		Debug    bool
		Trace    bool

		Conf *gorm.Config
		*gorm.DB
	}
	Option func(r *Orm)
)

func (r *Orm) GetOrm() *gorm.DB {
	return r.DB
}

func (r *Orm) OrmConnectionUpdate(conf OrmConf) *Orm {
	orm, err := NewMysqlOrm(conf)
	if err != nil {
		return r
	}
	return orm
}

func MustNewMysqlOrm(conf OrmConf, opts ...Option) *Orm {
	orm, err := NewMysqlOrm(conf, opts...)
	if err != nil {
		logx.Errorw("fail to load mysql orm", logx.LogField{Key: "err", Value: err})
		os.Exit(1)
	}
	return orm
}

func NewMysqlOrm(conf OrmConf, opts ...Option) (*Orm, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	opts = append([]Option{WithAddr(conf.Host, conf.Port)}, opts...)
	opts = append([]Option{WithAuth(conf.User, conf.Pass)}, opts...)
	opts = append([]Option{WithTrace(conf.Trace)}, opts...)
	opts = append([]Option{WithDBName(conf.Database)}, opts...)
	opts = append([]Option{WithDebug(conf.Debug)}, opts...)

	return newOrm(opts...)
}

func WithGormConf(conf *gorm.Config) Option {
	return func(r *Orm) {
		r.Conf = conf
	}
}

func WithTrace(trace bool) Option {
	return func(r *Orm) {
		r.Trace = trace
	}
}

func WithDebug(debug bool) Option {
	return func(r *Orm) {
		r.Debug = debug
	}
}

func WithAddr(host, port string) Option {
	return func(r *Orm) {
		r.Host = host
		r.Port = port
	}
}

func WithAuth(user, pass string) Option {
	return func(r *Orm) {
		r.Pass = pass
		r.User = user
	}
}

func WithDBName(db string) Option {
	return func(r *Orm) {
		r.Database = db
	}
}

func newOrm(opts ...Option) (*Orm, error) {
	m := &Orm{}
	for _, opt := range opts {
		opt(m)
	}
	conf := m.Conf
	if conf == nil {
		conf = &gorm.Config{
			Logger: gormv2.NewLogger(),
		}
	}
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Pass, m.Host, m.Port, m.Database)
	db, err := gorm.Open(mysql.Open(dsn), conf)
	if m.Trace {
		if _err := db.Use(otelgorm.NewPlugin()); _err != nil {
			// 用不了拉倒
			logx.Error(_err)
		}
	}
	if m.Debug {
		db = db.Debug()
	}
	m.DB = db
	return m, err
}
