package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var IP string

var currKeys1 StructKeys
var currKeys2 StructKeys
var keys1 BoxKeysModel
var keys2 BoxKeysModel

func TestInitializing(t *testing.T) {
	fmt.Println("Reading config file")

	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Log("Error occured while reading config")
	}
	json.Unmarshal(raw, &AppConfig)
	IP = AppConfig.IP + AppConfig.Port
	fmt.Printf("IP address is  : %s", IP)
}

func TestRegistration(t *testing.T) {
	jsonData := map[string]string{"yubikey": "akshjdk9823/=2n1923-=123", "mobileKey": "%ja0asbd12+kshjdk9823/=2n1923-=123"}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://192.168.0.108:2020/dacs/u2fRegistration", "application/json", bytes.NewBuffer([]byte(jsonValue)))
	if err != nil {
		t.Errorf("Connection Error : %d", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		type Result struct {
			Status  string
			Message string
		}
		var response Result
		errUnmarshal := json.Unmarshal(data, &response)
		if errUnmarshal != nil {
			t.Error("unmarshelling Error")
		} else {
			if response.Status == "Success" {
				t.Log("Success")
			} else {
				t.Logf(response.Message)
			}
		}
	}
}

func TestAuthenticaation(t *testing.T) {
	jsonData := map[string]string{"yubikey": "akshjdk9823/=2n1923-=123", "mobileKey": "%ja0asbd12+kshjdk9823/=2n1923-=123"}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://192.168.0.108:2020/dacs/u2fAuthentication", "application/json", bytes.NewBuffer([]byte(jsonValue)))
	if err != nil {
		t.Errorf("Connection Error : %d", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		type Result struct {
			Status  string
			Message string
		}
		var response Result
		errUnmarshal := json.Unmarshal(data, &response)
		if errUnmarshal != nil {
			t.Error("unmarshelling Error")
		} else {
			if response.Status == "Success" {
				t.Log("Success")
			} else {
				t.Logf(response.Message)
			}
		}
	}
}
