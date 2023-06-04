package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
	"vehicles/packages/infrastructure/datastore"
	ir "vehicles/packages/infrastructure/router"

	"github.com/gin-gonic/gin"
)

func Run(port string) error {
	redisSearchDB := datastore.CreateNewSearchRDB()
	redisSelectionDB := datastore.CreateNewSelectionRDB()
	surveyDB := datastore.CreateNewDBForSurvey()
	preferencesDB := datastore.CreateNewDBForPreferences()
	vehiclesDB := datastore.CreateNewDBForVehicles()

	router := gin.Default()
	router.LoadHTMLGlob("../server/pages/*html")
	router.Static("/styles", "../server/pages/styles")
	router.Static("/scripts", "../server/pages/scripts")

	car_photos_folders := []string{"polo", "megane", "avensis", "rio", "niva", "hover_h5", "freelander", "octavia", "mondeo", "7-series", "lancer", "antara"}
	numberOfFoldersWithPhotos := 12
	for i := 0; i < numberOfFoldersWithPhotos; i++ {
		dir := http.Dir("../server/pages/car_photos/" + car_photos_folders[i])
		num := strconv.Itoa(i + 1)
		router.StaticFS("/static"+num, dir)
	}

	router = ir.NewRouter(router, redisSearchDB, redisSelectionDB, surveyDB, preferencesDB, vehiclesDB)

	httpServer := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    1000 * time.Second,
		WriteTimeout:   1000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return httpServer.Shutdown(ctx)
}
