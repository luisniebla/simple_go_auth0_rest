package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string            `gorm:"string" json:"Name"` // will not export if not cpaitalized
	email         string            `json:"email"`
	user_metadata map[string]string `json:"user_metadata"`
	picture       string            `json:"picture"`
	user_id       string            `json:"user_id"`
}

func userInfoFromAuth0(authHeader string) (map[string]interface{}, error) {
	url := "https://dev-fteqbjgrbz4fpbco.us.auth0.com/api/v2/users/google-oauth2%7C114128656064441920176"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil) // http.Get doesn't allow custom headers
	if err != nil {
		return nil, err
	}

	bearerToken := strings.Split(authHeader, "Bearer ")[1]
	fmt.Println("Bearer", bearerToken)
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// interface{} can hold values of any type, so map looks like {string: interface}
	var userInfo map[string]interface{} // interfaces are sets of method signatures. value can hold any value that implements those methods
	json.Unmarshal(body, &userInfo)
	return userInfo, nil
}
