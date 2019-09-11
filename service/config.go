package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var data map[string]interface{}

const auth = "auth.json"
const services = "services.json"
const lastHeartbeat = "lastheartbeat.txt"

type Auth struct {
	Key  string `json:"key"`
	Host string `json:"host"`
}

type Mapper struct {
	ServiceMap []*ServiceMap
}

type Comms struct {
	Heartbeat time.Time
}

func ParseServices() (*Mapper, error) {

	jsonFile, err := os.Open("services.json")
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened services.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var mapper *Mapper
	err = json.Unmarshal(byteValue, &mapper)
	if err != nil {
		return nil, err
	}

	return mapper, nil
}

func Config(host, apikey string) (*Auth, error) {

	if len(host) > 0 && len(apikey) == 32 {
		return &Auth{Key: apikey, Host: host}, nil
	}

	return ParseAuth()
}

func ParseAuth() (*Auth, error) {
	jsonFile, err := os.Open("auth.json")
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened auth.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var auth *Auth
	err = json.Unmarshal(byteValue, &auth)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func LastHeartbeat() (time.Time, error) {
	return time.Now(), errors.New("Time unknown")
}

func HTTPGetConfig() {

}
