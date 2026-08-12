# errx

`errx` 是 `axion` 系列项目通用的结构化错误包，用来统一表达稳定错误码、错误类型、错误消息、原因链和创建时堆栈。

对外共享包路径：

- `github.com/axion-box/axion-errx/pkgs/errx`

## 基本用法

### 创建一个新的 err

```go
package main

import (
	"fmt"

	"github.com/axion-box/axion-errx/pkgs/errx"
)

func main() {
	err := errx.IllegalArgument.New("invalid user id: %d", 42)

	fmt.Println(err.Error())       // illegal_argument: invalid user id: 42
	fmt.Println(err.Type().Name()) // illegal_argument
	fmt.Println(err.Code())        // 1001
	fmt.Println(err.Message())     // invalid user id: 42
	fmt.Println(err.Stacktrace() != "") // true
}
```

`Type.New(...)` 会创建一个 `*errx.Error`，并在创建时自动捕获堆栈。

### Wrap 一个旧 err

日常场景优先用 `errx.WrapError(...)`。它会自动区分普通 `error` 和已有的 `errx.Error`：

- 普通 `error`：按 `external_error` 包装，并生成 `new: old.Error()`
- 已有 `errx.Error`：沿用原来的 `Type` 再包装一层，并生成 `new: old.Message()`

```go
package main

import (
	"errors"
	"fmt"

	"github.com/axion-box/axion-errx/pkgs/errx"
)

func main() {
	base := errors.New("dial tcp timeout")
	err := errx.WrapError(base, "cloud request failed")

	fmt.Println(err.Error())          // external_error: cloud request failed: dial tcp timeout
	fmt.Println(err.Type().Name())    // external_error
	fmt.Println(err.Message())        // cloud request failed: dial tcp timeout
	fmt.Println(err.Cause() == base)  // true
	fmt.Println(err.Stacktrace() != "") // true
}
```

如果传入的已经是 `errx.Error`，会复用原类型并拼接消息：

```go
base := errx.NotFound.New("profile missing")
err := errx.WrapError(base, "while refreshing cache")

fmt.Println(err.Type().Name()) // not_found
fmt.Println(err.Message())     // while refreshing cache: profile missing
```

`Type.Wrap(...)` 仍然会保留原始错误作为 `cause`，同时为外层错误补充新的类型、消息和堆栈。它更适合“调用方明确知道要强制指定哪一种错误类型”的场景。

如果包装时不给新消息，会回退到原始错误消息：

```go
base := errors.New("parse body failed")
err := errx.WrapError(base, "")

fmt.Println(err.Message()) // parse body failed
```

### 从 error 链里读取 errx 信息

```go
package main

import (
	"errors"
	"fmt"

	"github.com/axion-box/axion-errx/pkgs/errx"
)

func main() {
	base := errors.New("record missing")
	err := errx.WrapError(errx.NotFound.Wrap(base, "load profile failed"), "while reading profile view")

	fmt.Println(errx.Code(err))                     // 1004
	fmt.Println(errx.IsOfType(err, errx.NotFound)) // true

	if ex := errx.Cast(err); ex != nil {
		fmt.Println(ex.Type().Name())     // not_found
		fmt.Println(ex.Message())         // while reading profile view: load profile failed: record missing
		fmt.Println(ex.Stacktrace() != "") // true
	}
}
```

常用 helper：

- `errx.WrapError(err, ...)`：日常包装错误的推荐入口
- `errx.Cast(err)`：从错误链里取出第一条 `*errx.Error`
- `errx.Code(err)`：读取稳定错误码；普通 `error` 会回退为 `internal_error`
- `errx.IsOfType(err, t)`：判断错误链里是否包含指定类型
- `errx.Attrs(err)`：导出结构化字段，便于日志或协议层继续封装

补充说明：

- `Message()` 返回不带类型前缀的消息体
- `Error()` 返回适合直接打印的 `type: message`

## 预定义错误

`errx` 内置了一组共享错误类型，供多个仓库复用。当前共享错误码段为 `1000~1999`，这里已经占用：

| Type | Code |
| --- | ---: |
| `internal_error` | `1000` |
| `illegal_argument` | `1001` |
| `unauthorized` | `1002` |
| `forbidden` | `1003` |
| `not_found` | `1004` |
| `conflict` | `1005` |
| `data_unavailable` | `1006` |
| `rejected_operation` | `1007` |
| `unsupported_operation` | `1008` |
| `illegal_state` | `1009` |
| `illegal_format` | `1010` |
| `external_error` | `1011` |
| `initialization_failed` | `1012` |
| `method_not_allowed` | `1013` |
| `bad_request_body` | `1014` |
| `upstream_disconnected` | `1015` |

其中：

- `errx.NotImplemented` 是 `errx.UnsupportedOperation` 的别名
- `errx.UpstreamDisconnected` 只表示外部系统未发送可解释错误终态便断开；上游已经发送的错误应保留其原始分类
- 这些类型适合放在共享层复用，不适合承载某个业务仓库私有语义

## 定义业务错误

在业务仓库里，可以基于 `errx.NewType(name, code)` 定义自己的稳定错误类型：

```go
package billing

import "github.com/axion-box/axion-errx/pkgs/errx"

var InvoiceClosed = errx.NewType("invoice_closed", 21001)
```

然后像预定义类型一样继续使用：

```go
func CloseInvoice() error {
	return InvoiceClosed.New("invoice %s has already been closed", "INV-001")
}
```

约定：

- `NewType` 统一使用 `errx.NewType("name", code)` 形式
- `name` 和 `code` 在同一进程内必须唯一，重复注册会 `panic`
- `1000~1999` 是共享层保留错误码段；业务私有错误码应在各自仓库维护

## 开发

- 单元测试：`bash scripts/run_unit_tests.sh`
- 门禁：`bash scripts/run_gates.sh`
