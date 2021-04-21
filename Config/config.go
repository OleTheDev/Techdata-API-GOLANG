package Config

import (
	"github.com/spf13/viper"
)

var LiveMode = true
var Token string
var SOIN string
var Signature string
var Timestamp string
var Client_Secret string
var Client_ID string
var OauthURL string

func ApplyConfig() {
	if LiveMode {
		Token = viper.GetString("techdata_live.token")
		SOIN = viper.GetString("techdata_live.soin")
		Signature = viper.GetString("techdata_live.signature")
		Timestamp = viper.GetString("techdata_live.timestamp")
		Client_Secret = viper.GetString("techdata_live.client_secret")
		Client_ID = viper.GetString("techdata_live.client_id")
		OauthURL = viper.GetString("techdata_live.OauthURL")
	} else {
		Token = viper.GetString("techdata_dev.token")
		SOIN = viper.GetString("techdata_dev.soin")
		Signature = viper.GetString("techdata_dev.signature")
		Timestamp = viper.GetString("techdata_dev.timestamp")
		Client_Secret = viper.GetString("techdata_dev.client_secret")
		Client_ID = viper.GetString("techdata_dev.client_id")
		OauthURL = viper.GetString("techdata_dev.OauthURL")
	}
}
