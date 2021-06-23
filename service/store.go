package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
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
		db, err = gorm.Open(mysql.Open("root:@tcp(127.0.0.1:4000)/test?parseTime=true"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db.AutoMigrate(&Tips{})
	})
	return db
}

func PutNewTip(ctx context.Context, tips string) (string, error) {
	r := &Tips{
		Content: tips,
	}
	err := DB().Create(r).Error
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", r.ID), nil
}

func SearchTip(ctx context.Context, keyword string) (string, error) {
	var (
		tips []Tips
		err  error
	)
	if tracer := ctx.Value("tracer_id").(opentracing.TextMapCarrier); tracer != nil {
		DB().Exec(fmt.Sprintf("set @@session.tracer_id = '%s'", EncodeMap(tracer)))
	}
	result := DB().Where("Content LIKE ?", "%"+keyword+"%").Find(&tips)

	if tracer := ctx.Value("tracer_id").(opentracing.TextMapCarrier); tracer != nil {
		DB().Exec("set @@session.tracer_id = ''")
	}

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
