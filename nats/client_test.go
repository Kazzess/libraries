package mynats

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/Kazzess/libraries/core/rnd"
	"github.com/Kazzess/libraries/errors"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/require"
)

const natsServer = "nats://localhost:4222"

func runNATSServer() *server.Server {
	opts := &server.Options{
		Host:           "localhost",
		Port:           4222,
		NoLog:          true, // Turn off logs during testing
		NoSigs:         true,
		MaxControlLine: 2048,
		MaxPayload:     65536,
		MaxConn:        65536,
		PingInterval:   2 * time.Minute,
		MaxPingsOut:    2,
		WriteDeadline:  2 * time.Second,
		JetStream:      true,
		StoreDir:       "/tmp/nats" + rnd.RandomString(3),
	}

	s, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

	go s.Start()

	if !s.ReadyForConnections(5 * time.Second) {
		log.Fatalf("NATS server failed to start in time.")
	}

	return s
}

func TestClient_PublishSyncAndFetch(t *testing.T) {
	s := runNATSServer()
	defer s.Shutdown()

	config := NewConfig([]string{natsServer}, "test-consumer")

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	// Setup a stream first
	streamName := "testStream"
	err = client.CreateStream(streamName, WithSubjects("testA"))
	require.NoError(t, err)

	// Publish a message
	subject := streamName + ".testA"

	message1 := []byte("Hello, NATS!")
	err = client.PublishSync(context.Background(), subject, message1)
	require.NoError(t, err)

	message2 := []byte("GOODYBYE, NATS!")
	err = client.PublishSync(context.Background(), subject, message2)
	require.NoError(t, err)

	// Fetch the message
	messages, err := client.Fetch(subject, "test-consumer", 10)
	require.NoError(t, err)
	require.Equal(t, 2, len(messages))
	require.Equal(t, message1, messages[0].Data)
	require.Equal(t, message2, messages[1].Data)
}

func TestClient_SubscribeSync(t *testing.T) {
	s := runNATSServer()
	defer s.Shutdown()

	config := NewConfig([]string{natsServer}, "test-consumer")

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	// Setup a stream
	streamName := "testStreamSub"
	err = client.CreateStream(streamName, WithSubjects("testB"))
	require.NoError(t, err)

	// Publish a message
	subject := streamName + ".testB"
	message := []byte("Hello, NATS again!")
	err = client.PublishSync(context.Background(), subject, message)
	require.NoError(t, err)

	// Subscribe and handle the message
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.SubscribeSync(ctx, subject, "test-consumer", func(ctx context.Context, msg *nats.Msg) error {
		require.Equal(t, message, msg.Data)
		return nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		require.NoError(t, err)
	}
}

func TestClient_PublishAsync(t *testing.T) {
	s := runNATSServer()
	defer s.Shutdown()

	config := NewConfig([]string{natsServer}, "test-consumer")

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	// Setup a stream first
	streamName := "testStreamAsyncPub"
	err = client.CreateStream(streamName, WithSubjects("testC"))
	require.NoError(t, err)

	// Publish a message asynchronously
	subject := streamName + ".testC"
	message := []byte("Hello, NATS async!")

	err = client.PublishAsync(context.Background(), subject, message)
	require.NoError(t, err)
}

func TestClient_SubscribeAsync(t *testing.T) {
	s := runNATSServer()
	defer s.Shutdown()

	config := NewConfig([]string{natsServer}, "test-consumer")

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	// Setup a stream
	streamName := "testStreamAsyncSub"
	err = client.CreateStream(streamName, WithSubjects("testD"))
	require.NoError(t, err)

	// Publish a message
	subject := streamName + ".testD"
	message := []byte("Hello, NATS async again!")
	err = client.PublishSync(context.Background(), subject, message)
	require.NoError(t, err)

	messageReceived := make(chan bool)

	// Subscribe and handle the message asynchronously
	err = client.SubscribeAsync(
		context.Background(),
		subject,
		"test-consumer",
		func(ctx context.Context, msg *nats.Msg) error {
			require.Equal(t, message, msg.Data)
			close(messageReceived)
			return nil
		},
	)
	require.NoError(t, err)

	select {
	case <-messageReceived:
	case <-time.After(2 * time.Second):
		t.Fatalf("Timed out waiting for async message receipt")
	}
}
