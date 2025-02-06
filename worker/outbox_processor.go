package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"aiclinic/pkg/messaging"
	"aiclinic/pkg/model"
	"aiclinic/pkg/repository/postgres"
)

type OutboxProcessor struct {
	outboxRepo postgres.OutboxRepository
	broker     messaging.MessageBroker
}

func NewOutboxProcessor(outboxRepo postgres.OutboxRepository, broker messaging.MessageBroker) *OutboxProcessor {
	return &OutboxProcessor{
		outboxRepo: outboxRepo,
		broker:     broker,
	}
}

func (p *OutboxProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processEvents(ctx)
		}
	}
}

func (p *OutboxProcessor) processEvents(ctx context.Context) {
	events, err := p.outboxRepo.GetPendingEvents(ctx, 100)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get pending events")
		return
	}

	for _, event := range events {
		err := p.broker.Publish(ctx, event.EventType, event.Payload)
		if err != nil {
			log.Error().Err(err).Str("event_id", event.ID.String()).Msg("Failed to publish event")
			errStr := err.Error()
			_ = p.outboxRepo.UpdateStatus(ctx, event.ID, model.OutboxStatusFailed, &errStr)
			continue
		}

		err = p.outboxRepo.UpdateStatus(ctx, event.ID, model.OutboxStatusProcessed, nil)
		if err != nil {
			log.Error().Err(err).Str("event_id", event.ID.String()).Msg("Failed to update event status")
		}
	}
}
