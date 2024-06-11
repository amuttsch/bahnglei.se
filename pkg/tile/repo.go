package tile

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type tileRepo struct {
	db  *gorm.DB
	ctx context.Context
}

type OsmTile struct {
	gorm.Model
	X    uint `gorm:"index:idx_coordinate"`
	Y    uint `gorm:"index:idx_coordinate"`
	Z    uint `gorm:"index:idx_coordinate"`
	Data []byte
}

type Repo interface {
	Save(i *OsmTile) error
	Get(x uint, y uint, z uint) *OsmTile
}

func NewRepo(db *gorm.DB, ctx context.Context) *tileRepo {
	db.AutoMigrate(&OsmTile{})

	return &tileRepo{
		db:  db,
		ctx: ctx,
	}
}

func (r *tileRepo) Save(o *OsmTile) error {
	result := r.db.WithContext(r.ctx).Save(&o)
	return result.Error
}

func (r *tileRepo) Get(x uint, y uint, z uint) *OsmTile {
	var tile OsmTile
	err := r.db.WithContext(r.ctx).Where("x = ? AND y = ? AND z = ?", x, y, z).First(&tile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &tile
}
