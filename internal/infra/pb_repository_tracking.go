package infra

import (
	"example.com/stocker-back/internal/tracking"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type TrackingRepositoryPB struct {
	pb *pocketbase.PocketBase
}

func NewTrackingRepositoryPB(pb *pocketbase.PocketBase) *TrackingRepositoryPB {
	return &TrackingRepositoryPB{
		pb: pb,
	}
}

type RecordTracking struct {
	Ticker string `db:"ticker" json:"ticker"`
	Name   string `db:"name" json:"name"`
}

func (r RecordTracking) ToMap() map[string]any {
	return map[string]any{
		"ticker": r.Ticker,
		"name":   r.Name,
	}
}

func (r RecordTracking) ToModel() tracking.Tracking {
	return tracking.Tracking{
		Ticker: r.Ticker,
		Name:   r.Name,
	}
}

func (repo *TrackingRepositoryPB) GetTrackings() ([]tracking.Tracking, error) {
	records, err := repo.pb.Dao().FindRecordsByExpr("tracking")
	if err != nil {
		return nil, err
	}

	trackings := make([]tracking.Tracking, len(records))
	for idx := range records {
		trackings[idx].Ticker = records[idx].GetString("ticker")
		trackings[idx].Name = records[idx].GetString("name")
	}

	return trackings, nil
}

func (repo *TrackingRepositoryPB) SetTracking(tracking tracking.Tracking) error {
	collection, err := repo.pb.Dao().FindCollectionByNameOrId("tracking")
	if err != nil {
		return err
	}

	model := models.NewRecord(collection)
	record := convertTrackingToRecord(tracking).ToMap()
	model.Load(record)

	err = repo.pb.Dao().SaveRecord(model)
	if err != nil {
		repo.pb.Logger().Error("cannot write to `tracking`", "error", err.Error())
		return nil
	}

	return nil
}

func convertTrackingToRecord(tracking tracking.Tracking) RecordTracking {
	return RecordTracking{
		Ticker: tracking.Ticker,
		Name:   tracking.Name,
	}
}
