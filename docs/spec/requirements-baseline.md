# 需求基线

`axion-errx` 是 `axion` 系列项目通用的扩展 Error 定义仓库。

核心职责：

- 提供一个可派生的 Error 基础类型
- 支持错误 `Code`、`Type`、`Stacktrace`、`Error Message`
- 为组件内部与 HTTP 服务统一错误表达
- 收敛 `axion-agent` 与 `axion-bridge` 的通用错误码

范围边界：

- 本仓库不负责 HTTP handler、middleware、日志框架或配置系统
- 本仓库不负责业务私有错误码段
- 本仓库只维护共享错误抽象与共享通用错误类型
