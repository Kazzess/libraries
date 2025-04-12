package mynats

import (
	"context"
	"fmt"
	"log/slog"

	"git.adapticode.com/libraries/golang/logging"
	"github.com/nats-io/nats.go"

	"git.adapticode.com/libraries/golang/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func (c *Client) PublishAsync(ctx context.Context, subject string, data []byte, opts ...PublishOption) (err error) {
	if c.Config.tracing {
		var span trace.Span
		ctx, span = tracing.Continue(ctx, "PublishAsync")
		defer span.End()

		tracing.TraceValue(ctx, "subject", subject)
		tracing.TraceAny(ctx, "data", data)
	}

	msg := nats.NewMsg(subject)
	msg.Data = data

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

	ack, err := c.js.PublishMsgAsync(msg, opts...)
	if err != nil {
		return fmt.Errorf("PublishAsync: %w", err)
	}

	select {
	case <-ack.Ok():
		return nil
	case ackErr := <-ack.Err():
		return fmt.Errorf("PublishAsync: %w", ackErr)
	}
}

func (c *Client) SubscribeAsync(
	ctx context.Context,
	subject, consumerID string,
	handler SubscribeHandler,
	opts ...SubscribeOption,
) (err error) {
	if consumerID == "" {
		consumerID = c.Config.consumerID
	}

	_, err = c.js.QueueSubscribe(subject, consumerID, func(msg *nats.Msg) {
		var hErr error

		if msg.Header == nil {
			msg.Header = make(nats.Header)
		}

		traceCtx := otel.GetTextMapPropagator().
			Extract(context.Background(), propagation.HeaderCarrier(msg.Header))

		mdTimestamp, mdConsumer := metricsMetadata(msg)

		ObserveDeliveryTimeMs(msg.Subject, mdConsumer, mdTimestamp, true)
		observer := ObserveProcessingTimeMs(msg.Subject, consumerID, true)

		if hErr = msg.InProgress(); hErr != nil {
			slog.With("error", hErr.Error()).Error("set message in progress")

			return
		}

		hErr = handler(traceCtx, msg)
		if hErr != nil {
			logging.WithAttrs(
				ctx,
				logging.ErrAttr(hErr),
				logging.StringAttr("data", string(msg.Data)),
				logging.StringAttr("subject", msg.Subject),
				logging.AnyAttr("header", msg.Header),
			).Error("subscribe async handle error")
		}

		observer(&hErr)
	}, opts...)
	if err != nil {
		return fmt.Errorf("SubscribeAsync: %w", err)
	}

	return nil
}
