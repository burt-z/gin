# 基础镜像
FROM ubuntu:latest

# 打包进来这个镜像,放到这个目录,随意换
COPY gin-webook /app/gin-webook

WORKDIR /app

# 执行命令
ENTRYPOINT ["/app/gin-webook"]