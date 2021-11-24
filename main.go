package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	webHandler "bwastartup/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/gostartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService(campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	webUserHandler := webHandler.NewUserHandler(userService)
	webCampaignHandler := webHandler.NewCampaignHandler(campaignService, userService)

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.HTMLRender = loadTemplates("./web/templates")
	router.Use(cors.New(config))

	//for accessing file
	router.Static("/images", "./images")
	router.Static("/js", "./web/assets/js")
	router.Static("/css", "./web/assets/css")
	router.Static("/webfonts", "./web/assets/webfonts")

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/check-email", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/user/fetch", authMiddleware(authService, userService), userHandler.FetchCurrentUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-image", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/webhook", transactionHandler.GetWebhook)

	// Router for CMS web
	router.GET("/users", webUserHandler.Index)
	router.GET("/users/create", webUserHandler.Create)
	router.GET("/users/edit/:id", webUserHandler.Edit)
	router.GET("/users/avatar/:id", webUserHandler.EditAvatar)
	router.POST("/users", webUserHandler.Store)
	router.POST("/users/update/:id", webUserHandler.Update)
	router.POST("/users/avatar/:id", webUserHandler.UploadAvatar)

	router.GET("/campaigns", webCampaignHandler.Index)
	router.GET("/campaigns/create", webCampaignHandler.Create)
	router.GET("/campaigns/image/:id", webCampaignHandler.UploadImage)
	router.GET("/campaigns/edit/:id", webCampaignHandler.Edit)
	router.GET("/campaigns/show/:id", webCampaignHandler.Show)
	router.POST("/campaigns", webCampaignHandler.Store)
	router.POST("/campaigns/image/:id", webCampaignHandler.StoreImage)
	router.POST("/campaigns/update/:id", webCampaignHandler.Update)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// check header contains word "Bearer"
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// split header to get token
		splittedHeader := strings.Split(authHeader, " ")
		tokenString := ""
		if len(splittedHeader) == 2 {
			tokenString = splittedHeader[1]
		}

		// validate token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get claims and return with mapClaims
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// convert user id in claims data to int
		userId := int(claim["user_id"].(float64))

		// get user by id
		user, err := userService.GetUserByID(userId)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// set user data to context
		c.Set("currentUser", user)

	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
