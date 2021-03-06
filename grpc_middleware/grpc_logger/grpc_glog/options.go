// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_glog

import (
	"context"
	"github.com/donetkit/contrib-log/glog"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

type options struct {
	logger glog.ILoggerEntry
}

func evaluateServerOpt(opts []Option) *options {
	optCopy := &options{}
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

func evaluateClientOpt(opts []Option) *options {
	optCopy := &options{}
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

// CodeToLevel function defines the mapping between gRPC return codes and interceptor log level.
type CodeToLevel func(code codes.Code) logrus.Level

// DurationToField function defines how to produce duration fields for logging
type DurationToField func(duration time.Duration) (key string, value interface{})

// WithLogger customizes the function for mapping request durations to log fields.
func WithLogger(logger glog.ILogger) Option {
	return func(o *options) {
		if logger != nil {
			o.logger = logger.WithField("Grpc", "Grpc")
		}
	}
}

// DefaultCodeToLevel is the default implementation of gRPC return codes to log levels for server side.
func DefaultCodeToLevel(code codes.Code) logrus.Level {
	switch code {
	case codes.OK:
		return logrus.InfoLevel
	case codes.Canceled:
		return logrus.InfoLevel
	case codes.Unknown:
		return logrus.ErrorLevel
	case codes.InvalidArgument:
		return logrus.InfoLevel
	case codes.DeadlineExceeded:
		return logrus.WarnLevel
	case codes.NotFound:
		return logrus.InfoLevel
	case codes.AlreadyExists:
		return logrus.InfoLevel
	case codes.PermissionDenied:
		return logrus.WarnLevel
	case codes.Unauthenticated:
		return logrus.InfoLevel // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return logrus.WarnLevel
	case codes.FailedPrecondition:
		return logrus.WarnLevel
	case codes.Aborted:
		return logrus.WarnLevel
	case codes.OutOfRange:
		return logrus.WarnLevel
	case codes.Unimplemented:
		return logrus.ErrorLevel
	case codes.Internal:
		return logrus.ErrorLevel
	case codes.Unavailable:
		return logrus.WarnLevel
	case codes.DataLoss:
		return logrus.ErrorLevel
	default:
		return logrus.ErrorLevel
	}
}

// DefaultClientCodeToLevel is the default implementation of gRPC return codes to log levels for client side.
func DefaultClientCodeToLevel(code codes.Code) logrus.Level {
	switch code {
	case codes.OK:
		return logrus.DebugLevel
	case codes.Canceled:
		return logrus.DebugLevel
	case codes.Unknown:
		return logrus.InfoLevel
	case codes.InvalidArgument:
		return logrus.DebugLevel
	case codes.DeadlineExceeded:
		return logrus.InfoLevel
	case codes.NotFound:
		return logrus.DebugLevel
	case codes.AlreadyExists:
		return logrus.DebugLevel
	case codes.PermissionDenied:
		return logrus.InfoLevel
	case codes.Unauthenticated:
		return logrus.InfoLevel // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return logrus.DebugLevel
	case codes.FailedPrecondition:
		return logrus.DebugLevel
	case codes.Aborted:
		return logrus.DebugLevel
	case codes.OutOfRange:
		return logrus.DebugLevel
	case codes.Unimplemented:
		return logrus.WarnLevel
	case codes.Internal:
		return logrus.WarnLevel
	case codes.Unavailable:
		return logrus.WarnLevel
	case codes.DataLoss:
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}

// DefaultDurationToField is the default implementation of converting request duration to a log field (key and value).
var DefaultDurationToField = DurationToTimeMillisField

// DurationToTimeMillisField converts the duration to milliseconds and uses the key `grpc.time_ms`.
func DurationToTimeMillisField(duration time.Duration) (key string, value interface{}) {
	return "grpc.time_ms", durationToMilliseconds(duration)
}

// DurationToDurationField uses the duration value to log the request duration.
func DurationToDurationField(duration time.Duration) (key string, value interface{}) {
	return "grpc.duration", duration
}

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}

// MessageProducer produces a user defined log message
type MessageProducer func(ctx context.Context, format string, code codes.Code, err error)
