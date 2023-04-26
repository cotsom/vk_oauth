package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var clientID = ""
var clientSecret = ""
var redirectURI = "http://localhost:8080/me?"
var scope = []string{"account"}
var state = "12345"

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", index)
	r.GET("/me", me)

	r.Run(":8080")
}

func index(c *gin.Context) {
	scopeTemp := strings.Join(scope, "+")
	url := fmt.Sprintf("https://oauth.vk.com/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		clientID,
		redirectURI,
		scopeTemp,
		state,
	)

	c.HTML(http.StatusOK, "index.html", gin.H{"Info": url})
}

func me(c *gin.Context) {
	stateTemp := c.Query("state")

	// if stateTemp[len(stateTemp)-1] == '}' {
	// 	stateTemp = stateTemp[:len(stateTemp)-1]
	// }
	if stateTemp == "" {
		fmt.Errorf("state query param is not provided")
		return
	} else if stateTemp != state {
		fmt.Errorf("state query param do not match original one, got=%s", stateTemp)
		return
	}

	code := c.Query("code")

	if code == "" {
		fmt.Errorf("code query param is not provided")
		return
	}
	url := fmt.Sprintf("https://oauth.vk.com/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s", code, redirectURI, clientID, clientSecret)
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	defer resp.Body.Close()
	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	bytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)
	url = fmt.Sprintf("https://api.vk.com/method/%s?v=5.124&access_token=%s", "", token.AccessToken)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	defer resp.Body.Close()
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	c.HTML(http.StatusOK, "me.html", gin.H{"Info": string(bytes)})
}
