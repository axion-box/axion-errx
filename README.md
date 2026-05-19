# axion-errx

`axion-errx` 是 `axion` 系列项目通用的扩展 Error 定义仓库。

它的主要作用是提供一套共享的 `errx` 基础类型，方便 `axion-agent`、`axion-bridge` 以及后续组件统一表达以下信息：

- 错误类型 `Type`
- 稳定错误码 `Code`
- 用户可读错误消息 `Message`
- 创建时堆栈 `Stacktrace`
- 底层原因链 `Cause`

当前对外共享包暴露在：

- `github.com/axion-box/axion-errx/pkgs/errx`

## 设计目标

- 统一 `axion` 系列项目的通用错误模型
- 为组件内部和 HTTP 服务入口提供稳定的错误分类与结构化字段
- 收敛跨仓库复用的通用错误码，避免重复定义和漂移

## 错误码约定

- `1000~1999`: `axion-agent` / `axion-bridge` / 共享基础设施可复用的通用错误段

当前已合并的基础类型包括：

- `internal_error` `1000`
- `illegal_argument` `1001`
- `unauthorized` `1002`
- `forbidden` `1003`
- `not_found` `1004`
- `conflict` `1005`
- `data_unavailable` `1006`
- `rejected_operation` `1007`
- `unsupported_operation` `1008`
- `illegal_state` `1009`
- `illegal_format` `1010`
- `external_error` `1011`
- `initialization_failed` `1012`
- `method_not_allowed` `1013`
- `bad_request_body` `1014`

## 示例

```go
package main

import (
	"errors"
	"fmt"

	"github.com/axion-box/axion-errx/pkgs/errx"
)

func main() {
	base := errors.New("dial tcp timeout")
	err := errx.ExternalError.Wrap(base, "cloud request failed")

	fmt.Println(err.Code())
	fmt.Println(err.Type().Name())
	fmt.Println(err.Message())
	fmt.Println(err.Stacktrace() != "")
}
```

## 开发

- 单元测试：`bash scripts/run_unit_tests.sh`
- 门禁：`bash scripts/run_gates.sh`
