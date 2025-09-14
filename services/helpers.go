package services

import "gorm.io/gorm"

func first[T any](db *gorm.DB, id int) (*T, error) {
	var v T
	if err := db.First(&v, id).Error; err != nil {
		return nil, err
	}

	return &v, nil
}
