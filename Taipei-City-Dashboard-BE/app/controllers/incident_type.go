package controllers

import (
	"TaipeiCityDashboardBE/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateIncidentType(c *gin.Context) {
	var incidentType models.IncidentType

	if err := c.ShouldBindJSON(&incidentType); (err != nil) || (incidentType.Type == "") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	tmpIncident, err := models.CreateIncidentType(incidentType.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if tmpIncident.IsEmpty() {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Incident type " + incidentType.Type + " already exists!!"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": tmpIncident})
}


func UpdateIncidentType(c *gin.Context) {
	var incident models.Incident

	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	tmpIncident, err := models.UpdateIncidentType(incident.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": tmpIncident})
}