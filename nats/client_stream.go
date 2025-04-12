package mynats

import (
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"git.adapticode.com/libraries/golang/errors"
)

type StreamOption func(config *nats.StreamConfig)

func WithName(name string) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Name = name
	}
}

func WithDescription(desc string) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Description = desc
	}
}

func WithSubjects(subjects ...string) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Subjects = subjects
	}
}

func WithRetention(retention nats.RetentionPolicy) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Retention = retention
	}
}

func WithMaxConsumers(maxConsumers int) StreamOption {
	return func(c *nats.StreamConfig) {
		c.MaxConsumers = maxConsumers
	}
}

func WithMaxMsgs(maxMsgs int64) StreamOption {
	return func(c *nats.StreamConfig) {
		c.MaxMsgs = maxMsgs
	}
}

func WithMaxBytes(maxBytes int64) StreamOption {
	return func(c *nats.StreamConfig) {
		c.MaxBytes = maxBytes
	}
}

func WithDiscard(discard nats.DiscardPolicy) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Discard = discard
	}
}

func WithDiscardNewPerSubject(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.DiscardNewPerSubject = flag
	}
}

func WithMaxAge(maxAge time.Duration) StreamOption {
	return func(c *nats.StreamConfig) {
		c.MaxAge = maxAge
	}
}

func WithMaxMsgsPerSubject(maxMsgsPerSubject int64) StreamOption {
	return func(c *nats.StreamConfig) {
		c.MaxMsgsPerSubject = maxMsgsPerSubject
	}
}

func WithMaxMsgSize(maxMsgSize int32) StreamOption {
	return func(c *nats.StreamConfig) {
		c.MaxMsgSize = maxMsgSize
	}
}

func WithStorage(storage nats.StorageType) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Storage = storage
	}
}

func WithReplicas(replicas int) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Replicas = replicas
	}
}

func WithNoAck(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.NoAck = flag
	}
}

func WithTemplate(template string) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Template = template
	}
}

func WithDuplicates(duplicates time.Duration) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Duplicates = duplicates
	}
}

func WithPlacement(placement *nats.Placement) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Placement = placement
	}
}

func WithMirror(mirror *nats.StreamSource) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Mirror = mirror
	}
}

func WithSources(sources ...*nats.StreamSource) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Sources = sources
	}
}

func WithSealed(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.Sealed = flag
	}
}

func WithDenyDelete(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.DenyDelete = flag
	}
}

func WithDenyPurge(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.DenyPurge = flag
	}
}

func WithAllowRollup(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.AllowRollup = flag
	}
}

func WithRePublish(republish *nats.RePublish) StreamOption {
	return func(c *nats.StreamConfig) {
		c.RePublish = republish
	}
}

func WithAllowDirect(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.AllowDirect = flag
	}
}

func WithMirrorDirect(flag bool) StreamOption {
	return func(c *nats.StreamConfig) {
		c.MirrorDirect = flag
	}
}

func (c *Client) CreateStream(name string, options ...StreamOption) error {
	config := &nats.StreamConfig{
		Name: name,
	}

	for _, option := range options {
		option(config)
	}

	if config.Name == "" {
		return errors.New("stream name cannot be empty")
	}

	if len(config.Subjects) == 0 {
		return errors.New("at least one stream subjects must be specified")
	}

	var targetSubjects []string
	for _, subject := range config.Subjects {
		newSubj := subject
		if !strings.HasPrefix(subject, name) {
			newSubj = name + "." + subject
		}
		targetSubjects = append(targetSubjects, newSubj)
	}

	config.Subjects = targetSubjects

	_, err := c.js.AddStream(config)
	if err != nil {
		if errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
			return nil
		}

		return errors.Wrap(err, "js.AddStream")
	}

	return nil
}
