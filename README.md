# 农业种子萌发时序证据台 (task228-seedgerm)

种子实验人员用来标记每粒种子的萌发阶段，核对温湿度干预与异常停滞关系的后端服务。

## 业务域

- 导入按时间采集的种子图像摘要与环境曲线
- 服务检测胚根出现、叶鞘展开等阶段事件，识别停滞与污染
- 用户可修订阶段边界、加入人工观察并发布试验结果

非 OA / 非游戏 / 非前端页面，属实验观测证据系统。

## 标准命令

```bash
# 构建
CGO_ENABLED=0 GOTOOLCHAIN=local go build ./...

# 静态检查
CGO_ENABLED=0 GOTOOLCHAIN=local go vet ./...

# 测试
CGO_ENABLED=0 GOTOOLCHAIN=local go test ./...

# 自检（创建数据 -> 关闭重开 DB 验证持久化与重启恢复 -> 退出 0）
go run ./cmd/seedgerm --smoke-test

# 启动服务
go run ./cmd/seedgerm --addr :8080 --db ./seedgerm.db
```

## API 入口

全部以 `/api` 为前缀，详见 `internal/httpapi`。

- 试验：`POST /api/trials` `GET /api/trials/:id`
- 种子：`POST /api/trials/:id/seeds` `GET /api/trials/:id/seeds`
- 图像：`POST /api/seeds/:id/images` `GET /api/seeds/:id/images`
- 环境：`POST /api/trials/:id/env` `GET /api/trials/:id/env`
- 阶段：`POST /api/seeds/:id/stages` `POST /api/stages/:id/confirm`
- 复核：`POST /api/seeds/:id/observations` `POST /api/stages/:id/resolve`
- 结果：`POST /api/trials/:id/results` `POST /api/results/:id/publish`
- 自检：`GET /api/selfcheck`
