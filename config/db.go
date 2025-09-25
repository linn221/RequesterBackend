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
	"gorm.io/driver/sqlite"
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

	// Use SQLite by default
	// return connectSQLite()
	return ConnectMySQL()
	// connectRedis()
}
func GetBaseDir() string {
	return _BASE_DIR
}

func connectMySQL() *gorm.DB {
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
		panic("Fail To Connect MySQL Database")
	}
	migrate(db)
	return db
}

func connectSQLite() *gorm.DB {
	// Use app.db as the SQLite database file
	dbPath := "app.db"

	var err error
	db, err := gorm.Open(sqlite.Open(dbPath), initConfig())

	if err != nil {
		panic("Fail To Connect SQLite Database")
	}
	migrate(db)
	return db
}

// ConnectMySQL can be called to use MySQL instead of SQLite
func ConnectMySQL() *gorm.DB {
	// Load env from .env
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

	return connectMySQL()
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
	// Log to both file and standard output
	multiWriter := io.MultiWriter(f, os.Stdout)
	newLogger := logger.New(log.New(multiWriter, "\r\n", log.LstdFlags), logger.Config{
		Colorful:      true,
		LogLevel:      logger.Info, // This will show SQL queries
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
