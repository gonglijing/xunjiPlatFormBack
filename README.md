# SagooIOT

<div align="center">

![GoFrame](https://img.shields.io/badge/goframe-2.9-green)
![Go Version](https://img.shields.io/badge/go-1.23.0+-blue)
![License](https://img.shields.io/badge/license-GPL3.0-success)

**A lightweight enterprise-grade IoT platform developed in Go**

[English](README.md) | [中文](README_ZH.md)

</div>

---

## 简介

SagooIOT 是一个基于 Go 语言开发的轻量级企业级物联网平台，提供完整的物联网接入、管理、分析和应用解决方案，支持跨平台独立部署或分布式部署。

### 核心特性

- **快速部署** - 开箱即用，分钟级启动完整物联网平台
- **前后端分离** - GoFrame 2.9 + Vue 3，架构清晰易维护
- **多协议支持** - TCP、UDP、HTTP、WebSocket、MQTT、CoAP、OPC UA、Modbus 等
- **插件驱动** - 独特的热插拔插件系统，支持 C/C++、Python、Go 多语言开发
- **高性能数据处理** - 集成 TDengine 时序数据库，支持百万级数据点秒级处理
- **边缘计算** - 支持离线部署、本地规则执行

### 快速信息

| 项目 | 信息 |
|------|------|
| 前端项目 | [sagooiot-ui](https://github.com/sagoo-cloud/sagooiot-ui) |
| 默认账号 | admin / admin123456 |
| 开源协议 | GPL-3.0 |

---

## 技术栈

### 后端

| 技术 | 用途 |
|------|------|
| Go 1.23+ | 核心语言 |
| GoFrame 2.9 | Web 框架 |
| MySQL/PostgreSQL | 关系数据库 |
| Redis | 缓存和消息队列 |
| TDengine/InfluxDB | 时序数据库 |
| MQTT | 物联网协议 |
| Casbin | 权限管理 |
| gToken | 会话管理 |
| Asynq | 任务队列 |
| Gorilla WebSocket | 实时推送 |

### 前端

| 技术 | 用途 |
|------|------|
| Vue 3.x | 前端框架 |
| Element Plus | UI 组件库 |
| TypeScript 4.0+ | 类型语言 |
| Vite 2.0+ | 构建工具 |
| Pinia | 状态管理 |
| Axios | HTTP 客户端 |

---

## 项目结构

```
sagooiot/
├── main.go                    # 程序入口
├── api/                       # API 接口定义
│   └── v1/
├── command/                   # 命令行和服务器配置
│   ├── cmd.go                 # 主命令入口
│   ├── server.go              # HTTP 服务器配置
│   ├── initfunc.go            # 初始化函数
│   └── router/                # 路由配置
│       ├── system.go           # 系统模块路由
│       ├── iot.go              # IoT 模块路由
│       └── analysis.go         # 分析统计路由
├── internal/                  # 核心业务代码
│   ├── controller/            # 控制器层 (API 路由)
│   │   ├── alarm/             # 告警管理
│   │   ├── analysis/          # 分析统计
│   │   ├── device/            # 设备管理
│   │   ├── network/           # 网络通道
│   │   ├── notice/            # 通知服务
│   │   ├── oauth/             # OAuth 认证
│   │   ├── product/           # 产品管理
│   │   ├── system/            # 系统管理
│   │   ├── tdengine/          # 时序数据库
│   │   └── common/            # 公共控制器
│   ├── service/               # 业务逻辑层
│   ├── dao/                   # 数据访问层
│   ├── logic/                 # 业务逻辑
│   ├── model/                 # 数据模型
│   ├── mqtt/                  # MQTT 消息处理
│   ├── tasks/                 # 定时任务
│   ├── workers/               # 工作池
│   ├── websocket/             # WebSocket
│   ├── sse/                   # Server-Sent Events
│   └── queues/                # 消息队列
├── network/                   # 网络协议层
│   ├── network.go             # 网络模块入口
│   ├── core/                  # 核心协议处理
│   │   ├── server/            # TCP/UDP 服务器
│   │   ├── tunnel/            # 隧道管理
│   │   ├── device/            # 设备管理
│   │   └── router.go           # 协议路由
│   ├── model/                 # 网络模型
│   └── events/                # 事件处理
├── pkg/                       # 工具包
│   ├── response/              # 响应封装
│   ├── cache/                 # 缓存
│   ├── dcache/                # 分布式缓存
│   ├── gftoken/               # Token 认证
│   ├── mqttclient/            # MQTT 客户端
│   ├── oauth/                 # OAuth 认证
│   ├── plugins/               # 插件系统
│   ├── proxy/                 # 代理模块
│   ├── tsd/                   # 时序数据库
│   ├── utility/               # 工具类
│   ├── worker/                # 工作池
│   └── ...
├── module/                    # 模块扩展
│   ├── module.go              # 模块入口
│   └── hello/                 # 示例模块
├── manifest/                  # 配置文件
│   ├── config/                # 配置文件
│   ├── deploy/                # 部署配置
│   ├── docker-compose/        # Docker 编排
│   ├── i18n/                  # 国际化
│   ├── protobuf/             # Protobuf 定义
│   └── sql/                   # 数据库脚本
├── ui/                        # 前端界面 (Vue 3)
├── resource/                  # 资源文件
│   ├── rsa/                   # RSA 密钥
│   └── template/              # 模板文件
├── tools/                     # 开发工具
└── go.mod                     # Go 依赖管理
```

---

## 功能模块

### 设备与物联网核心

- 物模型管理 - 定义设备属性、事件、服务
- 产品管理 - 统一管理设备类型产品
- 设备管理 - 完整的设备生命周期管理
- 设备树 - 树形结构展示设备关系
- 设备标签 - 灵活的标签系统
- 实时数据 - 设备实时状态展示

### 协议与接入管理

- 多协议适配 - TCP、UDP、HTTP、WebSocket、MQTT、CoAP
- 协议网关 - 协议转换和网关管理
- 网络隧道 - 内网穿透和安全加密
- 插件系统 - 热插拔插件架构
- 数据采集协议 - Modbus、IEC61850、OPC UA、IEC104

### 数据处理与分析

- 数据中心 - 自定义数据模型
- 规则引擎 - 可视化规则编辑
- 数据转换 - ETL 管道
- 实时分析 - 流数据处理
- 时序数据库 - TDengine、InfluxDB 集成
- 数据导出 - Excel、CSV、JSON

### 告警与通知

- 告警规则 - 多条件触发
- 告警级别 - 自定义告警等级
- 告警日志 - 查询告警历史
- 通知模板 - 邮件、短信、webhook
- 消息推送 - WebSocket 实时推送

### 系统管理与权限

- 用户管理 - 用户 CRUD、角色分配
- 部门管理 - 树形组织结构
- 角色管理 - RBAC 权限控制
- 菜单管理 - 动态菜单配置
- 字典管理 - 系统字典维护
- 日志管理 - 操作日志、登录日志

### 开发与扩展工具

- 代码生成 - 前后端代码自动生成
- API 文档 - Swagger 自动生成
- 文件管理 - MinIO 集成
- 模块扩展 - 插件开发 SDK
- OpenAPI - 第三方系统集成

---

## 快速开始

### 环境要求

- Go 1.23.0+
- MySQL 5.7+ / 8.0+
- Redis 6.0+
- TDengine 3.0+ (可选)
- MQTT Broker (可选)

### 编译运行

```bash
# 克隆项目
git clone https://github.com/sagoo-cloud/sagooiot.git
cd sagooiot

# 创建目录
mkdir -p resource/log

# 配置数据库
# 导入 manifest/sql/ 目录下的初始化脚本

# 配置文件
# 修改 manifest/config/config.dev.yaml

# 运行
go run main.go
```

### Docker 部署

```bash
cd manifest/docker-compose
docker-compose up -d
```

访问 http://localhost:8000

---

## 数据流

```
设备 → 协议接入 → 协议解析 → 数据校验 → 业务处理 → 数据存储
              ↓
        实时推送 → WebSocket → 前端展示
              ↓
        规则引擎 → 告警判断 → 通知发送
              ↓
        数据分析 → 可视化报表
```

---

## 系统架构

```
┌─────────────────────────────────────┐
│         前端层 (Vue 3)               │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│       应用服务层 (GoFrame 2.9)        │
│  ┌─────────────────────────────────┐ │
│  │ 设备管理 │ 数据处理 │ 告警通知   │ │
│  │ 物模型   │ 规则引擎 │ 数据分析   │ │
│  │ 权限管理 │ 任务调度 │ 其他模块   │ │
│  └─────────────────────────────────┘ │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│       协议接入层 (多协议)             │
│  TCP │ MQTT │ UDP │ CoAP │ HTTP    │
│              ↓                      │
│        插件系统 (C/C++/Python/Go)    │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│       存储与消息中间件层              │
│  MySQL │ Redis │ TDengine │ MinIO  │
│  MQTT Broker │ 其他存储组件         │
└─────────────────────────────────────┘
```

---

## 社区与支持

- **官网文档** - http://iotdoc.sagoo.cn/
- **QQ 群** - 686637608
- **GitHub Issues** - 报告 bug 和功能建议

---

## 许可证

GPL-3.0

---

**如果这个项目对你有帮助，请点击右上角的 ⭐ Star 支持我们！**
