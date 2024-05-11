package models

import (
	"fmt"
)

type IncidentType struct {
	ID 					uint 			`gorm:"primaryKey"`
	Type				string		`json:"type" gorm:"not null"`
	Count				int				`json:"count"`
}

func (m IncidentType) IsEmpty() bool {
	return m.Type == "" && m.Count == 0
}

func CreateIncidentType(incidentType string) (newType IncidentType, err error){
	newType = IncidentType{
		Type:					incidentType,
		Count:				0,
	}
	if err := DBManager.Where("type = ?", incidentType).Error; err == nil {
		fmt.Printf("Incident type " + incidentType + " already exists!!")
		return IncidentType{}, nil
	}

	tempDB := DBManager.Table("incident_types")
	err = tempDB.Create(&newType).Error
	return newType, err
}

func UpdateIncidentType(incidentType string) (updType IncidentType, err error){
	if err := DBManager.Where("type = ?", incidentType).First(&updType, 1).Error; err != nil {
			// Handle error (e.g., incident not found)
			fmt.Printf("Incident type" + incidentType + " not found")
			return IncidentType{}, err
	}

	updType.Count += 1
	tempDB := DBManager.Table("incident_types")
	if err := tempDB.Save(&updType).Error; err != nil {
		// Handle error
		fmt.Printf("Failed to update incident type " + incidentType)
		return IncidentType{}, err
	}
	return updType, err
}
