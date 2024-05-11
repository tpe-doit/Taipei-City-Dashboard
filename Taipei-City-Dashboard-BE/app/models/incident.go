package models

import (
	"fmt"
)

type Incident struct {
  ID          uint      `gorm:"primaryKey"`
  Type        string    `json:"inctype"`
  Description string		`json:"description"`
  Distance    float64		`json:"distance"`
  Latitude    float64		`json:"latitude"`
  Longitude   float64		`json:"longitude"`
  Time        int64			`json:"reportTime"`
}

func GetAllIncident(pageSize int, pageNum int, sort string, order string) (incidents []Incident, totalIncidents int64, resultNum int64, err error) {
	tempDB := DBManager.Table("incidents")

	// Count the total amount of incidents
	tempDB.Count(&totalIncidents)

	tempDB.Count(&resultNum)

	// Sort the issues
	if sort != "" {
		tempDB = tempDB.Order(sort + " " + order)
	}
	
	// Paginate the issues
	if pageSize > 0 {
		tempDB = tempDB.Limit(pageSize)
		if pageNum > 0 {
			tempDB = tempDB.Offset((pageNum - 1) * pageSize)
		}
	}

	// Get the incidents
	err = tempDB.Find(&incidents).Error

	return incidents, totalIncidents, resultNum, err
}

func CreateIncident(incidentType, description string, distance, latitude, longitude float64, incidentTime int64) (incident Incident, err error) {

    // Create an example incident
    // incident = Incident{
    //     Type:        "Accident",
    //     Description: "Car crash on the highway",
    //     Distance:    2.5,
    //     Latitude:    40.7128,
    //     Longitude:   -74.0060,
    //     Time:        time.Now().Unix(),
    // }
		incident = Incident {
        Type:        incidentType,
        Description: description,
        Distance:    distance,
        Latitude:    latitude,
        Longitude:   longitude,
        Time:        incidentTime,
		}

    // Insert the incident into the database
    err = DBManager.Create(&incident).Error

		return incident, err
}

func DeleteIncident(id uint) (incident Incident, err error) {
	
	if err := DBManager.Where("ID = ?", id).First(&incident, 1).Error; err != nil {
		// Handle error (e.g., incident not found)
		fmt.Printf("Incident not found")
		return Incident{}, err
	}

	// Delete the incident
	err = DBManager.Delete(&incident).Error
	return incident, err
}
