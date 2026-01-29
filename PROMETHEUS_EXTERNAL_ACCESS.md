# Prometheus 外部访问配置指南

本文档说明如何从**另一台服务器上的 Prometheus** 访问 K8s 集群内 backend 服务的 metrics 和 pprof 端点。

## 1. 已暴露的端点

### Metrics 端点
- **路径**: `/metrics`
- **端口**: `8000`
- **访问地址**: `http://<pod-ip>:8000/metrics` 或通过 Service: `http://backend.aichat.svc.cluster.local:8000/metrics`

### pprof 性能分析端点
- **路径**: `/debug/pprof/*`
- **端口**: `8000`
- **常用端点**:
  - `/debug/pprof/` - 概览页面
  - `/debug/pprof/heap` - 堆内存分析
  - `/debug/pprof/profile` - CPU 性能分析（默认30秒采样）
  - `/debug/pprof/goroutine` - Goroutine 分析
  - `/debug/pprof/block` - 阻塞分析
  - `/debug/pprof/mutex` - 互斥锁分析

## 2. 外部 Prometheus 访问方案

### 方案 A：通过 NodePort Service（推荐用于内网环境）

在 `k8s-deployment.yaml` 中创建一个 NodePort Service，将 metrics 端口暴露到节点 IP：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-metrics
  namespace: aichat
  labels:
    app: backend
spec:
  type: NodePort
  selector:
    app: backend
  ports:
  - name: metrics
    port: 8000
    targetPort: 8000
    nodePort: 30800  # 外部访问端口（30000-32767 范围）
```

**Prometheus 配置** (`prometheus.yml`):
```yaml
scrape_configs:
  - job_name: 'aichat-backend'
    static_configs:
      - targets: ['<k8s-node-ip>:30800']  # 替换为你的 K8s 节点 IP
        labels:
          namespace: 'aichat'
          service: 'backend'
```

### 方案 B：通过 Ingress 暴露（不推荐，存在安全风险）

**注意**: metrics 可能包含敏感信息，不建议直接暴露到公网。如果必须使用，请配置 IP 白名单或 Basic Auth。

在 `k8s-deployment.yaml` 的 Ingress 中添加：

```yaml
  - host: metrics.aihelper.chat  # 使用独立域名
    http:
      paths:
      - path: /metrics
        pathType: Prefix
        backend:
          service:
            name: backend
            port:
              number: 8000
```

**Prometheus 配置**:
```yaml
scrape_configs:
  - job_name: 'aichat-backend'
    static_configs:
      - targets: ['metrics.aihelper.chat']  # 通过域名访问
        labels:
          namespace: 'aichat'
          service: 'backend'
    scheme: https  # 如果配置了 TLS
    tls_config:
      insecure_skip_verify: true  # 仅用于测试，生产环境应配置正确证书
```

### 方案 C：使用 Prometheus Federation（推荐用于生产环境）

在 K8s 集群内部署一个轻量级 Prometheus（或使用 Prometheus Operator），然后通过 Federation 将数据推送到外部 Prometheus。

**集群内 Prometheus 配置** (`prometheus-in-cluster.yml`):
```yaml
scrape_configs:
  - job_name: 'aichat-backend'
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - aichat
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: backend
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        replacement: $1:8000
```

**外部 Prometheus 配置** (`prometheus-external.yml`):
```yaml
scrape_configs:
  - job_name: 'federate-aichat'
    honor_labels: true
    metrics_path: '/federate'
    params:
      'match[]':
        - '{job="aichat-backend"}'
    static_configs:
      - targets: ['<cluster-internal-prometheus>:9090']
```

### 方案 D：使用 kubectl port-forward（仅用于测试/调试）

临时端口转发，适合快速测试：

```bash
# 转发 backend Service 的 8000 端口到本地 18000
kubectl port-forward -n aichat svc/backend 18000:8000

# 然后在外部 Prometheus 配置中
scrape_configs:
  - job_name: 'aichat-backend-local'
    static_configs:
      - targets: ['localhost:18000']  # 仅在同一台机器上有效
```

## 3. 推荐的完整配置（方案 A + 方案 C 结合）

### 步骤 1: 添加 NodePort Service

在 `k8s-deployment.yaml` 中添加：

```yaml
---
# Metrics 监控专用 Service（NodePort 类型，仅内网访问）
apiVersion: v1
kind: Service
metadata:
  name: backend-metrics
  namespace: aichat
  labels:
    app: backend
    component: metrics
spec:
  type: NodePort
  selector:
    app: backend
  ports:
  - name: metrics
    port: 8000
    targetPort: 8000
    nodePort: 30800  # 外部访问端口
```

### 步骤 2: 配置外部 Prometheus

在外部服务器的 `prometheus.yml` 中添加：

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  # aichat backend metrics
  - job_name: 'aichat-backend'
    scrape_interval: 10s
    static_configs:
      - targets: ['<k8s-node-ip>:30800']  # 替换为实际 K8s 节点 IP
        labels:
          cluster: 'production'
          namespace: 'aichat'
          service: 'backend'
    metrics_path: '/metrics'
    
  # 如果需要抓取多个 Pod 的指标（绕过 Service 负载均衡）
  - job_name: 'aichat-backend-pods'
    kubernetes_sd_configs:
      - role: pod
        api_server: 'https://<k8s-api-server>:6443'  # K8s API Server 地址
        tls_config:
          ca_file: '/path/to/ca.crt'
          cert_file: '/path/to/client.crt'
          key_file: '/path/to/client.key'
    relabel_configs:
      - source_labels: [__meta_kubernetes_namespace]
        action: keep
        regex: aichat
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: backend
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_ip, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        replacement: $1:8000
```

### 步骤 3: 验证配置

```bash
# 1. 检查 NodePort Service 是否创建成功
kubectl get svc -n aichat backend-metrics

# 2. 测试 metrics 端点是否可访问（从外部服务器）
curl http://<k8s-node-ip>:30800/metrics

# 3. 测试 pprof 端点
curl http://<k8s-node-ip>:30800/debug/pprof/

# 4. 重启 Prometheus 并检查 targets
# 访问 Prometheus UI: http://<prometheus-server>:9090/targets
```

## 4. 安全建议

1. **防火墙规则**: 限制 NodePort 端口（30800）仅允许 Prometheus 服务器 IP 访问
2. **网络策略**: 使用 K8s NetworkPolicy 限制 metrics 端点的访问来源
3. **TLS 加密**: 如果通过公网访问，建议使用 Ingress + TLS
4. **认证授权**: 考虑在 metrics 端点添加 Basic Auth 或 Token 认证

## 5. pprof 使用示例

### 通过 NodePort 访问 pprof

```bash
# CPU 性能分析（采样30秒）
go tool pprof http://<k8s-node-ip>:30800/debug/pprof/profile?seconds=30

# 堆内存分析
go tool pprof http://<k8s-node-ip>:30800/debug/pprof/heap

# Goroutine 分析
go tool pprof http://<k8s-node-ip>:30800/debug/pprof/goroutine

# 生成可视化报告
go tool pprof -http=:8080 http://<k8s-node-ip>:30800/debug/pprof/heap
```

### 通过 Ingress 访问（如果配置了）

```bash
# 访问 pprof 概览页面
curl http://api.aihelper.chat/debug/pprof/

# CPU 分析
go tool pprof http://api.aihelper.chat/debug/pprof/profile
```

## 6. 故障排查

### 问题 1: Prometheus 无法连接
- 检查防火墙是否开放 NodePort 端口
- 验证 K8s 节点 IP 是否正确
- 检查 Service 的 selector 是否匹配 Pod labels

### 问题 2: metrics 端点返回 404
- 确认 backend 服务已正确部署
- 检查 Pod 日志确认服务正常启动
- 验证 `/metrics` 路由是否已注册

### 问题 3: 抓取到的指标为空
- 检查 Pod annotations 是否正确配置
- 验证 metrics 中间件是否已启用
- 查看 Prometheus 的 `/targets` 页面确认抓取状态
