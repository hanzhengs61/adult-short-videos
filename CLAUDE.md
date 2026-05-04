# CLAUDE.md

## 语言
所有对话、文档、代码注释统一使用中文。

## 工作方式
- 直接写代码，不走 brainstorm→plan→subagent 链路
- 不自动 git commit，由用户手动审查提交
- 优先编辑已有文件，不轻易新建
- 嵌套层数不超过 4 层
- 单文件不超过 400 行，超出须拆分
- 每次只解决一个具体问题
- 修改或新增代码时，对非显而易见的逻辑、字段定义、常量值加中文注释说明
- 设计要预设通用性，不只做最小修改（项目仍在早期）
- 输出简洁，推理详尽

## 特殊架构约定
- 冷热分离：视频分 cold（源站 URL）和 hot（本地 CDN URL），`heat_calculator.go` 每 5 分钟计算热度
- 防盗链代理：冷数据走 `/api/storage/proxy`，代理伪造请求头，M3U8 分片 URL 递归重写
- 四层结构：handler → logic → repository → model
