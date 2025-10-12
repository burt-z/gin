# 源码地址 
https://gitee.com/wei-jie-zhu-code/geektime-basic-go/tree/master/webook


# 层级
internal
- domain 领域对象,
- service 领域服务,业务的完整处理过程
- repository 领域存储,存储数据的抽象


# gin 本项目
- docker compose up/down
- go run main.go

# webook-fe 前端项目
- npm run dev

# 注意事项
CreateUser(db) 会创建数据表

#安装 kubectl 工具
- brew install kubectl
- kubectl version --client 查看版本

#使用 k8s部署

## 生成 linux 二进制文件
GOOS=linux GOARCH=arm go build -o webook .
- GOOS:平台
- GOARCH:架构
- -o 输出文件 
- webook 输出文件名
- . 当前目录
## 将 webook 塞进去

## 构建镜像命令
docker build -t webook:v1 .
## 命令执行失败
### 如果命令执行失败,在终端命令行内执行 docker pull ubuntu:latest
### 在 docker desktop 的 docker engine 里面增加配置加速器
{
"registry-mirrors": [
"https://docker.mirrors.ustc.edu.cn",
"https://hub-mirror.c.163.com",
"https://mirror.baidubce.com",
"https://registry.docker-cn.com"
],
"max-concurrent-downloads": 3,
"max-download-attempts": 5
}

## 使用 Makefile 文件构建
执行命令 make docker
代码需要有变化,代码无变化 image 更新时间不变.

## gin 的 k8s 部署
项目根目录下,gin>>
kubectl apply -f k8s_webook_deployment.yaml
查看是否成功
kubectl get deployments

启动 service
kubectl apply -f k8s_webook_service.yaml
查看是否成功
kubectl get services

# 查看
kubectl get pv
kubectl get pvc

# 删除
kubectl delete pv 名字
kubectl delete pvc 名字

# k8s 部署 mysql
kubectl apply -f  k8s_gin_webook_mysql_service.yaml

kubectl apply -f  k8s_gin_webook_mysql_deployment.yaml
kubectl apply -f  k8s_gin_webook_mysql_pvc.yaml
kubectl apply -f  k8s_gin_webook_mysql_pv.yaml

# k8s 部署 redis
kubectl apply -f k8s_gin_webook_redis_service.yaml
kubectl apply -f k8s_gin_webook_redis_deployment.yaml

# 部署 ingress-nginx
kubectl apply -f k8s_gin_webook_ingress_nginx.yaml
里面有个 host: live.webook.com 的配置,需要配置 host文件,配合解析
sudo vim /etc/hosts
增加下面的配置
127.0.0.1 live.webook.com



# 本地连接 k8s 里面的 redis
redis-cli -h localhost -p 30003


# 安装 ingress-nginx 的区别
本地
--install ingress-nginx ingress-nginx/ingress-nginx \

教程里面
--install ingress-nginx ingress-nginx \


ide 里面访问 service 里面暴露的 port,targetPort 对应 deployment里面的containers下面的 port

# //go:build k8s 的作用,有两个变量 Config,加了 go:build 不会冲突
go build -tags=k8s

# 生成 mock文件 在极客目录下  /jike/gin/
mockgen -source=gin/internal/service/code.go -package=svcmock -destination=gin/internal/service/mocks/code.mock.go
mockgen -source=gin/internal/service/user.go -package=svcmock -destination=gin/internal/service/mocks/user.mock.go






需要一个 service,一个 development 管理 3 个 pod

