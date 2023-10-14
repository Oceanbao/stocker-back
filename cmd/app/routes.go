package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func routeUpdateDailyETF(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/update-daily-etf", func(c echo.Context) error {
		days := c.QueryParam("days")
		if days == "" {
			return c.JSON(http.StatusOK, map[string]any{"message": "require days param"})
		}

		daysInt, err := strconv.Atoi(days)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"message": "cannot conver query to int"})
		}

		go updateDailyCollectionETF(app, daysInt)

		return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
		// })
	}, /* optional middlewares */ apis.RequireRecordAuth("users"))
}

func routeUpdateDaily(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/update-daily", func(c echo.Context) error {
		go updateDailyCollection(app)

		return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
	}, /* optional middlewares */ apis.RequireRecordAuth("users"))
}

func routeTrack(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/track", func(c echo.Context) error {
		// 1. Get all records from Collection `track`: 'code, name, started'.
		var tempRecords = []recordTrack{}
		err := app.Dao().DB().
			Select("id", "code", "name", "started").
			From("track").
			All(&tempRecords)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{"message": "error in getting all records from database `track`"},
			)
		}
		if len(tempRecords) == 0 {
			return c.JSON(http.StatusOK, map[string]string{"message": "ok", "data": ""})
		}

		return c.JSON(http.StatusOK, map[string]any{"message": "ok", "data": tempRecords})

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
