# event-go
An event system based on Go.

## 功能特性
- 事件发布/订阅机制
- 线程安全的并发处理
- 支持三种发布策略：
  - DiscardNew（默认）：队列满时丢弃新事件
  - Block：阻塞式发布
  - Serial：串行顺序处理
- 内置对象池减少GC压力
- 完善的panic恢复机制
- 支持动态订阅管理

## 安装
```bash
go get -u github.com/zhoudm1743/event-go
```

## 快速开始
```go
package main

import (
	"go.uber.org/zap"
	"event_go"
)

func main() {
	logger, _ := zap.NewProduction()
	es := event_go.NewEventService(logger)
	
	// 订阅事件
	unsubscribe := es.Subscribe("user.created", func(e event_go.Event) {
		logger.Info("收到事件", zap.String("name", e.Name), zap.Any("data", e.Data))
	})
	defer unsubscribe()

	// 发布事件
	es.Publish(event_go.Event{
		Name: "user.created",
		Data: map[string]interface{}{"id": 1001, "name": "张三"},
	})
}
```

## 配置选项
```go
// 使用自定义配置
es := NewEventService(logger,
	WithWorkerCount(8),           // 设置工作协程数
	WithChannelSize(5000),        // 事件通道容量
	WithPublishStrategy(Serial),  // 使用串行策略
)
```

## 串行模式说明
当使用WithSerialStrategy时：
1. 事件按接收顺序严格串行处理
2. 每个事件处理完成后才会处理下一个
3. 适用场景：
   - 需要严格顺序保证
   - 资源敏感型操作
   - 状态依赖型处理

```go
// 串行模式示例
es := NewEventService(logger, WithSerialStrategy())
```

## 注意事项
1. 事件对象在发布后会被回收，订阅者不应保留引用
2. 串行模式性能低于并行模式，需根据场景选择
3. 建议使用zap日志库以获得完整调试信息
4. 发布阻塞时间过长时考虑调整通道容量
5. 使用前请调用logger.Sync()确保日志落盘
```