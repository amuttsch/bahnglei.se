package importer

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/repo/country"
	"gorm.io/gorm"
)

type importerRepo struct {
	db  *gorm.DB
	ctx context.Context
}

type ImporterModel struct {
	gorm.Model
	Country         country.Country
	CountryID       string
	NumberStations  int64
	NumberPlatforms int64
	State           string
}

type Repo interface {
	Save(i *ImporterModel) error
}

func New(db *gorm.DB, ctx context.Context) *importerRepo {
	db.AutoMigrate(&ImporterModel{})

	return &importerRepo{
		db:  db,
		ctx: ctx,
	}
}

func (r *importerRepo) Save(i *ImporterModel) error {
	result := r.db.WithContext(r.ctx).Save(&i)
	return result.Error
}
