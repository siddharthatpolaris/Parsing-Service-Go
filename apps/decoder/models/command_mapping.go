package models

type CommandMapping struct {
	ID          int32 `gorm:"primaryKey"`
	CmdName     string `gorm:"type:varchar(50)"`
	CmdID       int32
	TopicID     int32
	SP          int32
	DP          int32
	ActualCmdID int32
}

func (CommandMapping) TableName() string {
	return "command_mapping"
}

// fwVersion
type SwVersion struct {
	ID      int32 `gorm:"primaryKey"`
	Version string	`gorm:"type:varchar(30)"`
}

func (SwVersion) TableName() string {
	return "sw_version"
}

// fwversion cmd mapping
type CommandMappingSwVersion struct {
	ID               int32          `gorm:"primaryKey"`
	CommandMappingID int32          
	SwVersionID      int32         
	CommandMapping   CommandMapping `gorm:"foreignKey:CommandMappingID;constraint:OnDelete:CASCADE"`
	SwVersion        SwVersion      `gorm:"foreignKey:SwVersionID;constraint:OnDelete:CASCADE"`
}

func (CommandMappingSwVersion) TableName() string {
	return "command_mapping_sw_version"
}
