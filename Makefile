.PHONY: mock
# 这个是一个别名
mock:
	@mockgen -package=redismock -destination=internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable

