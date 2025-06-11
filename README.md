# FM-Flow 部署自动化平台

![Go Version](https://img.shields.io/badge/Go-1.23.6-blue)
![Vue Version](https://img.shields.io/badge/Vue-3.5-brightgreen)
![License](https://img.shields.io/badge/License-MIT-green)

基于可视化流程编排的智能部署系统，支持跨平台、多语言的分布式部署管理。

## ✨ 核心功能

### 流程引擎
- **可视化编排** - 拖拽式流程设计（基于 LogicFlow 引擎）
- **智能路由** - 自动处理节点依赖关系（见 `flowProcess.go` 的 NextNodes 逻辑）
- **版本控制** - 流程模板版本化管理（参考 `storage.go` 的 MemoryStorage 实现）

### 部署执行
```go
// 节点执行核心逻辑（internal/server/flowProcess.go）
go fp.executeNode(node, server) // 并发执行节点任务
 
