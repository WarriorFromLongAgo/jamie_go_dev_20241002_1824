package main

import (
	"go-project/business/workflow/service"
	"go.uber.org/zap"

	"go-project/main/app"
)

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}
	app.SetConfig(cfg)

	log, err := app.NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	app.SetLogger(log)
	app.Log.Info("InitializeLog success")

	db := app.InitializeDB(cfg)
	app.SetDB(db)
	app.Log.Info("InitializeDB InitializeDB")

	defer func() {
		if db != nil {
			db, _ := db.DB()
			err = db.Close()
			if err != nil {
				app.Log.Error("main db.Close", zap.Any("err", err))
				return
			}
			app.Log.Info("main db.Close success")
		}
	}()

	service.NewService(app.Log, app.DB)

	app.RunServer(cfg)
}
