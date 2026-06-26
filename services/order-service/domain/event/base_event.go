package event

import "time"

type DomainEvent interface {
	EventName() string
	OccurentAt() time.Time
	AggregatorID() string
}

type BaseEvent struct {
	eventName    string
	occurentAt   time.Time
	aggregatorID string
}

func NewBaseEvent(name, aggregatorID string) BaseEvent {
	return BaseEvent{
		eventName:    name,
		occurentAt:   time.Now().UTC(),
		aggregatorID: aggregatorID,
	}
}

func (b BaseEvent) EventName() string     { return b.eventName }
func (b BaseEvent) OccurentAt() time.Time { return b.occurentAt }
func (b BaseEvent) AggregatorID() string  { return b.aggregatorID }
