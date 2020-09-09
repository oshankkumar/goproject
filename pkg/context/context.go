package context

import (
	"context"

	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
)

type key int

const (
	logKey key = iota + 1
	reqIDKey
)

func WithLogger(parent context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(parent, logKey, logger)
}

func WithReqID(ctx context.Context, reqid uuid.UUID) context.Context {
	return context.WithValue(ctx, reqIDKey, reqid)
}

func Logger(ctx context.Context) logrus.FieldLogger {
	reqLogger, ok := ctx.Value(logKey).(logrus.FieldLogger)
	if !ok {
		return logrus.StandardLogger()
	}
	return reqLogger
}
