package config

import (
	"encoding/json"
	"flag"
	"os"
)

// Configuration ...
var Configuration *Config

// Config ...
type Config struct {
	Database       DatabaseConfiguration `json:"database"`
	ServerPort     string                `json:"serverPort"`
	LevelLog       string                `json:"logLevel"`
	GmailUser      string                `json:"gmail_user"`
	GmailPassword  string                `json:"gmail_password"`
	AccessKey      string                `json:"access_key"`
	RefreshKey     string                `json:"refresh_key"`
	GoogleClient   string                `json:"google_client_id"`
	GoogleSecret   string                `json:"google_secret"`
	FacebookAppId  string                `json:"facebook_appid"`
	FacebookSecret string                `json:"facebook_secret"`
	Contents       string                `json:"contents"`
	ContentsHttps  string                `json:"contents_https"`
	FrontendURL    string                `json:"frontend_url"`
	TLS            bool
	Local          bool
	Host           string
	File           string
	FrontHost      string
	AdminHost      string
}

// DatabaseConfiguration ...
type DatabaseConfiguration struct {
	DBName     string `json:"dbname"`
	DBUser     string `json:"dbuser"`
	DBPassword string `json:"dbpassword"`
	DBPort     string `json:"dbport"`
	DBHost     string `json:"dbhost"`
}

func init() {
	var (
		env   string
		tls   bool
		local bool
	)

	flag.StringVar(&env, "env", "dev", "to connect to right environment")
	flag.BoolVar(&tls, "tls", false, "to decide http or https")
	flag.BoolVar(&local, "local", false, "to decide local cert or public cert")

	flag.Parse()

	data, err := Asset(env + "_config.json")
	if err != nil {
		os.Exit(1)
	}

	_ = json.Unmarshal(data, &Configuration)

	Configuration.Host = "api.pugcha.com"
	Configuration.File = "./static/production/"
	Configuration.FrontHost = "pugcha.com"
	Configuration.AdminHost = "admin.pugcha.com"

	Configuration.Local = local

	if Configuration.Local == true {
		Configuration.Host = "dev.api.pugcha.com"
		Configuration.AdminHost = "dev.admin.pugcha.com"
		Configuration.File = "./static/development/"
		Configuration.FrontHost = "dev.pugcha.com"
	}

	if tls == true {
		Configuration.TLS = true
		Configuration.ServerPort = ":443"
		Configuration.Contents = Configuration.ContentsHttps
	}
}
