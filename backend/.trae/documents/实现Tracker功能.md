# 实现Tracker功能

## 1. 定义Proto文件
在`api/base/tracker.proto`中定义服务和消息类型，包括：
- Tracker服务，包含Batch和List两个接口
- Tracker消息类型，对应SysTracker模型
- BatchRequest和BatchReply消息类型，用于批量新增
- ListRequest和ListReply消息类型，用于分页查询

## 2. 生成Proto代码
使用make命令生成proto对应的go代码

## 3. 实现Biz层
- 在`internal/biz/base/`目录下创建`tracker.go`文件
- 定义TrackerRepo接口，包含BatchCreate和List方法
- 实现TrackerUseCase，调用data层接口

## 4. 实现Data层
- 在`internal/data/`目录下创建`tracker.go`文件
- 实现TrackerRepo接口，与数据库交互
- 实现批量创建和分页查询功能

## 5. 实现Service层
- 在`internal/service/base/`目录下创建`tracker.go`文件
- 实现Tracker服务的gRPC和HTTP接口
- 调用biz层处理业务逻辑

## 6. 更新依赖注入
- 更新`cmd/backend/wire.go`文件，注入tracker相关的依赖

## 7. 测试验证
- 运行服务，测试两个接口是否正常工作

现在我将按照这个计划开始实现Tracker功能。