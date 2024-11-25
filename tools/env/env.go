package env

import (
	"os"
	"strconv"
)

type Configuration struct {
	Port          int    `json:"port"`
	MailServerURL string `json:"mailServerUrl"`
	MongoURL          string `json:"mongoUrl"`

}

func new() *Configuration {
	return &Configuration{
		Port:          3010, // Default port
		MailServerURL: "http://localhost:2525", // Default Mail Server URL
		MongoURL:          "mongodb://localhost:27017",
	}
}

func load() *Configuration {
	result := new()

	if value := os.Getenv("MAIL_SERVER_URL"); len(value) > 0 {
		result.MailServerURL = value
	}

	if value := os.Getenv("PORT"); len(value) > 0 {
		if intVal, err := strconv.Atoi(value); err == nil {
			result.Port = intVal
		}
	}

	return result
}

// Get retrieves the configuration, loading it if necessary
func Get() *Configuration {
	if config == nil {
		config = load()
	}
	return config
}

var config *Configuration
