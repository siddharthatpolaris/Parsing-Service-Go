package models

type CommandMapping struct {
	ID          int32  `gorm:"primaryKey"`
	CmdName     string `gorm:"type:varchar(50)"`
	CmdID       int32
	TopicID     int32
	SP          int32
	DP          int32
	ActualCmdID int32
	GroupID     int32
}

func (CommandMapping) TableName() string {
	return "command_mapping"
}

type SwVersionGroup struct {
	ID        int32  `gorm:"primaryKey"`
	Version   string `gorm:"type:varchar(30)"`
	GroupIDCM int32
	GroupIDDL int32
}

func (SwVersionGroup) TableName() string {
	return "sw_version_group"
}

// // fwversion cmd mapping
// type CommandMappingSwVersion struct {
// 	ID               int32          `gorm:"primaryKey"`
// 	CommandMappingID int32
// 	SwVersionID      int32
// 	CommandMapping   CommandMapping `gorm:"foreignKey:CommandMappingID;constraint:OnDelete:CASCADE"`
// 	SwVersion        SwVersion      `gorm:"foreignKey:SwVersionID;constraint:OnDelete:CASCADE"`
// }

// func (CommandMappingSwVersion) TableName() string {
// 	return "command_mapping_sw_version"
// }

// type SwVersionGroup struct {
// 	ID int32 `gorm:"primaryKey"`
// 	SwVersionID string	`gorm:"type:varchar(30)"`
// 	GroupID int32
// }

// func (SwVersionGroup) TableName() string {
// 	return "sw_version_group"
// }
