package database

import (
	"parsing-service/pkg/config"
	"parsing-service/pkg/logger"
	"parsing-service/scripts"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/driver/postgres"
	"gorm.io/plugin/dbresolver"
)

var (
	DB *gorm.DB
	DBErr error
)

type Database struct {
	*gorm.DB
}


// SetupConnection creates a database connection
func SetUpDbConnection(logger logger.ILogger) error {
	masterDSN, replicaDSN := config.DbConfiguration()

	gormLogMode := viper.GetBool("MASTER_DB_LOG_MODE")
	debug := viper.GetBool("DEBUG")

	gormLogLevel := gormLogger.Silent
	if gormLogMode {
		gormLogLevel = gormLogger.Info
	}

	db, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		DBErr = err
		logger.Fatalf("DB Conn Error: %v", err)
		return err
	}

	if !debug {
		err := db.Use(dbresolver.Register(dbresolver.Config{
			Replicas : []gorm.Dialector{
				postgres.Open(replicaDSN),
			},
			Policy: dbresolver.RandomPolicy{},
		}))
		if err != nil {
			DBErr = err
			logger.Fatalf("DB Resolver Error: %v", err)
			return err
		}
	}

	scripts.ExecuteScriptsPriorToMigrations(db)
	logger.Infof("Staring AutoMigrate")
	err = db.AutoMigrate(migrationModels...)
	if err != nil {
		logger.Fatalf("AutoMigrate: %v", err)
		return err
	}
	logger.Infof("AutoMigrate completed successfully")
	scripts.ExecuteScriptsPostToMigrations(db)

	DB = db
	return nil
}	


// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}

//GetDBError returns the database connection error
func GetDBError() error {
	return DBErr
}
