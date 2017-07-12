package utility

import (
	"encoding/json"
	"os"
)

// mczal: For Externalized Configuration
var Configuration EntConfiguration

type EntConfiguration struct {
	Secret                 string `json:"secret"`
	Port                   string `json:"port"`
	SendgridKey            string `json:"sendgridKey"`
	BasePathLocation       string `json:"basePathLocation"`
	StaticBasePathLocation string `json:"staticBasePathLocation"`
	NotifWhenPointAchieve  []int  `json:"notifWhenPointAchieve"`
	BonusPoint             int    `json:"bonusPoint"`
	Dynamodb_dbname        string `json:"dynamodb_dbname"`
}

// mczal:  When error, expecting application to stop (panic)
func InitializeConfig() {
	file, err := os.Open("./env.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Configuration)
	if err != nil {
		panic(err)
	}
}

func ContainsInt(intSlice []int, searchInt int) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}
