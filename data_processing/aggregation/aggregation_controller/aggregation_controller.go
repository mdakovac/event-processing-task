package aggregation_controller

import (
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/aggregation/aggregation_service"
	"github.com/gin-gonic/gin"
)

func SetupAggregationController(r *gin.Engine, service aggregation_service.AggregationServiceType) {
	r.GET("/materialized", func(c *gin.Context) {
		a := service.GetAggregation()
		c.JSON(200, a)
	})
}
