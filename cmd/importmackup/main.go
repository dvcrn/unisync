package main

import (
	"fmt"
	"github.com/dvcrn/unisync/cmd/importmackup/generate"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

const mackupPath = "./mackup/mackup/applications"
const targetPath = "./apps/mackup/"

func gen(configPath string) error {
	fmt.Println("parsing: ", configPath)
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	cfg, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"configuration_files"},
	}, content)
	if err != nil {
		fmt.Println("was not able to parse", configPath)
		return nil
	}

	appSection, err := cfg.GetSection("application")
	if err != nil {
		return err
	}
	appName, err := appSection.GetKey("name")
	if err != nil {
		return err
	}
	configurationFiles := cfg.Section("configuration_files").Body()

	generatedConfig, err := generate.GenerateConfig(appName.String(), strings.Split(configurationFiles, "\n"))
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(generatedConfig)
	if err != nil {
		return err
	}

	targetFile := path.Join(targetPath, generatedConfig.FriendlyName+".yaml")
	if err := ioutil.WriteFile(targetFile, out, 0644); err != nil {
		return err
	}

	fmt.Println("generated config for: ", appName.String(), "path: ", targetFile)

	return nil
}

func main() {
	files, err := ioutil.ReadDir(mackupPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fullPath := path.Join(mackupPath, f.Name())
		fmt.Println(fullPath)
		if err := gen(fullPath); err != nil {
			panic(err)
		}
	}
}
