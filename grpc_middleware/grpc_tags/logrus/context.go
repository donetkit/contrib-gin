package ctx_logrus

import (
	"context"
	"github.com/donetkit/contrib-gin/grpc_middleware/grpc_logger/grpc_glog/ctxlogrus"

	"github.com/sirupsen/logrus"
)

// AddFields adds logrus fields to the logger.
// Deprecated: should use the ctxlogrus.Extract instead
func AddFields(ctx context.Context, fields logrus.Fields) {
	ctxlogrus.AddFields(ctx, fields)
}

// Extract takes the call-scoped logrus.Entry from grpc_logrus middleware.
// Deprecated: should use the ctxlogrus.Extract instead
func Extract(ctx context.Context) *logrus.Entry {
	return ctxlogrus.Extract(ctx)
}

// ToContext adds the logrus.Entry to the context for extraction later.
// Deprecated: should use ctxlogrus.ToContext instead
func ToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	return ctxlogrus.ToContext(ctx, entry)
}
