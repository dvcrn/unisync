package generate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateConfig_SingleFile(t *testing.T) {
	configFiles := []string{"Library/Preferences/com.trankynam.XtraFinder.plist"}

	appConfig, err := GenerateConfig("xtrafinder", "XtraFinder", configFiles)

	assert.NoError(t, err)
	assert.Equal(t, "XtraFinder", appConfig.Name)
	assert.Equal(t, "xtrafinder", appConfig.FriendlyName)
	assert.Len(t, appConfig.Files, 1)
	assert.Equal(t, appConfig.Files[0].BasePath, "~/Library/Preferences/")
	assert.Len(t, appConfig.Files[0].IncludedFiles, 1)
	assert.Equal(t, appConfig.Files[0].IncludedFiles[0], "com.trankynam.XtraFinder.plist")
}

func TestGenerateConfig_FriendlyName(t *testing.T) {
	// table test
	testCases := []struct {
		name     string
		expected string
	}{
		{"XtraFinder", "xtrafinder"},
		{"1Password 4", "1password4"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			appConfig, err := GenerateConfig(testCase.name, "", []string{})
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, appConfig.FriendlyName)
		})
	}
}

func TestGenerateConfig_MultipleFilesSameRoot(t *testing.T) {
	configFiles := []string{"Library/Application Support/xbar/xbar.config.json", "Library/Application Support/xbar/plugins"}

	appConfig, err := GenerateConfig("xtrafinder", "XtraFinder", configFiles)

	assert.NoError(t, err)
	assert.Equal(t, "XtraFinder", appConfig.Name)
	assert.Equal(t, "xtrafinder", appConfig.FriendlyName)
	assert.Len(t, appConfig.Files, 1)
	assert.Equal(t, appConfig.Files[0].BasePath, "~/Library/Application Support/xbar/")
	assert.Len(t, appConfig.Files[0].IncludedFiles, 2)
	assert.Contains(t, appConfig.Files[0].IncludedFiles, "xbar.config.json")
	assert.Contains(t, appConfig.Files[0].IncludedFiles, "plugins")
}

func TestGenerateConfig_MultipleFilesMultiRoot(t *testing.T) {
	configFiles := []string{"Library/Application Support/xbar/xbar.config.json", "Library/Application Support/xbar/plugins", "Library/Preferences/com.trankynam.XtraFinder.plist"}

	appConfig, err := GenerateConfig("xtrafinder", "XtraFinder", configFiles)

	assert.NoError(t, err)
	assert.Equal(t, "XtraFinder", appConfig.Name)
	assert.Equal(t, "xtrafinder", appConfig.FriendlyName)
	assert.Len(t, appConfig.Files, 2)
	assert.Equal(t, appConfig.Files[0].BasePath, "~/Library/Application Support/xbar/")

	assert.Len(t, appConfig.Files[0].IncludedFiles, 2)
	assert.Contains(t, appConfig.Files[0].IncludedFiles, "xbar.config.json")
	assert.Contains(t, appConfig.Files[0].IncludedFiles, "plugins")

	assert.Len(t, appConfig.Files[1].IncludedFiles, 1)
	assert.Contains(t, appConfig.Files[1].IncludedFiles, "com.trankynam.XtraFinder.plist")
}
