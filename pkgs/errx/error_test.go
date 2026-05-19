package errx_test

import (
	"errors"
	"testing"

	"github.com/axion-box/axion-errx/pkgs/errx"
)

func TestNewErrorCarriesTypeMessageCodeAndStacktrace(t *testing.T) {
	t.Parallel()

	err := errx.IllegalArgument.New("invalid user id: %d", 42)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid user id: 42" {
		t.Fatalf("unexpected error string: %q", err.Error())
	}
	if err.Type() != errx.IllegalArgument {
		t.Fatalf("unexpected type: got %v", err.Type())
	}
	if err.Code() != 1001 {
		t.Fatalf("unexpected code: got %d", err.Code())
	}
	if err.Message() != "invalid user id: 42" {
		t.Fatalf("unexpected message: %q", err.Message())
	}
	if err.Cause() != nil {
		t.Fatalf("unexpected cause: %v", err.Cause())
	}
	if err.Unwrap() != nil {
		t.Fatalf("unexpected unwrap value: %v", err.Unwrap())
	}
	if err.Stacktrace() == "" {
		t.Fatal("expected stacktrace")
	}
}

func TestWrapPreservesCauseAndSupportsHelpers(t *testing.T) {
	t.Parallel()

	base := errors.New("dial tcp timeout")
	err := errx.ExternalError.Wrap(base, "cloud request failed")

	if !errors.Is(err, base) {
		t.Fatal("expected wrapped cause in error chain")
	}
	if err.Cause() != base {
		t.Fatal("expected cause to be preserved")
	}
	if err.Unwrap() != base {
		t.Fatal("expected unwrap to expose cause")
	}
	if errx.Cast(err) != err {
		t.Fatal("expected cast to return wrapped errx error")
	}
	if !errx.IsOfType(err, errx.ExternalError) {
		t.Fatal("expected errx.IsOfType to match")
	}
	if errx.Code(err) != errx.ExternalError.Code() {
		t.Fatalf("unexpected code from helper: got %d", errx.Code(err))
	}
}

func TestWrapFallsBackToCauseMessageWhenMessageEmpty(t *testing.T) {
	t.Parallel()

	base := errors.New("parse body failed")
	err := errx.BadRequestBody.Wrap(base, "")

	if err.Message() != "parse body failed" {
		t.Fatalf("unexpected fallback message: %q", err.Message())
	}
}

func TestWrapReturnsNilForNilCause(t *testing.T) {
	t.Parallel()

	if got := errx.ExternalError.Wrap(nil, "ignored"); got != nil {
		t.Fatalf("expected nil wrap result, got %#v", got)
	}
}

func TestAttrsForStructuredPlainAndNilErrors(t *testing.T) {
	t.Parallel()

	structured := errx.NotFound.New("record missing")
	structuredAttrs := errx.Attrs(structured)
	if len(structuredAttrs) != 5 {
		t.Fatalf("unexpected attr count: got %d", len(structuredAttrs))
	}
	if structuredAttrs[0].Key != "error_type" || structuredAttrs[0].Value != "not_found" {
		t.Fatalf("unexpected type attr: %+v", structuredAttrs[0])
	}
	if structuredAttrs[1].Key != "error_code" || structuredAttrs[1].Value != 1004 {
		t.Fatalf("unexpected code attr: %+v", structuredAttrs[1])
	}
	if structuredAttrs[2].Key != "error_message" || structuredAttrs[2].Value != "record missing" {
		t.Fatalf("unexpected message attr: %+v", structuredAttrs[2])
	}
	if structuredAttrs[3].Key != "error_cause" || structuredAttrs[3].Value != "" {
		t.Fatalf("unexpected cause attr: %+v", structuredAttrs[3])
	}
	if structuredAttrs[4].Key != "error_stacktrace" || structuredAttrs[4].Value == "" {
		t.Fatalf("unexpected stacktrace attr: %+v", structuredAttrs[4])
	}

	receiverAttrs := structured.Attrs()
	if len(receiverAttrs) != 5 {
		t.Fatalf("unexpected receiver attr count: got %d", len(receiverAttrs))
	}
	if receiverAttrs[0] != structuredAttrs[0] || receiverAttrs[1] != structuredAttrs[1] || receiverAttrs[2] != structuredAttrs[2] {
		t.Fatalf("expected receiver attrs to match helper attrs: got %+v want %+v", receiverAttrs, structuredAttrs)
	}

	plain := errors.New("plain failure")
	plainAttrs := errx.Attrs(plain)
	if plainAttrs[0].Value != "internal_error" {
		t.Fatalf("unexpected fallback type: %+v", plainAttrs[0])
	}
	if plainAttrs[1].Value != 1000 {
		t.Fatalf("unexpected fallback code: %+v", plainAttrs[1])
	}
	if plainAttrs[2].Value != "plain failure" {
		t.Fatalf("unexpected fallback message: %+v", plainAttrs[2])
	}
	if plainAttrs[3].Value != "" {
		t.Fatalf("unexpected fallback cause: %+v", plainAttrs[3])
	}
	if plainAttrs[4].Value != "" {
		t.Fatalf("unexpected fallback stacktrace: %+v", plainAttrs[4])
	}

	nilAttrs := errx.Attrs(nil)
	if nilAttrs[0].Value != "internal_error" {
		t.Fatalf("unexpected nil type: %+v", nilAttrs[0])
	}
	if nilAttrs[1].Value != 1000 {
		t.Fatalf("unexpected nil code: %+v", nilAttrs[1])
	}
	if nilAttrs[2].Value != "" || nilAttrs[3].Value != "" || nilAttrs[4].Value != "" {
		t.Fatalf("unexpected nil attrs payload: %+v", nilAttrs)
	}
}

func TestErrorHelpersHandleNilAndPlainErrors(t *testing.T) {
	t.Parallel()

	var err *errx.Error
	if err.Error() != "" {
		t.Fatalf("unexpected nil Error() result: %q", err.Error())
	}
	if err.Message() != "" {
		t.Fatalf("unexpected nil Message() result: %q", err.Message())
	}
	if err.Cause() != nil {
		t.Fatalf("unexpected nil Cause() result: %v", err.Cause())
	}
	if err.Type() != nil {
		t.Fatalf("unexpected nil Type() result: %v", err.Type())
	}
	if err.Code() != 0 {
		t.Fatalf("unexpected nil Code() result: %d", err.Code())
	}
	if err.Stacktrace() != "" {
		t.Fatalf("unexpected nil Stacktrace() result: %q", err.Stacktrace())
	}
	if err.Unwrap() != nil {
		t.Fatalf("unexpected nil Unwrap() result: %v", err.Unwrap())
	}
	if got := err.Attrs(); len(got) != 5 {
		t.Fatalf("unexpected nil receiver attrs count: %d", len(got))
	}

	var nilType *errx.Type
	errWithNilType := nilType.New("detached")
	if errWithNilType.Type() != nil {
		t.Fatalf("expected nil type, got %v", errWithNilType.Type())
	}
	if errWithNilType.Code() != 0 {
		t.Fatalf("expected zero code for nil type, got %d", errWithNilType.Code())
	}

	if errx.Cast(nil) != nil {
		t.Fatalf("expected nil cast result, got %#v", errx.Cast(nil))
	}
	if errx.Cast(errors.New("plain")) != nil {
		t.Fatalf("expected plain error cast to be nil, got %#v", errx.Cast(errors.New("plain")))
	}
	if errx.Code(nil) != 0 {
		t.Fatalf("unexpected nil helper code: %d", errx.Code(nil))
	}
	if errx.Code(errors.New("plain")) != errx.InternalError.Code() {
		t.Fatalf("unexpected plain helper code: %d", errx.Code(errors.New("plain")))
	}
	if errx.IsOfType(errors.New("plain"), errx.NotFound) {
		t.Fatal("did not expect plain error to match type")
	}
	if errx.IsOfType(errx.NotFound.New("missing"), nil) {
		t.Fatal("did not expect nil type matcher to succeed")
	}
}
