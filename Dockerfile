# 基础镜像
FROM ubuntu:latest

# 打包进来这个镜像,放到这个目录,随意换
COPY webook /app/webook

WORKDIR /app

# 执行命令
ENTRYPOINT ["/app/webook"]