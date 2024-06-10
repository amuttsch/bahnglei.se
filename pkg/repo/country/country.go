package country

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type countryRepo struct {
	db  *gorm.DB
	ctx context.Context
}

type Country struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	IsoCode   string         `gorm:"primaryKey"`
	Name      string
	OsmUrl    string
}

type Repo interface {
	Save(country Country) error
	Count() int64
}

func New(db *gorm.DB, ctx context.Context) *countryRepo {
	db.AutoMigrate(&Country{})
	db.Migrator().CreateIndex(&Country{}, "iso_code")

	return &countryRepo{
		db:  db,
		ctx: ctx,
	}
}

func (c *countryRepo) Save(country Country) error {
	result := c.db.WithContext(c.ctx).Save(&country)
	return result.Error
}

func (c *countryRepo) Count() int64 {
	var count int64
	c.db.WithContext(c.ctx).Table("countries").Count(&count)
	return count
}
