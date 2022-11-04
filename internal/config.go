package internal

type Config struct {
	Apps       []string `yaml:"apps"`
	TargetPath string   `yaml:"targetPath"`
}
