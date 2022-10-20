package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Fields map[string]interface{}

type Logger interface {
	Debug(ctx context.Context, msg string, fields Fields)
	Error(ctx context.Context, msg string, err error)
	Log(ctx context.Context, msg string, fields Fields)
}

type JSONLogger struct {
	EnableDebug bool
	Output      io.Writer
}

func (jl *JSONLogger) Debug(ctx context.Context, msg string, fields Fields) {
	if !jl.EnableDebug {
		return
	}
	jl.init()
	jl.log(ctx, msg, fields)
}

func (jl *JSONLogger) Error(ctx context.Context, msg string, err error) {
	jl.init()
	jl.log(ctx, msg, Fields{"Error": err.Error()})
}

func (jl *JSONLogger) Log(ctx context.Context, msg string, fields Fields) {
	jl.init()
	jl.log(ctx, msg, fields)
}

func (jl *JSONLogger) init() {
	if jl.Output == nil {
		jl.Output = os.Stdout
	}
}

func (jl *JSONLogger) log(ctx context.Context, msg string, fieldData interface{}) {
	body := map[string]interface{}{
		"Message": msg,
		"Fields":  fieldData,
	}
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: json marshaling failed: %v: %#v\n", err, body)
	}
	fmt.Fprintln(jl.Output, string(b))
}
