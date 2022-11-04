package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dvcrn/uniconfig/internal"
	"gopkg.in/yaml.v3"
)

//go:embed apps/*
var content embed.FS

type actionType string

const (
	actionTypeUnknown        actionType = ""
	actionTypeSync           actionType = "sync"
	actionTypeInitFromTarget actionType = "initFromTarget"
	actionTypeInitFromApp    actionType = "initFromApp"
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
		fmt.Println("  init-from-target - run initial sync, targetPath -> appPath")
		fmt.Println("  init-from-app - run initial sync, appPath -> targetPath")
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
	case "init-from-app":
		action = actionTypeInitFromApp

	}
}
func main() {
	switch action {
	case actionTypeSync:
	case actionTypeInitFromApp:
	case actionTypeInitFromTarget:
		configurations := map[string]*internal.AppConfig{}

		entries, err := content.ReadDir("apps")
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			c, err := content.ReadFile(fmt.Sprintf("apps/%s", entry.Name()))
			if err != nil {
				log.Fatal(err)
			}

			var t *internal.AppConfig
			err = yaml.Unmarshal(c, &t)
			if err != nil {
				log.Fatal(err)
			}

			configurations[t.FriendlyName] = t
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		configFile, err := os.ReadFile(homeDir + "/.config/uniconfig/uniconfig.yaml")
		if err != nil {
			log.Fatal(err)
		}

		var config *internal.Config
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			log.Fatalf("err trying to parse yaml config file: %s\n", err)
		}

		fmt.Println("starting uniconfig...")
		fmt.Printf("targetPath: %s, apps: %s\n", config.TargetPath, config.Apps)

		syncer := internal.NewSyncer(config.TargetPath)

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

			case actionTypeInitFromApp:
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
