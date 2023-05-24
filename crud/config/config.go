package config

import (
	"fmt"

	"github.com/daimall/tools/crud/common"
	"github.com/spf13/viper"
)

func init() {
	rootPath := common.GetRootPath()
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(rootPath + "/conf")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}
}
