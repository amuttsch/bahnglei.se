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

type StopPosition struct {
	gorm.Model
	StationID uint
	Platform  string
	Lat       float64
	Lng       float64
	Neighbors string
}

type Station struct {
	gorm.Model
	Country      country.Country
	CountryID    string
	Name         string
	Lat          float64
	Lng          float64
	Operator     string
	Wikidata     string
	Wikipedia    string
	Tracks       int
	Platforms    []Platform
	StopPosition []StopPosition
	OsmTile      []byte
}

type Repo interface {
	Save(i *Station) error
	Get(id uint) *Station
	Search(term string) []Station
	Count() int64
}

func New(db *gorm.DB, ctx context.Context) *stationRepo {
	db.AutoMigrate(&Station{}, &Platform{}, &StopPosition{})

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
	r.db.WithContext(r.ctx).Preload("Platforms").Preload("StopPosition").First(&station, id)
	return &station
}

func (r *stationRepo) Search(term string) []Station {
	var stations []Station
	r.db.WithContext(r.ctx).Where("name ILIKE ?", "%"+term+"%").Order("tracks DESC").Find(&stations)
	return stations
}

func (r *stationRepo) Count() int64 {
	var count int64
	r.db.WithContext(r.ctx).Table("stations").Count(&count)
	return count
}
