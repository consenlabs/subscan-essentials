package configs

import (
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/container/pool"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/itering/subscan/util"
	"strings"
)

type (
	MysqlConf struct {
		Conf struct {
			Host string
			User string
			Pass string
			DB   string
		}
		Api  *sql.Config
		Task *sql.Config
		Test *sql.Config
	}
	RedisConf struct {
		Config *redis.Config
		DbName int
	}
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (dc *MysqlConf) MergeConf() {
	// checkErr(paladin.Get("mysql.toml").UnmarshalTOML(dc)) // todo 通过环境变量的方式进行配置
	dc.Api = &sql.Config{}
	dc.Task = &sql.Config{}
	dc.mergeEnvironment()
}

func (rc *RedisConf) MergeConf() {
	// checkErr(paladin.Get("redis.toml").UnmarshalTOML(rc)) // todo 同上
	rc.Config = &redis.Config{
		Config: &pool.Config{
			Active:      10,
			Idle:        10,
			IdleTimeout: 1000000000,
			WaitTimeout: 1000000000,
		},
		Name:         "substrate",
		Proto:        "tcp",
		DialTimeout:  1000000000,
		ReadTimeout:  1000000000,
		WriteTimeout: 1000000000,
		SlowLog:      0,
	}
	rc.mergeEnvironment()
}

func (dc *MysqlConf) mergeEnvironment() {

	// dbHost := util.GetEnv("MYSQL_HOST", dc.Conf.Host)
	// dbUser := util.GetEnv("MYSQL_USER", dc.Conf.User)
	// dbPass := util.GetEnv("MYSQL_PASS", dc.Conf.Pass)
	// dbName := util.GetEnv("MYSQL_DB", dc.Conf.DB)
	// dc.Api.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName) + dc.Api.DSN
	// dc.Task.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName) + dc.Task.DSN

	// TODO 为了适应ops 那边的配置文件格式
	dc.Api.DSN = util.GetEnv("MYSQL", "")
	dc.Task.DSN = util.GetEnv("MYSQL", "")
}

func (rc *RedisConf) mergeEnvironment() {
	// rc.Config.Addr = util.GetEnv("REDIS_ADDR", rc.Config.Addr)
	// rc.DbName = util.StringToInt(util.GetEnv("REDIS_DATABASE", "0"))

	// todo 为了适应ops 那边的配置文件格式
	redisUrl := util.GetEnv("REDIS", "") // 环境变量中的类似：redis://10.0.0.12:6379/10
	// 拆开 address 和 db
	spl := strings.Split(strings.TrimPrefix(redisUrl, "redis://"), "/")
	rc.Config.Addr = spl[0]
	rc.DbName = util.StringToInt(spl[1])
}
