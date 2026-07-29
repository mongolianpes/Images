package main

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func MarshalUser(u User) ([]byte, error) {
	return json.Marshal(u)
}

func UnmarshalUser(data []byte) (User, error) {
	var u User
	err := json.Unmarshal(data, &u)
	return u, err
}

func JsoniterMarshalUser(u User) ([]byte, error) {
	return jsonIter.Marshal(u)
}

func JsoniterUnmarshalUser(data []byte) (User, error) {
	var u User
	err := jsonIter.Unmarshal(data, &u)
	return u, err
}
