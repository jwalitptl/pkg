package event

import "github.com/gin-gonic/gin"

type EventHandler interface {
	RegisterRoutesWithEvents(r *gin.RouterGroup, eventTracker *EventTrackerMiddleware)
}
