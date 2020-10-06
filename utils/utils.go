package utils

import (
"encoding/json"
"io/ioutil"
"log"
"math/rand"
)

const ConfigJsonPath string = "config.json"
const CommandConfigPath = "commands.json"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ConfigJson struct {
	BOT_TOKEN	string `json:"bot_token"`
	OWNER_ID	int    `json:"owner_id"`
	DB_URL		string	`json:"db_url"`
	DB_Name		string	`json:"db_name"`
	DB_Col		string	`json:"db_collection"`
}
type CommandJson struct {
	START    string `json:"start"`
	HELP     string `json:"help"`
	ADD     string `json:"add"`
}

var CommandConfig *CommandJson = InitCommandConfig()
var Config *ConfigJson = InitConfig()

func InitCommandConfig() *CommandJson {
	file, err := ioutil.ReadFile(CommandConfigPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}
	var Config CommandJson
	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal(err)
	}
	return &Config
}

func InitConfig() *ConfigJson {
	file, err := ioutil.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}

	var Config ConfigJson
	err = json.Unmarshal([]byte(file), &Config)
	if err != nil {
		log.Fatal(err)
	}
	return &Config
}

func GetBotToken() string {
	return Config.BOT_TOKEN
}

func GetDbUrl() string {
	return Config.DB_URL
}

func GetDbCollection() string {
	return Config.DB_Col
}

func GetDbName() string {
	return Config.DB_Name
}

func IsUserOwner(userId int) bool {
	return Config.OWNER_ID == userId
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetStartCommand() string {
	return CommandConfig.START
}

func GetHelpCommand() string {
	return CommandConfig.HELP
}

func GetADDCommand() string {
	return CommandConfig.ADD
}

