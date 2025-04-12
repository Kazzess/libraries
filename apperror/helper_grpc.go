package apperror

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func ErrorInfoFromDetails(sErr *status.Status, reasonHandlers map[string]func() error) error {
	for _, detail := range sErr.Details() {
		switch dt := detail.(type) {
		case *errdetails.ErrorInfo:
			for reason, handler := range reasonHandlers {
				if dt.Reason == reason {
					return handler()
				}
			}

		}
	}

	return nil
}
