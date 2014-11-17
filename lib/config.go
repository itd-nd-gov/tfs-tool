package lib

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ecornell/tfs-tool/utils"
)

type configT struct {
	BaseURL    string
	User       string
	Password   string
	Collection string
}

var Config configT

func LoadConfig() {
	_, err := toml.DecodeFile("tfs-tool.cfg", &Config)
	utils.StopOnErr(err)

	if !strings.HasPrefix(Config.Password, "~~") {
		ciphertext, err := utils.Encrypt([]byte(Config.Password))
		utils.StopOnErr(err)

		Config.Password = "~~" + hex.EncodeToString(ciphertext)

		var configBuffer bytes.Buffer
		toml.NewEncoder(&configBuffer).Encode(Config)

		err = ioutil.WriteFile("tfs-tool.cfg", configBuffer.Bytes(), 0644)
		utils.StopOnErr(err)

	} else {
		encPassword, _ := hex.DecodeString(Config.Password[2:])
		result, err := utils.Decrypt(encPassword)
		utils.StopOnErr(err)
		Config.Password = string(result)
	}

}

func getUser() string {
	if Flags.UserID != "" {
		return Flags.UserID
	}
	return Config.User
}

func getPassword() string {
	if Flags.Password != "" {
		return Flags.Password
	}
	return Config.Password
}
