package commonconfig

import (
	"log"
	"os"
	"path/filepath"
	"system/internal/commonmodel"

	"gopkg.in/yaml.v3"
)

func OpenYaml() (data []byte) {
	basePath, _ := os.Getwd()
	data, err := os.ReadFile(filepath.Join(basePath, "config.yaml"))
	if err != nil {
		log.Fatalf("Read File config.yaml Failed: %v", err)
	}
	return data
}
func GetCommonConfig(data []byte) commonmodel.CommonConfig {
	newConfig := struct {
		CommonConfig commonmodel.CommonConfig `yaml:"common_config"`
	}{}
	if err := yaml.Unmarshal(data, &newConfig); err != nil {
		log.Fatalf("Unmarshal Common Config Failed: %v", err)
	}
	return newConfig.CommonConfig
}
