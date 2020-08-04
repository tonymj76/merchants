package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Router returns the api resources
func Router() *gin.Engine {
	service := new(Service)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	merchantAPI := router.Group("/api/v1/merchants")
	{
		merchantAPI.POST("/", service.CreateMerchant)
		merchantAPI.GET("/", service.ListMerchants)
		merchantAPI.PATCH("/", service.UpdateMerchant)
		merchantAPI.GET("/:merchant_id", service.GetMerchantByID)
		merchantAPI.DELETE("/:merchant_id", service.DeleteMerchantByID)
	}

	outletAPI := router.Group("api/v1/outlets")
	{
		outletAPI.POST("/", service.CreateMerchantOutlet)
		outletAPI.PATCH("/", service.UpdateMerchantOutlet)
		outletAPI.GET("/:merchant_id", service.GetMerchantOutlets)
		outletAPI.DELETE("/", service.DeleteMerchantOutlet)
	}

	terminalAPI := router.Group("api/v1/terminals")
	{
		terminalAPI.POST("/", service.CreateMerchantTerminal)
		terminalAPI.PATCH("/", service.UpdateMerchantTerminal)
		terminalAPI.GET("/:merchant_id", service.GetMerchantTerminals)
		terminalAPI.DELETE("/", service.DeleteMerchantTerminal)
	}

	return router
}

/*outletAPI := router.Group("api/v1/:merchant_id/outlets")
router.Use(CORSMiddleware())
router.Use(MyCustomMiddleWare())
    router.Use(gin.Logger())
		router.Use(gin.Recovery())
*/
