package models

import (
	base "gotransact/apps/Base/models"

	"gorm.io/gorm"
)

type TransactionHistory struct {
	gorm.Model
	base.Base
	TransactionID uint `gorm:"not null"`
	//Transaction   TransactionRequest `gorm:"foreignKey:TransactionID"`
	Status          string `gorm:"type:varchar(20);not null"`
	Description     string `gorm:"type:text"`
	Amount          string `gorm:"type:string;not null"`
	TransactionHash string `gorm:"type:varchar(100)"`
}

// // Enum for status
// const (
// 	StatusPending    = "pending"
// 	StatusProcessing = "processing"
// 	StatusSuccess    = "success"
// 	StatusFailed     = "failed"
// )
