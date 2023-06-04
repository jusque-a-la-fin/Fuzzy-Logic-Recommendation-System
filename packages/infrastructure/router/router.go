package router

import (
	"database/sql"
	"net/http"
	"vehicles/packages/registry"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewRouter(router *gin.Engine, redisSearchDB *redis.Client, redisSelectionDB *redis.Client, surveyDB *sql.DB, preferencesDB *sql.DB, vehiclesDB *sql.DB) *gin.Engine {

	router.GET("main", func(c *gin.Context) {
		registry.NewSearchController(surveyDB, redisSearchDB, c).ShowMainPage(c, "main_page.html")
	})
	router.POST("main", func(c *gin.Context) {
		registry.NewSearchController(surveyDB, redisSearchDB, c).GetCars(c)
		c.Redirect(http.StatusFound, "http://localhost:8080/search")
	})
	router.POST("fingerprint", func(c *gin.Context) { registry.NewUserController(c).SetFingerprint(c) })

	router.GET("search", func(c *gin.Context) {
		registry.NewSearchController(surveyDB, redisSearchDB, c).PassCarsData(c, redisSearchDB, surveyDB)
	})
	router.GET("search/card/:id", func(c *gin.Context) {
		registry.NewSearchController(surveyDB, redisSearchDB, c).ShowCarCard(c, redisSearchDB)
	})
	router.POST("search", func(c *gin.Context) {
		registry.NewQuestionController(surveyDB, redisSearchDB, c).GetAnswer(c, redisSearchDB, surveyDB)
	})

	selection := router.Group("/selection")
	{

		selection.GET("priorities", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).ChoosePriorities(c)
		})

		selection.POST("priorities", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).SetPriorities(c)
			c.JSON(http.StatusOK, gin.H{"message": "Данные успешно получены"})
		})
		selection.GET("price", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).ChoosePrice(c)
		})
		selection.POST("price", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).SetPrice(c)
			c.JSON(http.StatusOK, gin.H{"message": "Данные успешно получены"})
		})
		selection.GET("manufacturers", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).ChooseManufacturers(c)
		})
		selection.POST("manufacturers", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).SetManufacturers(c)
		})
		selection.GET("", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).GetSelection(c)
		})
		selection.GET("card/:id", func(c *gin.Context) {
			registry.NewSelectionController(preferencesDB, vehiclesDB, redisSelectionDB, c).ShowCarCard(c)
		})
	}

	return router
}
