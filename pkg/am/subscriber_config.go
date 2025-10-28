package am

import "time"

type AckType int

const (
	AckTypeAuto AckType = iota
	AckTypeManual
)

type DeliveryPolicyType int

const (
	AllDeliveryPolicy DeliveryPolicyType = iota
	NewDeliveryPolicy
)

var defaultAckWait = 30 * time.Second
var defaultMaxDeliver = 5

type SubscriberConfig struct {
	msgFilter     []string
	groupName     string
	ackType       AckType
	ackWait       time.Duration
	maxRedeliver  int
	deliverPolicy DeliveryPolicyType
}

func NewSubscriberConfig(options []SubscriberOption) SubscriberConfig {
	cfg := SubscriberConfig{
		msgFilter:     make([]string, 0),
		groupName:     "",
		ackType:       AckTypeManual,
		ackWait:       defaultAckWait,
		maxRedeliver:  defaultMaxDeliver,
		deliverPolicy: NewDeliveryPolicy,
	}

	for _, option := range options {
		option.configureSubscriberConfig(&cfg)
	}

	return cfg
}

type SubscriberOption interface {
	configureSubscriberConfig(*SubscriberConfig)
}

func (s SubscriberConfig) MessageFilters() []string          { return s.msgFilter }
func (s SubscriberConfig) GroupName() string                 { return s.groupName }
func (s SubscriberConfig) AckType() AckType                  { return s.ackType }
func (s SubscriberConfig) AckWait() time.Duration            { return s.ackWait }
func (s SubscriberConfig) MaxRedeliver() int                 { return s.maxRedeliver }
func (s SubscriberConfig) DeliverPolicy() DeliveryPolicyType { return s.deliverPolicy }

type MessageFilter []string

func (s MessageFilter) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.msgFilter = s }

type GroupName string

func (g GroupName) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.groupName = string(g) }

func (a AckType) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.ackType = a }

type AckWait time.Duration

func (w AckWait) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.ackWait = time.Duration(w) }

type MaxDeliver int

func (d MaxDeliver) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.maxRedeliver = int(d) }

func (d DeliveryPolicyType) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.deliverPolicy = d }
