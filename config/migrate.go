package config

import (
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	// Then migrate the other tables
	// err := db.AutoMigrate(&requests.Endpoint{}, &requests.ImportJob{}, &requests.MyRequest{})
	// if err != nil {
	// 	panic("Error migrating other tables: " + err.Error())
	// }
}
