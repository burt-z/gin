.PHONY: docker
# 这个是一个别名
docker:
	@rm gin-webook || true
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o gin-webook .
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -tags=k8s -o gin-webook .
	@docker rmi -f gin-webook:v0.0.1
	@docker build -t gin-webook:v0.0.1 .



#.PHONY: mock
## 这个是一个别名
#mock:
#	@mockgen -package=redismock -destination=internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
#
