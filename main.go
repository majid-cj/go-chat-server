package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/majid-cj/go-chat-server/config"
	"github.com/majid-cj/go-chat-server/router"

	"github.com/iris-contrib/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}
	os.Mkdir(os.Getenv("UPLOADS"), os.FileMode(0766))
}

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		appConfig.ErrChan <- err
	}

	defer appConfig.Persistence.Client.Disconnect(appConfig.AppContext)
	defer appConfig.Auth.DB.Close()

	CORS := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           5000,
	})

	appConfig.App.Logger().SetLevel("debug")
	appConfig.App.Use(logger.New())
	appConfig.App.Use(recover.New())
	appConfig.App.UseRouter(CORS)

	appConfig.App.HandleDir("/uploads", iris.Dir("./uploads"))
	appConfig.App.I18n.Load("./locales/*/*.yaml")
	appConfig.App.I18n.SetDefault("en")

	router.APIVersionOne(appConfig)

	go func() {
		for range appConfig.ErrChan {
			appConfig.Log.Errorf("%+v", <-appConfig.ErrChan)
		}
	}()

	go func() {
		err := appConfig.App.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")), iris.WithOptimizations)
		if err != nil {
			appConfig.Log.Errorf("Error starting server %+v", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	appConfig.Log.Error("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	appConfig.App.Shutdown(ctx)
}
