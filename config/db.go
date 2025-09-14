package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const SearchLimit = 10

var _BASE_DIR string

func ConnectDB() *gorm.DB {
	// Load env from .env
	// godotenv.Load(".env")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	_BASE_DIR = dir
	environmentPath := filepath.Join(dir, ".env")
	err = godotenv.Load(environmentPath)
	if err != nil {
		log.Fatal(err)
	}

	return connectDatabase()
	// connectRedis()
}
func GetBaseDir() string {
	return _BASE_DIR
}

func connectDatabase() *gorm.DB {
	databaseConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err := gorm.Open(mysql.Open(databaseConfig), initConfig())

	if err != nil {
		panic("Fail To Connect Database")
	}
	migrate(db)
	return db
}

// InitConfig Initialize Config
func initConfig() *gorm.Config {
	return &gorm.Config{
		Logger:         initLog(),
		NamingStrategy: initNamingStrategy(),
	}
}

// InitLog Connection Log Configuration
func initLog() logger.Interface {
	f, _ := os.Create("gorm.log")
	newLogger := logger.New(log.New(io.MultiWriter(f), "\r\n", log.LstdFlags), logger.Config{
		Colorful:      true,
		LogLevel:      logger.Info,
		SlowThreshold: time.Second,
	})
	return newLogger
}

// InitNamingStrategy Init NamingStrategy
func initNamingStrategy() *schema.NamingStrategy {
	return &schema.NamingStrategy{
		SingularTable: false,
		TablePrefix:   "",
	}
}
