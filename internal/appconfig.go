package internal

type FileConfig struct {
	BasePath      string   `yaml:"basePath"`
	IncludedFiles []string `yaml:"includedFiles"`
	IgnoredFiles  []string `yaml:"ignoredFiles"`
}

type AppConfig struct {
	Name         string       `yaml:"name"`
	FriendlyName string       `yaml:"friendlyName"`
	Files        []FileConfig `yaml:"files"`
}
