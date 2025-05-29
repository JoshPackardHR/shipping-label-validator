package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/JoshuaPackardHR/shipping-label-validator/docs"
	"github.com/JoshuaPackardHR/shipping-label-validator/gpt"
	"github.com/JoshuaPackardHR/shipping-label-validator/internal/shipping"
	"github.com/JoshuaPackardHR/shipping-label-validator/ups"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

//go:generate swag init
//go:generate swag fmt

//	@title			Shipping Label Validator API
//	@version		1.0
//	@description	Public API for Shipping Label Validator
//	@termsOfService	http://happyreturns.com/terms/

//	@contact.name	API Support
//	@contact.url	http://www.happyreturns.com/support
//	@contact.email	dev@happyreturns.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@schemes	http
//	@host		localhost:8080
//	@BasePath	/api/latest

//	@securityDefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file: " + err.Error())
	}

	environment := os.Getenv("APP_ENV")

	// Configure swagger docs
	switch environment {
	case "production":
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.Host = "shipping-label-validator-api.happyreturns.com"
	case "staging":
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.Host = "shipping-label-validator-api-staging.happyreturns.com"
	}

	// Initialize gin router
	if environment != "local" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Enable CORS for all origins
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Authorization,Content-Type,access-control-allow-origin,access-control-allow-headers"},
		AllowCredentials: true,
	}))

	router.GET("/heartbeat", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	latest := router.Group("/api/latest")
	upsClient, err := ups.NewClient(os.Getenv("UPS_CLIENT_ID"), os.Getenv("UPS_CLIENT_SECRET"))
	if err != nil {
		log.Fatalf("Failed to initialize UPS client: %v", err)
	}

	gptClient, err := initGPTClient()
	if err != nil {
		log.Fatalf("Failed to initialize GPT client: %v", err)
	}

	shipping.NewHandler(
		shipping.NewManager(upsClient, gptClient),
	).RegisterRoutes(latest.Group("/shipping"))

	httpPort := ":" + os.Getenv("HTTP_PORT")
	if httpPort == ":" {
		httpPort = ":8080"
	}
	bind := httpPort
	if environment == "local" {
		bind = fmt.Sprintf("127.0.0.1%s", httpPort)
	}

	go func() {
		err := router.Run(bind)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Capture Ctrl-c to shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

}

func initGPTClient() (gpt.GPT, error) {
	if os.Getenv("GPT") == "gemini" {
		gemini, err := gpt.NewGemini(os.Getenv("GEMINI_MODEL"), os.Getenv("GEMINI_API_KEY"))
		if err != nil {
			return nil, err
		}

		return gemini, nil
	}

	openAi, err := gpt.NewOpenAI(os.Getenv("OPENAI_MODEL"), os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		return nil, err
	}

	return openAi, nil
}
