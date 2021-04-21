package Techdata

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"TTApi/Config"

	"github.com/spf13/viper"
)

/*
	General API Request function
*/
func ApiRequest(token string, SOIN string, Signature string, Timestamp string, url string, method string) []byte {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	//Header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("SOIN", SOIN)
	req.Header.Add("Signature", Signature)
	req.Header.Add("TimeStamp", Timestamp)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	//Check for any redirects so we can still keep the same headers
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

/*
	Oauth function to generate token
*/
func Oauth(client_id string, client_secret string) []byte {
	//Body for the request
	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", Config.OauthURL, strings.NewReader(data.Encode())) // URL-encoded payload

	if err != nil {
		log.Fatal(err)
	}

	//Header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//Error handling
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

/*
	Login function being called from the API {URL}/techdata/auth
*/
func AuthLogin(w http.ResponseWriter, r *http.Request) {
	authenticator := Oauth(Config.Client_ID, Config.Client_Secret)

	var auth Auth
	json.Unmarshal([]byte(authenticator), &auth)

	if auth.Result != nil {
		log.Fatal(auth.ErrorMessage)
	}

	timeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	//Create the signature
	data := auth.Access_token + ":" + Config.SOIN + ":" + timeStamp
	signature := b64.StdEncoding.EncodeToString([]byte(data))

	//Write the yml file
	if Config.LiveMode {
		viper.SetDefault("techdata_live.timestamp", timeStamp)
		viper.SetDefault("techdata_live.token", auth.Access_token)
		viper.SetDefault("techdata_live.signature", signature)
	} else {
		viper.SetDefault("techdata_dev.timestamp", timeStamp)
		viper.SetDefault("techdata_dev.token", auth.Access_token)
		viper.SetDefault("techdata_dev.signature", signature)
	}

	viper.WriteConfig()
	Config.ApplyConfig()

	//Send responds
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		{
			"status": "success",
			"description": "Token Generated",
			"token": "` + auth.Access_token + `",
			"signature": "` + signature + `"
		}
	`))
}

/*
	Struct for the AuthLogin
*/
type Auth struct {
	Access_token string  `json:"access_token"`
	Result       *string `json:"Result"`
	ErrorMessage *string `json:"ErrorMessage"`
	//token_type   *string `json:"token_type"`
	//expires_in   *int    `json:"expires_in"`
}
