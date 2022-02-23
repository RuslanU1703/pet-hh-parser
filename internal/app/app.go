package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pet-app/internal/data"
	"pet-app/internal/data/handler"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	TokenAdmin      string `mapstructure:"TOKEN_ADMIN"`
	TokenVK         string `mapstructure:"TOKEN_VK"`
	StringGroupID   string `mapstructure:"GROUP_VK_ID"`
	StringProfileID string `mapstructure:"PROFILE_VK_ID"`
	GroupID         int
	ProfileID       int
}

var (
	config              Config
	callFrequencyTime   = time.Duration(3 * time.Hour)
	count               = 0 // capture first search result
	userAgent           = "my-api-agent"
	searchURLs          = make(map[int]data.SearchURL)
	baseSearchResultMap = make(map[string]data.CastomVacancy)
	allChanges          = make(map[string]interface{})
)

func init() {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../.") // <- path for testing
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err.Error())
	}
	config.GroupID, err = strconv.Atoi(config.StringGroupID)
	if err != nil {
		log.Fatal(err.Error())
	}
	config.ProfileID, err = strconv.Atoi(config.StringProfileID)
	if err != nil {
		log.Fatal(err.Error())
	}
}
func Run() {
	handler := handler.New(searchURLs, &count, config.TokenAdmin)
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler.Init(),
	}
	go func() {
		for {
			switch len(searchURLs) {
			case 0:
				log.Println("[app]No urls")
			default:
				if count == 0 {
					// first search on this urls
					log.Println("[app]Search with updated parameters")
					baseSearchResultMap = getData(searchURLs)
				} else {
					searchResultMap := getData(searchURLs)
					ok := reflect.DeepEqual(baseSearchResultMap, searchResultMap)
					if !ok {
						fixChanges(searchResultMap, baseSearchResultMap)
						baseSearchResultMap = searchResultMap
					}
				}
			}
			count++
			time.Sleep(callFrequencyTime)
		}
	}()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("[app]Server forced to shutdown:", err)
	}
	log.Print("[app]Server stopped") // correctly
}
