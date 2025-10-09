package config

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type config struct {
	Db    DBConfig
	Redis RedisConfig
}
