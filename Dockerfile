# 多阶段构建：前端构建阶段
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装前端依赖
RUN npm ci --only=production=false

# 复制前端源代码
COPY frontend/ ./

# 构建前端（生成静态文件）
RUN npm run build

# 多阶段构建：后端构建阶段
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制 Go 依赖文件
COPY backend/go.mod backend/go.sum ./

# 下载 Go 依赖
RUN go mod download

# 复制后端源代码
COPY backend/ ./

# 编译后端（静态链接，生成单个可执行文件）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o explore \
    explore.go

# 最终运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 从后端构建阶段复制可执行文件
COPY --from=backend-builder /app/backend/explore /app/explore

# 从后端构建阶段复制配置文件
COPY --from=backend-builder /app/backend/etc /app/etc

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /app/frontend/dist /app/frontend

# 创建日志目录
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 8877

# 设置时区
ENV TZ=Asia/Shanghai

# 启动后端服务（后端会提供静态文件服务）
CMD ["/app/explore", "-f", "/app/etc/explore.yaml"]

