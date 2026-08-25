基于 Go 实现的煤矿井下带式输送机集控平台项目，一款后端服务，完成皮带逆煤流启停、瓦斯与温度联锁、制动保护、检修放行与煤流计量的集中控制。

## 构建与运行

使用固定 Go 工具链构建：

    go build -mod=vendor ./...

启动服务：

    go run ./cmd/minebelt -port 8080 -data ./data

健康检查：

    curl http://127.0.0.1:8080/healthz

## Docker

    bash build_benzhi_docker.sh
    docker run --rm -p 8080:8080 minebelt:latest
