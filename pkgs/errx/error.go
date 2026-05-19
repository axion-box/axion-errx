package errx

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error 是带类型码、原因链和堆栈的结构化业务错误。
type Error struct {
	typ        *Type
	message    string
	cause      error
	stacktrace string
}

// Attr 表示一条结构化错误字段。
type Attr struct {
	// Key 表示当前结构化字段的键名。
	Key string

	// Value 表示当前结构化字段对应的值。
	Value any
}

// New 基于当前错误类型创建一条新的结构化错误。
func (t *Type) New(message string, args ...any) *Error {
	return &Error{
		typ:        t,
		message:    formatMessage(message, args...),
		stacktrace: captureStacktrace(),
	}
}

// Wrap 用当前错误类型包装一个已有错误，并补充新的用户可读消息。
func (t *Type) Wrap(err error, message string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		typ:        t,
		message:    coalesceMessage(formatMessage(message, args...), err),
		cause:      err,
		stacktrace: captureStacktrace(),
	}
}

// Error 返回错误的用户可读消息。
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Message 返回错误消息本体，不包含类型和堆栈信息。
func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Cause 返回当前结构化错误包装的根因错误。
func (e *Error) Cause() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Type 返回当前错误绑定的错误类型。
func (e *Error) Type() *Type {
	if e == nil {
		return nil
	}
	return e.typ
}

// Code 返回当前错误类型的稳定数值编码。
func (e *Error) Code() int {
	if e == nil || e.typ == nil {
		return 0
	}
	return e.typ.Code()
}

// Stacktrace 返回创建这条结构化错误时捕获的堆栈文本。
func (e *Error) Stacktrace() string {
	if e == nil {
		return ""
	}
	return e.stacktrace
}

// Attrs 返回当前错误对象对应的结构化字段集合。
func (e *Error) Attrs() []Attr {
	return buildAttrs(buildAttrsArgs{
		errorType:  e.Type().Name(),
		errorCode:  e.Code(),
		message:    e.Message(),
		cause:      e.Cause(),
		stacktrace: e.Stacktrace(),
	})
}

// Unwrap 让 Error 可以接入标准库 errors 链。
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Cast 在错误链上向内查找第一条结构化 Error。
func Cast(err error) *Error {
	for err != nil {
		var target *Error
		if errors.As(err, &target) {
			return target
		}
		err = errors.Unwrap(err)
	}
	return nil
}

// Code 返回当前错误链对应的稳定错误码。
func Code(err error) int {
	if err == nil {
		return 0
	}
	if ex := Cast(err); ex != nil {
		return ex.Code()
	}
	return InternalError.Code()
}

// IsOfType 判断错误链里是否包含指定的错误类型。
func IsOfType(err error, t *Type) bool {
	if t == nil {
		return false
	}
	cast := Cast(err)
	return cast != nil && cast.Type() == t
}

// Attrs 返回当前错误链对应的结构化字段集合。
func Attrs(err error) []Attr {
	if ex := Cast(err); ex != nil {
		return ex.Attrs()
	}

	message := ""
	if err != nil {
		message = err.Error()
	}
	return buildAttrs(buildAttrsArgs{
		errorType: InternalError.Name(),
		errorCode: InternalError.Code(),
		message:   message,
	})
}

type buildAttrsArgs struct {
	errorType  string
	errorCode  int
	message    string
	cause      error
	stacktrace string
}

func buildAttrs(args buildAttrsArgs) []Attr {
	cause := ""
	if args.cause != nil {
		cause = args.cause.Error()
	}
	return []Attr{
		{Key: "error_type", Value: args.errorType},
		{Key: "error_code", Value: args.errorCode},
		{Key: "error_message", Value: args.message},
		{Key: "error_cause", Value: cause},
		{Key: "error_stacktrace", Value: args.stacktrace},
	}
}

func formatMessage(message string, args ...any) string {
	if len(args) == 0 {
		return message
	}
	return fmt.Sprintf(message, args...)
}

func coalesceMessage(message string, cause error) string {
	if message != "" {
		return message
	}
	if cause != nil {
		return cause.Error()
	}
	return ""
}

func captureStacktrace() string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	lines := make([]string, 0, n)

	for {
		frame, more := frames.Next()
		lines = append(lines, fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return strings.Join(lines, "\n")
}
