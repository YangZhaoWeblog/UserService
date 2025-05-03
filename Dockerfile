
# 语法声明
FROM golang:1.23.0 AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 安装必要的构建工具
RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler \
    git \
    make \
    && rm -rf /var/lib/apt/lists/*

# 设置 SSH 配置
RUN mkdir -p /root/.ssh && chmod 700 /root/.ssh
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# 配置 Git 使用 SSH 代替 HTTPS
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

# 设置 Go 私有模块
# ENV GOPRIVATE=github.com/Sider-ai/*

COPY . /src
WORKDIR /src

# Dockerfile snippet
RUN go version

# 调试 SSH 连接
# RUN --mount=type=ssh ssh -T git@github.com


# 合并多个 make 命令到一个 RUN 指令中，并挂载 SSH
RUN --mount=type=ssh \
    make init && \
    make all && \
    make build

# deploy
FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  netbase \
  vim \
  net-tools \
  curl \
  && rm -rf /var/lib/apt/lists/* \
  && apt-get autoremove -y && apt-get autoclean -y

# 从上一步的构建目录中拿取文件
COPY --from=builder /src/bin /app
COPY --from=builder /src/configs /app/configs

WORKDIR /app

EXPOSE 8010
EXPOSE 9010
VOLUME /data/conf

CMD ["./server"]