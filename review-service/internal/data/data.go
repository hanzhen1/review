package data

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"review-service/internal/conf"
	"review-service/internal/data/query"
	"strings"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewData, NewReviewRepo, NewESClient, NewRedisClient)

// Data .
type Data struct {
	//db *gorm.DB
	query *query.Query
	log   *log.Helper
	es    *elasticsearch.TypedClient //"github.com/elastic/go-elasticsearch/v8"
	rdb   *redis.Client
}

// NewData .
func NewData(db *gorm.DB, esClient *elasticsearch.TypedClient, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	//非常重要! 为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)
	return &Data{
		query: query.Q,
		log:   log.NewHelper(logger),
		es:    esClient,
		rdb:   rdb,
	}, cleanup, nil
}

// NewESClient ES Client构造函数
func NewESClient(conf *conf.Elasticsearch) (*elasticsearch.TypedClient, error) {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: conf.Addresses,
	}
	// 创建客户端连接
	return elasticsearch.NewTypedClient(cfg)
}
func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.GetDriver()) { //strings.ToLower 转为小写
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.GetSource()))
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.GetSource()))
	}
	return nil, errors.New("connect DB fail unsupported db driver")
}
func NewRedisClient(cfg *conf.Data) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		ReadTimeout:  cfg.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: cfg.Redis.WriteTimeout.AsDuration(),
	})
}
