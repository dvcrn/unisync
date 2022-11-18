package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dvcrn/unisync/internal"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

//go:embed apps/*
var content embed.FS

type actionType string

const (
	actionTypeUnknown        actionType = ""
	actionTypeSync           actionType = "sync"
	actionTypeInitFromTarget actionType = "initFromTarget"
	actionTypeInitFromLocal  actionType = "initFromLocal"
	actionTypeList           actionType = "list"
)

var action actionType = actionTypeUnknown

func init() {
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  sync - run sync between all enabled apps")
		fmt.Println("  init-from-target - run/force initial sync, targetPath -> appPath")
		fmt.Println("  init-from-local - run/force initial sync, local -> targetPath")
		fmt.Println("  list - list available apps")
	}

	if len(flag.Args()) < 1 {
		log.Println("no argument given. currently supported: sync")
		flag.Usage()
		os.Exit(1)
	}

	switch strings.ToLower(flag.Arg(0)) {
	case "sync":
		action = actionTypeSync
	case "init-from-target":
		action = actionTypeInitFromTarget
	case "init-from-local":
		action = actionTypeInitFromLocal
	case "list":
		action = actionTypeList

	}
}

func readAppConfigs(basePath string) (map[string]*internal.AppConfig, error) {
	configurations := map[string]*internal.AppConfig{}

	entries, err := content.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		c, err := content.ReadFile(fmt.Sprintf("%s/%s", basePath, entry.Name()))
		if err != nil {
			return nil, err
		}

		var t *internal.AppConfig
		err = yaml.Unmarshal(c, &t)
		if err != nil {
			return nil, err
		}

		configurations[t.FriendlyName] = t
	}

	return configurations, nil
}

func loadApps() (map[string]*internal.AppConfig, error) {
	configurations := map[string]*internal.AppConfig{}

	configEntries, err := readAppConfigs("apps/mackup")
	if err != nil {
		return nil, err
	}

	for k, v := range configEntries {
		configurations[k] = v
	}

	// apps overwrite mackup, so mackup gets parsed first
	configEntries, err = readAppConfigs("apps")
	if err != nil {
		return nil, err
	}

	for k, v := range configEntries {
		configurations[k] = v
	}

	return configurations, nil
}

func main() {
	switch action {
	case actionTypeList:
		configurations, err := loadApps()
		if err != nil {
			log.Fatal(err)
		}

		keys := maps.Keys(configurations)
		sort.Strings(keys)

		fmt.Println("Available apps:")
		for _, k := range keys {
			fmt.Printf(" - %s (%s)\n", configurations[k].Name, configurations[k].FriendlyName)
		}

	case actionTypeSync,
		actionTypeInitFromLocal,
		actionTypeInitFromTarget:

		configurations, err := loadApps()
		if err != nil {
			log.Fatal(err)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		configFile, err := os.ReadFile(homeDir + "/.config/unisync/unisync.yaml")
		if err != nil {
			configFile, err = os.ReadFile(homeDir + "/.config/unisync/unisync.yml")
			if err != nil {
				log.Fatal(err)
			}
		}

		var config *internal.Config
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			log.Fatalf("err trying to parse yaml config file: %s\n", err)
		}

		preferDirection := internal.PreferMode(config.PreferDirection)
		if config.PreferDirection != internal.PreferModeLocal && config.PreferDirection != internal.PreferModeTarget {
			fmt.Println("preferDirection not set or invalid, defaulting to 'preferDirection = target'")
			preferDirection = internal.PreferModeTarget
		}

		fmt.Println("starting unisync...")
		fmt.Printf("targetPath: %s, apps: %s\n", config.TargetPath, config.Apps)

		syncer := internal.NewSyncer(config.TargetPath, preferDirection)

		for _, app := range config.Apps {
			fmt.Printf("--- syncing app: %s ---\n", app)
			found, exists := configurations[app]
			if !exists {
				fmt.Printf("config for app '%s' does not exist, skipping\n", app)
				continue
			}

			switch action {
			case actionTypeSync:
				if err := syncer.SyncApp(found); err != nil {
					fmt.Printf("err syncing app: %s\n", err)
				}

			case actionTypeInitFromLocal:
				if err := syncer.SyncAppBToA(found); err != nil {
					fmt.Printf("err syncing app: %s\n", err)
				}

			case actionTypeInitFromTarget:
				if err := syncer.SyncAppAToB(found); err != nil {
					fmt.Printf("err syncing app: %s\n", err)
				}
			}

		}

	default:
		fmt.Println("action not found")
		flag.Usage()
		os.Exit(1)
	}
}
