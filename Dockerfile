FROM golang:1.24.2-alpine AS builder

LABEL authors="mundo"

# 设置工作目录
WORKDIR /mundo/chat

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 安装 Git
RUN apk update && apk add --no-cache git

# 设置 GitHub 访问令牌用于认证
ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

RUN go mod tidy

# 下载 Go 项目的依赖
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 编译 Go 项目
RUN go build -o main .

# 使用更小的 Alpine 镜像作为运行时镜像
FROM alpine:latest

# 从构建阶段复制编译好的二进制文件到运行时镜像
COPY --from=builder /mundo/chat/main /main