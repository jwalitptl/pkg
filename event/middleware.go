package event

import (
	"github.com/gin-gonic/gin"
)

type EventTrackerMiddleware struct {
	config    *config.EventTrackingConfig
	extractor FieldExtractor
	emitter   eventservice.Service
}

func NewEventTrackerMiddleware(
	config *config.EventTrackingConfig,
	emitter eventservice.Service,
) *EventTrackerMiddleware {
	return &EventTrackerMiddleware{
		config:    config,
		extractor: &DefaultFieldExtractor{},
		emitter:   emitter,
	}
}

// Copy rest of the middleware code from pkg/middleware/event_tracker.go

type EventContext struct {
	Resource   string
	Operation  string
	OldData    interface{}
	NewData    interface{}
	Additional map[string]interface{}
}

func (m *EventTrackerMiddleware) TrackEvent(resource, operation string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			c.Next()
			return
		}

		resourceConfig, exists := m.config.Endpoints[resource]
		if !exists {
			c.Next()
			return
		}

		var endpointConfig config.EndpointConfig
		switch operation {
		case "create":
			endpointConfig = resourceConfig.Create
		case "update":
			endpointConfig = resourceConfig.Update
		case "delete":
			endpointConfig = resourceConfig.Delete
		default:
			c.Next()
			return
		}

		if !endpointConfig.Enabled {
			c.Next()
			return
		}

		c.Set("eventCtx", &EventContext{
			Resource:  resource,
			Operation: operation,
		})

		c.Next()

		eventCtxI, exists := c.Get("eventCtx")
		if !exists {
			return
		}

		eventCtx := eventCtxI.(*EventContext)
		var payload map[string]interface{}

		if endpointConfig.TrackChanges && eventCtx.OldData != nil && eventCtx.NewData != nil {
			payload = m.extractor.ExtractChanges(eventCtx.OldData, eventCtx.NewData, endpointConfig.TrackedFields)
		} else if eventCtx.NewData != nil {
			payload = m.extractor.ExtractFields(eventCtx.NewData, endpointConfig.TrackedFields)
		}

		for k, v := range eventCtx.Additional {
			payload[k] = v
		}

		m.emitter.Emit(eventservice.EventType(endpointConfig.EventType), payload)
	}
}
