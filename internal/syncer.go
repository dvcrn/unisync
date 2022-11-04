package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dvcrn/uniconfig/internal/engine"
)

var globalIgnore = []string{
	"Name .DS_Store",
}

func toAbsolutePath(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	if strings.Contains(path, "~") {
		return filepath.Abs(strings.Replace(path, "~", homeDir, -1))
	}

	return filepath.Abs(path)
}

type Syncer struct {
	TargetPath string
}

func NewSyncer(targetPath string) *Syncer {
	return &Syncer{
		TargetPath: targetPath,
	}
}

type syncMode int8

const (
	syncModeNormal syncMode = iota
	syncModeAToB
	syncModeBToA
)

func (s *Syncer) sync(appConfig *AppConfig, syncMode syncMode) error {
	for _, fileConfig := range appConfig.Files {
		targetPath := filepath.Join(s.TargetPath, appConfig.FriendlyName)

		absPathA, err := toAbsolutePath(targetPath)
		if err != nil {
			return err
		}

		absPathB, err := toAbsolutePath(fileConfig.BasePath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absPathA); os.IsNotExist(err) {
			if err = os.MkdirAll(absPathA, 0700); err != nil {
				return err
			}
		}

		ignoredFiles := append([]string{}, fileConfig.IgnoredFiles...)
		ignoredFiles = append(ignoredFiles, globalIgnore...)

		unison := engine.NewUnison()
		switch syncMode {
		case syncModeNormal:
			if err = unison.Sync(absPathA, absPathB, fileConfig.IncludedFiles, ignoredFiles); err != nil {
				return err
			}
		case syncModeAToB:
			if err = unison.SyncAToB(absPathA, absPathB, fileConfig.IncludedFiles, ignoredFiles); err != nil {
				return err
			}
		case syncModeBToA:
			if err = unison.SyncBToA(absPathA, absPathB, fileConfig.IncludedFiles, ignoredFiles); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Syncer) SyncApp(appConfig *AppConfig) error {
	return s.sync(appConfig, syncModeNormal)
}

func (s *Syncer) SyncAppAToB(appConfig *AppConfig) error {
	return s.sync(appConfig, syncModeAToB)
}

func (s *Syncer) SyncAppBToA(appConfig *AppConfig) error {
	return s.sync(appConfig, syncModeBToA)
}
