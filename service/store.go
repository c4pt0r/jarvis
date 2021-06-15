package main

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Data model
type Tips struct {
	gorm.Model
	ID      int64 `gorm:"AUTO_INCREMENT`
	Content string
}

var (
	once sync.Once
	db   *gorm.DB
)

func DB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open("javis.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db.AutoMigrate(&Tips{})
	})
	return db
}

func PutNewTip(tips string) (string, error) {
	r := &Tips{
		Content: tips,
	}
	err := DB().Create(r).Error
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", r.ID), nil
}

func SearchTip(keyword string) (string, error) {
	var (
		tips []Tips
		err  error
	)
	result := DB().Where("Content LIKE ?", "%"+keyword+"%").Find(&tips)
	err = result.Error
	if err != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "no such record", nil
		}
		return "", err
	}
	var output string
	for _, r := range tips {
		output += fmt.Sprintf("KB-%d %s\n", r.ID, r.Content)
	}
	return output, nil
}
