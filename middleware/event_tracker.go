package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/jwalitptl/pkg/event"
)

type EventTrackerMiddleware struct {
	config    *event.EventTrackingConfig
	extractor event.FieldExtractor
	emitter   eventservice.Service
}

func NewEventTrackerMiddleware(
	config *event.EventTrackingConfig,
	emitter eventservice.Service,
) *EventTrackerMiddleware {
	return &EventTrackerMiddleware{
		config:    config,
		extractor: &event.DefaultFieldExtractor{},
		emitter:   emitter,
	}
}

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

		var endpointConfig event.EndpointConfig
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

		// Store the event context
		c.Set("eventCtx", &EventContext{
			Resource:  resource,
			Operation: operation,
		})

		c.Next()

		// After handler execution, get the event context
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

		// Add additional data to payload
		for k, v := range eventCtx.Additional {
			payload[k] = v
		}

		m.emitter.Emit(eventservice.EventType(endpointConfig.EventType), payload)
	}
}
