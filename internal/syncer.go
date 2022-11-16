package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dvcrn/unisync/internal/engine"
)

type PreferMode string

const (
	PreferModeTarget = "target"
	PreferModeLocal  = "local"
)

var globalIgnore = []string{
	"Name .DS_Store",
}

func ToAbsolutePath(path string) (string, error) {
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
	PreferMode PreferMode
}

func NewSyncer(targetPath string, targetMode PreferMode) *Syncer {
	return &Syncer{
		TargetPath: targetPath,
		PreferMode: targetMode,
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

		absPathA, err := ToAbsolutePath(targetPath)
		if err != nil {
			return err
		}

		absPathB, err := ToAbsolutePath(fileConfig.BasePath)
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
			prefer := absPathA
			if s.PreferMode == PreferModeLocal {
				prefer = absPathB
			}

			if err = unison.Sync(absPathA, absPathB, prefer, fileConfig.IncludedFiles, ignoredFiles); err != nil {
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
