package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	CONFIG_FILE_NAME = "rest-db"
	CONFIG_TYPE = "yaml"
	ENV_PREFIX = ""
)

func InitConfig() error {
	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	var cfgPath = os.Getenv("MPC_CFG_PATH")
	if cfgPath != "" {
		viper.AddConfigPath(cfgPath)
	} else {
		viper.AddConfigPath("./")
		gopath := os.Getenv("GOPATH")
		for _, p := range filepath.SplitList(gopath) {
			peerpath := filepath.Join(p, "/src/config")
			viper.AddConfigPath(peerpath)
		}
	}

	viper.SetConfigName(CONFIG_FILE_NAME)
	viper.SetConfigType(CONFIG_TYPE)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}