package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"
	//"github.com/kushtaka/kushtakad/service/telnet"
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

func ValidateAuth(host, apikey string) (*Auth, error) {

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

func HTTPServicesConfig(host, key string) ([]interface{}, error) {
	url := host + "/api/v1/config.json"

	spaceClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := spaceClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tmpMap []TmpMap
	err = json.Unmarshal(body, &tmpMap)
	if err != nil {
		return nil, err
	}

	for _, v := range tmpMap {
		switch v.Type {
		case "telnet":
			rf := reflect.ValueOf(v.Service)
			log.Info(rf)
			//rf.FieldByName("Port")
			//log.Info("telnet", rf.Kind())
			//log.Info("telnet", rf.FieldByName("Port"))
			//back, ok := v.Service.(telnet.TelnetService)
			//log.Info("telnet", v.Service, back.ID, ok)
		}
	}

	return nil, nil
}
