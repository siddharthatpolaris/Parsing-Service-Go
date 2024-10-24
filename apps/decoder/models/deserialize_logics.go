package models

import (
	"github.com/jackc/pgtype"
)

type DeserializeLogics struct {
	ID          int32 `gorm:"primaryKey"`
	CmdID       int32
	ArgumentKey pgtype.JSONB
	SourcePort  int32
	IndexNo     int32
	Length      int32
}

func (DeserializeLogics) TableName() string {
	return "deserialize_logics"
}

// fw-deserializelogic mapping table
type DeserializeLogicSwVersion struct {
	ID                 int32             `gorm:"primaryKey"`
	DeserializeLogicsID int32             
	SwVersionID        int32            
	DeserializeLogics  DeserializeLogics `gorm:"foreignKey:DeserializeLogicsID;constraint:OnDelete:CASCADE"`
	SwVersion          SwVersion         `gorm:"foreignKey:SwVersionID;constraint:OnDelete:CASCADE"`
}

func (DeserializeLogicSwVersion) TableName() string {
	return "deserialize_logic_sw_version"
}
