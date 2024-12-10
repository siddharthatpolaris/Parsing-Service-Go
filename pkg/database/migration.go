package database

import(
	decoder "parsing-service/apps/decoder/models"
)

var migrationModels = []interface{}{
	&decoder.CommandMapping{},
	// &decoder.SwVersion{},
	// &decoder.CommandMappingSwVersion{},
	&decoder.DeserializeLogics{},
	// &decoder.DeserializeLogicSwVersion{},

}