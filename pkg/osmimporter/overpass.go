package osmimporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type ElementType string

const (
	ElementTypeNode     ElementType = "node"
	ElementTypeWay      ElementType = "way"
	ElementTypeRelation ElementType = "relation"
)

type overpassResponse struct {
	OSM3S struct {
		TimestampOSMBase   time.Time `json:"timestamp_osm_base"`
		TimestampAreasBase time.Time `json:"timestamp_areas_base"`
		Copyright          string    `json:"copyright"`
	} `json:"osm3s"`
	Elements []overpassResponseElement `json:"elements"`
	Remark   string                    `json:"remark"`
}

type overpassResponseElement struct {
	Type      ElementType `json:"type"`
	ID        int64       `json:"id"`
	Lat       float64     `json:"lat"`
	Lon       float64     `json:"lon"`
	Timestamp *time.Time  `json:"timestamp"`
	Version   int64       `json:"version"`
	Changeset int64       `json:"changeset"`
	User      string      `json:"user"`
	UID       int64       `json:"uid"`
	Nodes     []int64     `json:"nodes"`
	Members   []struct {
		Type ElementType `json:"type"`
		Ref  int64       `json:"ref"`
		Role string      `json:"role"`
	} `json:"members"`
	Geometry []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"geometry"`
	Bounds *struct {
		MinLat float64 `json:"minlat"`
		MinLon float64 `json:"minlon"`
		MaxLat float64 `json:"maxlat"`
		MaxLon float64 `json:"maxlon"`
	} `json:"bounds"`
	Tags map[string]string `json:"tags"`
}

type Overpass struct {
	ctx            context.Context
	db             *pgxpool.Pool
	repo           *repository.Queries
	client         *http.Client
	interpreterApi string
}

func NewOverpassApi(ctx context.Context, db *pgxpool.Pool, repo *repository.Queries, interpreterApi string) *Overpass {
	client := http.DefaultClient
	return &Overpass{
		ctx:            ctx,
		db:             db,
		repo:           repo,
		interpreterApi: interpreterApi,
		client:         client,
	}
}

func (o *Overpass) fetch(query string) (*overpassResponse, error) {
	resp, err := o.client.PostForm(o.interpreterApi, url.Values{
		"data": {query},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var overpassResp overpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&overpassResp); err != nil {
		return nil, fmt.Errorf("overpass engine error: %w", err)
	}

	if overpassResp.Remark != "" {
		logrus.Warnf("Got remark from API: %s", resp.Remark)
	}

	return &overpassResp, nil
}
