//go:build !k8s

package config

var Config config = config{
	Db: DBConfig{
		DSN: "root:root@tcp(gin-webook-mysql:3309)/webook",
	},
	Redis: RedisConfig{
		Addr: "gin-book-redis:10379",
	},
}
