package main

import (
	"aurora/initialize"
	"embed"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/acheong08/endless"
	"github.com/joho/godotenv"
)

//go:embed web/*
var staticFiles embed.FS

func main() {
	_ = godotenv.Load(".env")
	gin.SetMode(gin.ReleaseMode)
	router := initialize.RegisterRouter()
	subFS, err := fs.Sub(staticFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	router.StaticFS("/web", http.FS(subFS))
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	viper.SetConfigFile("/data/options.json")
	viper.SetConfigType("json")

	err1 := viper.ReadInConfig() // 查找并读取配置文件
	if err1 != nil {
		log.Printf("配置文件不存在，请检查: %s \n", err1)
	}

	proxy := viper.GetString("proxy_url")
	auth := viper.GetString("authorization")

	if proxy != "" {
		os.Setenv("PROXY_URL", proxy)
	}

	if auth != "" {
		os.Setenv("Authorization", auth)
	}

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	log.Printf("server running on %s:%s", host, port)
	log.Printf("Authorization key %s, PROXY_URL is %s", os.Getenv("Authorization"), os.Getenv("PROXY_URL"))

	if tlsCert != "" && tlsKey != "" {
		_ = endless.ListenAndServeTLS(host+":"+port, tlsCert, tlsKey, router)
	} else {
		_ = endless.ListenAndServe(host+":"+port, router)
	}
}
