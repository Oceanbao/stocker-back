package infra

import (
	"cmp"
	"slices"

	"example.com/stocker-back/internal/screener"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type ScreenRepositoryPB struct {
	pb *pocketbase.PocketBase
}

func NewScreenRepositoryPB(pb *pocketbase.PocketBase) *ScreenRepositoryPB {
	return &ScreenRepositoryPB{
		pb: pb,
	}
}

type RecordScreen struct {
	Ticker string  `db:"ticker" json:"ticker"`
	Kdj    float64 `db:"kdj" json:"kdj"`
}

func (r RecordScreen) ToMap() map[string]any {
	return map[string]any{
		"ticker": r.Ticker,
		"kdj":    r.Kdj,
	}
}

func (r RecordScreen) ToModel() screener.Screen {
	return screener.Screen{
		Ticker: r.Ticker,
		Kdj:    r.Kdj,
	}
}

func (repo *ScreenRepositoryPB) GetScreens() ([]screener.Screen, error) {
	records, err := repo.pb.Dao().FindRecordsByExpr("screen")
	if err != nil {
		return []screener.Screen{}, err
	}

	screens := make([]screener.Screen, 0, len(records))
	for idx := range records {
		s := screener.Screen{
			Ticker: records[idx].GetString("ticker"),
			Kdj:    records[idx].GetFloat("kdj"),
		}
		screens = append(screens, s)
	}

	// Sort ascending.
	slices.SortFunc(screens, func(a, b screener.Screen) int {
		return cmp.Compare(a.Kdj, b.Kdj)
	})

	return screens, nil
}

func (repo *ScreenRepositoryPB) SetScreens(screens []screener.Screen) error {
	// First clear all records.
	records, err := repo.pb.Dao().FindRecordsByExpr("screen")
	if err != nil {
		return err
	}

	for _, rec := range records {
		if err := repo.pb.Dao().DeleteRecord(rec); err != nil {
			return err
		}
	}

	collection, err := repo.pb.Dao().FindCollectionByNameOrId("screen")
	if err != nil {
		return err
	}

	err = repo.pb.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, data := range screens {
			recordData, err := data.ToMap()
			if err != nil {
				return err
			}
			record := models.NewRecord(collection)
			record.Load(recordData)

			err = txDao.SaveRecord(record)
			if err != nil {
				repo.pb.Logger().Error("cannot write to `screen`", "error", err.Error())
				continue
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
