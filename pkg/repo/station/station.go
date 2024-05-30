package station

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/repo/country"
	"gorm.io/gorm"
)

type stationRepo struct {
	db  *gorm.DB
	ctx context.Context
}

type Platform struct {
	gorm.Model
	StationID uint
	Positions string
}

type Station struct {
	gorm.Model
	Country   country.Country
	CountryID string
	Name      string
	Lat       float64
	Lng       float64
	Operator  string
	Wikidata  string
	Wikipedia string
	Platforms []Platform
}

type Repo interface {
	Save(i *Station) error
	Get(id uint) *Station
}

func New(db *gorm.DB, ctx context.Context) *stationRepo {
	db.AutoMigrate(&Station{}, &Platform{})

	return &stationRepo{
		db:  db,
		ctx: ctx,
	}
}

func (r *stationRepo) Save(i *Station) error {
	result := r.db.WithContext(r.ctx).Save(&i)
	return result.Error
}

func (r *stationRepo) Get(id uint) *Station {
	var station Station
	r.db.WithContext(r.ctx).First(&station, id)
	return &station
}
