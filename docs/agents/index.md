# Agent 规则总入口

本仓库采用精简版 agent harness，只保留小型共享库真正需要的约束：

- 先对齐仓库定位，再落代码
- 共享包变更优先考虑跨仓库复用与兼容性
- 文档、脚本、测试只保留当前仓库确实在用的最小集合

执行口径：

- 勘误：见 [errata.md](errata.md)
- 小仓库保持骨架轻量，避免照搬大型仓库模板
- 代码变更应伴随单元测试；只改文档时至少保证门禁可通过
- 公共类型、错误码、导出 API 一旦对外暴露，应优先保持兼容
- 新增导出 API 时，同步更新 README 或相关规范文档
- 验证统一从 `scripts/agents/check_all_gates.py` 或 `bash scripts/run_gates.sh` 进入

常用命令：

- `bash scripts/run_unit_tests.sh`
- `bash scripts/run_gates.sh`
- `go test ./...`
