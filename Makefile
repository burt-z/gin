.PHONY: docker
# 这个是一个别名
docker:
	@rm gin-webook || true
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -tags=k8s -o gin-webook .
	@docker rmi -f gin-webook:v0.0.1
	@docker build -t gin-webook:v0.0.1 .

# 视频的是在 jike 这个目录里面,和 gin 同级,测试的时候新建一个 MakeFile文件就行
#.PHONY: mock
 ## 这个是一个别名
 #mock:
 #	@mockgen -source=gin/internal/service/code.go -package=svcmock -destination=gin/internal/service/mocks/code.mock.go
 #	@mockgen -source=gin/internal/service/user.go -package=svcmock -destination=gin/internal/service/mocks/user.mock.go