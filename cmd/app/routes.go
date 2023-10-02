package main

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func routeTrack(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/track", func(c echo.Context) error {
		// 1. Get all records from Collection `track`: 'code, name, started'.
		var tempRecords = []struct {
			ID      string `db:"id" json:"id"`
			Code    string `db:"code" json:"code"`
			Name    string `db:"name" json:"name"`
			Started string `db:"started" json:"started"`
		}{}
		err := app.Dao().DB().
			Select("id", "code", "name", "started").
			From("track").
			All(&tempRecords)
		if err != nil {
			return c.JSON(500, map[string]string{"message": "error in getting all records from database `track`"})
		}
		if len(tempRecords) == 0 {
			return c.JSON(200, map[string]string{"message": "ok", "data": ""})
		}

		return c.JSON(200, map[string]any{"message": "ok", "data": tempRecords})

		// if key := c.PathParam("key"); key != keyStored {
		// 	return c.JSON(http.StatusForbidden, map[string]string{"message": "Not allowed."})
		// }

		// var err error
		// // 1. Update all `alert` records.
		// // 1.1 Get all `stocks` records.
		// var tempRecords = []struct {
		// 	ID   string `db:"id" json:"id"`
		// 	Code string `db:"code" json:"code"`
		// 	Name string `db:"name" json:"name"`
		// 	Cap  string `db:"cap" json:"cap"`
		// }{}
		// err = app.Dao().DB().
		// 	Select("id", "code", "name", "cap").
		// 	From("stocks").
		// 	All(&tempRecords)
		// if err != nil {
		// 	log.Println("error in reading database `stocks`")
		// 	return err
		// }
		// // 1.2 For each `alert` record, update its fields.
		// err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		// 	for _, x := range tempRecords {
		// 		record, errFindStockByID := app.Dao().FindRecordById("stocks", x.ID)
		// 		if errFindStockByID != nil {
		// 			log.Println("error in finding record in 'stocks': ", errFindStockByID)
		// 			return errFindStockByID
		// 		}

		// 		codeNew := x.Code
		// 		codeNew = strings.ReplaceAll(codeNew, "sh", "1.")
		// 		codeNew = strings.ReplaceAll(codeNew, "sz", "0.")
		// 		record.Set("code", codeNew)

		// 		if err = txDao.SaveRecord(record); err != nil {
		// 			log.Println("error in updating record: ", x.Code, err)
		// 			return err
		// 		}
		// 	}

		// 	return nil
		// })
		// if err != nil {
		// 	log.Println("error in transaction of updating stocks records: ", err)
		// 	return err
		// }
	}, /* optional middlewares */ apis.RequireRecordAuth("users"))
}

// e.Router.POST("/upload", func(c echo.Context) error {
// 	data := struct {
// 		Daily []dailyData `json:"daily"`
// 	}{}
// 	if err := c.Bind(&data); err != nil {
// 		log.Println("error in reading body json")
// 		log.Println("error: ", err)
// 		return apis.NewBadRequestError("error in reading body json", err)
// 	}

// 	// var daily []dailyData
// 	// err := json.Unmarshal(dailyFile, &daily)
// 	// if err != nil {
// 	// 	log.Println("error in reading embed file")
// 	//  return apis.NewBadRequestError("erro in reading embed file", err)
// 	// }

// 	// b, err := json.Marshal(daily[0])
// 	// if err != nil {
// 	// 	log.Println("error in marshalling json")
// 	//  return apis.NewBadRequestError("error in marshalling json", err)
// 	// }

// 	collection, err := app.Dao().FindCollectionByNameOrId("daily")
// 	if err != nil {
// 		log.Println("error in finding collection 'daily'", collection)
// 		return apis.NewBadRequestError("error in finding collection `daily`", err)
// 	}

// 	app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
// 		for _, stock := range data.Daily {
// 			for _, entry := range stock.Data {
// 				record := models.NewRecord(collection)
// 				record.Load(map[string]any{
// 					"code":  stock.Code,
// 					"date":  entry.Date,
// 					"open":  entry.Open,
// 					"high":  entry.High,
// 					"low":   entry.Low,
// 					"close": entry.Close,
// 				})

// 				if err = txDao.SaveRecord(record); err != nil {
// 					log.Println("error in writing record: ", stock.Name, entry.Date)
// 					// log.Println("error: ", err)
// 					// return apis.NewBadRequestError(fmt.Sprintf("error in writing record: %v-%v", stock.Name, entry.Date), err)
// 				}
// 			}
// 		}

// 		return nil
// 	})

// 	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
// })

// e.Router.GET("/upload", func(c echo.Context) error {
// 	data := []dailyRecordDB{}

// 	err := app.Dao().DB().
// 		Select("id", "code", "date", "open", "high", "low", "close").
// 		From("daily").
// 		All(&data)
// 	if err != nil {
// 		log.Println("error in reading database")
// 		return apis.NewBadRequestError("error in reading database", err)
// 	}

// 	grouped := lop.GroupBy(data, func(d dailyRecordDB) (code string) {
// 		return d.Code
// 	})
// 	log.Println("len(data): ", len(data))
// 	log.Println("len(grouped): ", len(grouped))

// 	codeDaily := make([]string, len(grouped))
// 	index := 0
// 	for k := range grouped {
// 		var key string
// 		parts := strings.Split(k, ".")
// 		if parts[0] == "0" {
// 			key = fmt.Sprintf("%s%s", "sz", parts[1])
// 		} else {
// 			key = fmt.Sprintf("%s%s", "sh", parts[1])
// 		}
// 		codeDaily[index] = key
// 		index++
// 	}

// 	var dataTemp = []struct {
// 		Code string `db:"code" json:"code"`
// 	}{}
// 	err = app.Dao().DB().
// 		Select("code").
// 		From("stocks").
// 		All(&dataTemp)
// 	if err != nil {
// 		log.Println("error in reading database")
// 		return apis.NewBadRequestError("error in reading database", err)
// 	}
// 	codeStock := make([]string, len(dataTemp))
// 	for idx, x := range dataTemp {
// 		codeStock[idx] = x.Code
// 	}

// 	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
// })
