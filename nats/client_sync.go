package mynats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/Kazzess/libraries/errors"
	"github.com/Kazzess/libraries/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func (c *Client) PublishSync(ctx context.Context, subject string, data []byte, opts ...PublishOption) (err error) {
	if c.Config.tracing {
		var span trace.Span
		ctx, span = tracing.Start(ctx, "PublishSync")
		defer span.End()

		tracing.TraceValue(ctx, "subject", subject)
		tracing.TraceAny(ctx, "data", data)
	}

	msg := nats.NewMsg(subject)
	msg.Data = data

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

	_, err = c.js.PublishMsg(msg, opts...)
	if err != nil {
		return fmt.Errorf("PublishSync: %w", err)
	}

	return nil
}

func (c *Client) SubscribeSync(
	ctx context.Context,
	subject, consumerID string,
	handler SubscribeHandler,
	opts ...SubscribeOption,
) error {
	if consumerID == "" {
		consumerID = c.Config.consumerID
	}

	sub, err := c.js.QueueSubscribeSync(subject, consumerID, opts...)
	if err != nil {
		return fmt.Errorf("SubscribeSync: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		msg, nextMsgErr := sub.NextMsgWithContext(ctx)
		if nextMsgErr != nil {
			return errors.Wrap(nextMsgErr, "NextMsgWithContext")
		}

		var span trace.Span
		ctx, span = tracing.Start(ctx, "SubscribeSync")
		tracing.TraceValue(ctx, "subject", subject)

		if msg.Header == nil {
			msg.Header = make(nats.Header)
		}

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Header))

		mdTimestamp, mdConsumer := metricsMetadata(msg)

		ObserveDeliveryTimeMs(msg.Subject, mdConsumer, mdTimestamp, false)
		observer := ObserveProcessingTimeMs(msg.Subject, consumerID, false)

		err = handler(ctx, msg)
		if err != nil {
			span.SetAttributes(MessageAttributes(msg)...)
			span.End()

			return errors.Wrap(err, "process handler")
		}

		observer(&err)

		span.End()
	}
}
