package infra

import (
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

func (repo *ScreenRepositoryPB) GetScreenAll() ([]screener.Screen, error) {
	records, err := repo.pb.Dao().FindRecordsByExpr("screen")
	if err != nil {
		return []screener.Screen{}, err
	}

	screens := make([]screener.Screen, len(records))
	for idx := range records {
		screens[idx].Ticker = records[idx].GetString("ticker")
		screens[idx].Kdj = records[idx].GetFloat("kdj")
	}

	return screens, nil
}

func (repo *ScreenRepositoryPB) SetScreenAll(screens []screener.Screen) error {
	collection, err := repo.pb.Dao().FindCollectionByNameOrId("screen")
	if err != nil {
		return err
	}

	err = repo.pb.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, data := range screens {
			recordData := convertScreenToRecord(data).ToMap()
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

func convertScreenToRecord(screen screener.Screen) RecordScreen {
	return RecordScreen{
		Ticker: screen.Ticker,
		Kdj:    screen.Kdj,
	}
}
