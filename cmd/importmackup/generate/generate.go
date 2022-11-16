package generate

import (
	"fmt"
	"github.com/dvcrn/unisync/internal"
	"path"
	"strings"
)

func sanitizeName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "")
}

func GenerateConfig(appName string, files []string) (*internal.AppConfig, error) {
	appConfig := internal.AppConfig{
		Name:         appName,
		FriendlyName: sanitizeName(appName),
	}

	fileRoots := map[string][]string{}

	for _, file := range files {
		// check if absolute path first
		if file[0] != '/' {
			file = fmt.Sprintf("~/%s", file)
		}

		dir := path.Dir(file) + "/"
		if _, ok := fileRoots[dir]; !ok {
			fileRoots[dir] = []string{}
		}

		base := path.Base(file)
		fileRoots[dir] = append(fileRoots[dir], base)
	}

	for root, configs := range fileRoots {
		appConfig.Files = append(appConfig.Files, internal.FileConfig{
			BasePath:      root,
			IncludedFiles: configs,
		})
	}

	return &appConfig, nil
}
