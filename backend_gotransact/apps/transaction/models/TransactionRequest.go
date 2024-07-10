package models

import (
	base "gotransact/apps/Base/models"

	"gorm.io/gorm"
)

type TransactionRequest struct {
	gorm.Model
	base.Base
	UserID uint `gorm:""`
	//User                   User           `gorm:"foreignKey:UserID"`
	Status                 string `gorm:"type:varchar(20);not null"`
	PaymentGatewayMethodID uint   `gorm:"not null"`
	//PaymentGateway         PaymentGateway `gorm:"foreignKey:PaymentGatewayMethodID"`
	Description        string             `gorm:"type:text"`
	Amount             string             `gorm:"type:string;not null"`
	TransactionHistory TransactionHistory `gorm:"foreignKey:TransactionID"`
	TransactionHash    string             `gorm:"type:varchar(100)"`
}

// Enum for status
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
)
