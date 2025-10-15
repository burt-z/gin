//go:build !k8s

package config

var Config config = config{
	Db: DBConfig{
		DSN: "root:root@tcp(localhost:3309)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
