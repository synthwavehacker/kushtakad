package service

var data map[string]interface{}

// send map["config"]Config
type Config struct {
	ApiKey string
	Host   string
}
