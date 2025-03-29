# My Go Server

一个基于Gin框架的简洁Go语言后端服务器。

## 项目结构

```
my-go-server/
├── config/        # 配置相关
├── handlers/      # HTTP请求处理器
├── models/        # 数据模型
├── routes/        # 路由定义
├── utils/         # 工具函数
├── main.go        # 应用入口
└── README.md      # 项目文档
```

## 启动服务

```bash
# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

服务将在 http://localhost:8080 上运行。

## API 端点

- `GET /api/v1/ping` - 健康检查，返回 "pong" 