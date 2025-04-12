package apperror

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Kazzess/libraries/logging"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GRPCUnaryInterceptor(systemCode string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if err != nil {
			var requestAttr slog.Attr

			if stringer, ok := req.(fmt.Stringer); ok {
				requestAttr = logging.StringAttr("request", stringer.String())
			} else {
				requestAttr = logging.StringAttr("request", fmt.Sprintf("%v", req))
			}

			logging.WithAttrs(ctx, logging.ErrAttr(err), requestAttr).Error("request failed")

			if s, ok := status.FromError(err); ok {
				return nil, s.Err()
			}

			internalError := NewInternalError(
				systemCode,
				WithMessage("unknown internal system error"),
			).WithTrace(ctx)

			sentry.CaptureException(internalError)

			return nil, internalError
		}

		return resp, nil
	}
}
