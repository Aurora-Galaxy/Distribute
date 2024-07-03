FROM golang:1.21.1-alpine AS builder

MAINTAINER Aurora_Galaxy

WORKDIR /go/src/register_service

COPY . .

WORKDIR /go/src/register_service/cmd

# 在构建阶段编译 Go 应用,禁用cgo 使用静态链接
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o register_service .

RUN go build -o register_service .

# 使用更小的基础镜像 Alpine，作为最终运行镜像
FROM alpine:latest

# 复制编译好的可执行文件到运行镜像
COPY --from=builder /go/src/register_service/cmd/register_service  /register_service

EXPOSE 3000

#CMD ["./register_service"]
ENTRYPOINT ["/register_service"]