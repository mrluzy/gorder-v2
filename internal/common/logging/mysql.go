package logging

import (
	"context"
	"fmt"
	"github.com/mrluzy/gorder-v2/common/util"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	Method   = "method"
	Args     = "args"
	Cost     = "cost_ms"
	Response = "response"
	Error    = "err"
)

type ArgFormatter interface {
	FormatArg() (string, error)
}

func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
	fields := logrus.Fields{
		Method: method,
		Args:   formatMySQLArgs(args),
	}
	start := time.Now()
	return fields, func(resp any, err *error) {
		level, msg := logrus.InfoLevel, "mysql_success"
		fields[Cost] = time.Since(start).Milliseconds()
		fields[Response] = resp

		if err != nil && (*err != nil) {
			level, msg = logrus.ErrorLevel, "mysql_error"
			fields[Error] = (*err).Error()
		}
		logrus.WithContext(ctx).WithFields(fields).Logf(level, "%s", msg)
	}

}

func formatMySQLArgs(args []any) string {
	var item []string
	for _, arg := range args {
		item = append(item, formatMySQLArg(arg))
	}
	return strings.Join(item, "||")
}

func formatMySQLArg(arg any) string {
	var (
		str string
		err error
	)

	defer func() {
		if err != nil {
			str = fmt.Sprintf("unsupported type in formatMySQLArg||err=%s", err.Error())
		}
	}()

	switch v := arg.(type) {
	case ArgFormatter:
		str, err = v.FormatArg()
	default:
		str, err = util.MarshalString(v)
	}

	return str
}
