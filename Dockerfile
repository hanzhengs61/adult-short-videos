# ========================================
# 多阶段构建 Dockerfile
# ========================================

# 第一阶段：构建
FROM golang:1.25.5-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要工具
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 第二阶段：运行
FROM alpine:latest

# 安装 ca-certificates（HTTPS 需要）
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# 创建日志目录
RUN mkdir -p /root/logs

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]