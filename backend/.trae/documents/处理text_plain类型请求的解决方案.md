# 处理text/plain类型请求的解决方案

## 问题分析
Kratos 2.8.0版本的HTTP服务器不支持ContentTypeParser选项，导致无法直接处理text/plain类型的请求。

## 解决方案
使用自定义RequestDecoder来处理text/plain类型的请求，将其转换为application/json类型，然后使用默认的JSON解码器进行解码。

## 实现步骤

1. **修改internal/server/http.go文件**：
   - 移除之前添加的ContentTypeParser选项
   - 添加自定义RequestDecoder，处理text/plain类型的请求

2. **实现自定义RequestDecoder**：
   - 检查请求的Content-Type
   - 如果是text/plain类型，将请求体转换为JSON格式
   - 使用默认的JSON解码器进行解码

3. **更新HTTP服务器配置**：
   - 将自定义RequestDecoder添加到HTTP服务器选项中

4. **测试验证**：
   - 测试修改后的代码，确保text/plain类型的请求能够被正确处理

现在我将按照这个计划开始实现。