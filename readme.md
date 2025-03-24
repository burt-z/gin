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
## 使用 Makefile 文件构建
执行命令 make docker
代码需要有变化,代码无变化 image 更新时间不变.

## gin 的 k8s 部署
kubectl apply -f k8s_webook_deployment.yaml
查看是否成功
kubectl get developments

kubectl apply -f k8s_webook_service.yaml
查看是否成功
kubectl get services






需要一个 service,一个 development 管理 3 个 pod

