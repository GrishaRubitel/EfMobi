package main

import (
	httpH "choomandco/efimobi/http_handler"
	models "choomandco/efimobi/models"
	"errors"

	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var loger = logrus.New()
var tokenAddress string
var nullTable bool = true

func main() {
	envMap, err := godotenv.Read("../../.env")
	if err != nil {
		loger.Fatal("Error loading .env file")
		return
	}

	tokenAddress = envMap["TOKEN_FINDER_ADDRESS_OUTER"]

	conn := envMap["POSTGRES_CONN"]
	fmt.Print(conn)
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		loger.Fatal("DB connection failed: ", err)
		return
	}

	models.MigrateModels(db)

	router := gin.Default()
	router.Use(cors.Default())

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/ping", func(c *gin.Context) {
			code, resp, err := httpH.CallTokenAPI(tokenAddress + "/token")
			if err != nil {
				httpH.ResponseReturner(code, "", err, c)
			} else {
				loger.Info("everything is fine, token - " + resp)
				httpH.ResponseReturner(http.StatusOK, "sounds good", nil, c)
			}
		})

		apiGroup.POST("/execute_dml", func(c *gin.Context) {
			if nullTable {
				resp := executeSQLFile(db, envMap["DML_FILE"])
				if resp == "" {
					httpH.ResponseReturner(http.StatusInternalServerError, "", errors.New("internal server error"), c)
				} else {
					httpH.ResponseReturner(http.StatusCreated, resp, nil, c)
				}
			} else {
				httpH.ResponseReturner(http.StatusConflict, "", errors.New("dml already executed"), c)
			}
		})

		soundGroup := apiGroup.Group("/sound")
		{
			soundGroup.GET("/info", func(c *gin.Context) {
				queryParams := httpH.ReadQueryParams(c)
				code, resp, err := getSoundData(db, queryParams)
				httpH.ResponseReturner(code, resp, err, c)
			})

			soundGroup.GET("/whole_lib", func(c *gin.Context) {
				queryParams := httpH.ReadQueryParams(c)
				code, resp, err := getAllSoundsData(db, queryParams)
				httpH.ResponseReturner(code, resp, err, c)
			})

			soundGroup.PATCH("/delete", func(c *gin.Context) {
				code, bodyParams, err := httpH.ReadBodyData(c)
				if err != nil || code != http.StatusOK {
					httpH.ResponseReturner(code, "", err, c)
				} else {
					code, resp, err := patchDeleteSound(db, bodyParams)
					httpH.ResponseReturner(code, resp, err, c)
				}
			})

			soundGroup.POST("/add_track", func(c *gin.Context) {
				code, bodyParams, err := httpH.ReadBodyData(c)
				if err != nil || code != http.StatusOK {
					httpH.ResponseReturner(code, "", err, c)
				} else {
					code, resp, err := postNewSound(db, bodyParams)
					httpH.ResponseReturner(code, resp, err, c)
				}
			})

			soundGroup.GET("/lyrics", func(c *gin.Context) {
				queryParams := httpH.ReadQueryParams(c)
				code, resp, err := getSoundLyrics(db, queryParams)
				httpH.ResponseReturner(code, resp, err, c)
			})

			soundGroup.PATCH("/update", func(c *gin.Context) {
				code, bodyParams, err := httpH.ReadBodyData(c)
				if err != nil || code != http.StatusOK {
					httpH.ResponseReturner(code, "", err, c)
				} else {
					code, resp, err := patchExistingSound(db, bodyParams)
					httpH.ResponseReturner(code, resp, err, c)
				}
			})
		}

		artistGroup := apiGroup.Group("/artist")
		{
			artistGroup.GET("/info", func(c *gin.Context) {
				queryParams := httpH.ReadQueryParams(c)
				code, resp, err := getArtistData(db, queryParams)
				httpH.ResponseReturner(code, resp, err, c)
			})
		}
	}

	router.Run(envMap["MAIN_SERVER_ADDRESS"])
}
