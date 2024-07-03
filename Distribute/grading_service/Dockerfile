FROM golang:1.21.1-alpine AS builder

MAINTAINER Aurora_Galaxy

WORKDIR /go/src/grading_service

COPY . .

WORKDIR /go/src/grading_service/cmd

# 在构建阶段编译 Go 应用,禁用cgo 使用静态链接
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grading_service .

RUN go build -o grading_service .

# 使用更小的基础镜像 Alpine，作为最终运行镜像
FROM alpine:latest

# 复制编译好的可执行文件到运行镜像
COPY --from=builder /go/src/grading_service/cmd/grading_service  /grading_service

EXPOSE 6000

ENTRYPOINT ["/grading_service"]