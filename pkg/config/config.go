package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// loadConfig - функция заполняющая структуру конфига (structConfg)
// из файла (Path) или возвращающая ошибку если это сделать не удалось
func LoadConfig(Path string, structConfg interface{}) error {
	configBytes, err := os.ReadFile(Path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configBytes, structConfg)
	if err != nil {
		return err
	}
	return nil
}
