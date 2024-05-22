package service

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type APIConfig struct {
	Path          string `json:"path"`
	Method        string `json:"method"`
	ThirdPartyURL string `json:"third_party_url"`
}

type Config struct {
	APIs []APIConfig `json:"apis"`
}

var (
	config *Config
	mu     sync.Mutex
)

func loadConfig() (*Config, error) {
	data, err := os.ReadFile("apis.json")
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func createHandler(thirdPartyURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := http.NewRequest(c.Request.Method, thirdPartyURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call third party API"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

func updateRoutes(router *gin.Engine) {
	mu.Lock()
	defer mu.Unlock()

	newConfig, err := loadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	newRouter := gin.Default()

	for _, api := range newConfig.APIs {
		switch api.Method {
		case http.MethodGet:
			newRouter.GET(api.Path, createHandler(api.ThirdPartyURL))
		case http.MethodPost:
			newRouter.POST(api.Path, createHandler(api.ThirdPartyURL))
		default:
			log.Printf("Unsupported method %s for path %s", api.Method, api.Path)
		}
	}

	router = newRouter
	config = newConfig
}

func LoadThirdApis(router *gin.Engine) {
	config, _ = loadConfig()
	updateRoutes(router)

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalf("Failed to create watcher: %v", err)
		}
		defer watcher.Close()

		err = watcher.Add("apis.json")
		if err != nil {
			log.Fatalf("Failed to watch file: %v", err)
		}

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Config file modified, reloading...")
					updateRoutes(router)
				}
			case err := <-watcher.Errors:
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

}
