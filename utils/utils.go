package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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
	GH_TOKEN	string  `json:"gh_token"`
}
type CommandJson struct {
	START    string `json:"start"`
	HELP     string `json:"help"`
	ADD     string `json:"add"`
	ALL		string	`json:"all"`
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

func GetGhToken() string{
	return Config.GH_TOKEN
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

func GetAllCommand() string {
	return CommandConfig.ALL
}

func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}