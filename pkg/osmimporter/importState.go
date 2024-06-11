package osmimporter

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/country"
	"gorm.io/gorm"
)

type importerRepo struct {
	db  *gorm.DB
	ctx context.Context
}

type ImportState struct {
	gorm.Model
	Country         country.Country
	CountryID       string
	NumberStations  int64
	NumberPlatforms int64
	State           string
}

type Repo interface {
	Save(i *ImportState) error
}

func NewRepo(db *gorm.DB, ctx context.Context) *importerRepo {
	db.AutoMigrate(&ImportState{})

	return &importerRepo{
		db:  db,
		ctx: ctx,
	}
}

func (r *importerRepo) Save(i *ImportState) error {
	result := r.db.WithContext(r.ctx).Save(&i)
	return result.Error
}
