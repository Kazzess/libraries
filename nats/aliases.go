package mynats

import "github.com/nats-io/nats.go"

type (
	Msg             = nats.Msg
	JSOpt           = nats.JSOpt
	PublishOption   = nats.PubOpt
	SubscribeOption = nats.SubOpt
	Subscription    = nats.Subscription
	AckWait         = nats.AckWait
	MaxWait         = nats.MaxWait
)

const (
	FileStorage   = nats.FileStorage
	MemoryStorage = nats.MemoryStorage

	LimitsPolicy    = nats.LimitsPolicy
	InterestPolicy  = nats.InterestPolicy
	WorkQueuePolicy = nats.WorkQueuePolicy

	DiscardOld = nats.DiscardOld
	DiscardNew = nats.DiscardNew
)

var (
	ReconnectWait           = nats.ReconnectWait
	MaxReconnects           = nats.MaxReconnects
	Timeout                 = nats.Timeout
	AckNone                 = nats.AckNone
	AckAll                  = nats.AckAll
	AckExplicit             = nats.AckExplicit
	Bind                    = nats.Bind
	UserCredentials         = nats.UserCredentials
	UserInfo                = nats.UserInfo
	Token                   = nats.Token
	MaxAckPending           = nats.MaxAckPending
	DeliverAll              = nats.DeliverAll
	DeliverLast             = nats.DeliverLast
	DeliverLastPerSubject   = nats.DeliverLastPerSubject
	DeliverNew              = nats.DeliverNew
	Durable                 = nats.Durable
	ManualAck               = nats.ManualAck
	OrderedConsumer         = nats.OrderedConsumer
	InactiveThreshold       = nats.InactiveThreshold
	ErrConnectionClosed     = nats.ErrConnectionClosed
	ErrConsumerNameRequired = nats.ErrConsumerNameRequired
	ErrStreamNameRequired   = nats.ErrStreamNameRequired
	ErrConsumerNotFound     = nats.ErrConsumerNotFound
	ErrTimeout              = nats.ErrTimeout
)
