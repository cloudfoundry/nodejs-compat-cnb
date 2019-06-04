package resources

type OverrideDependency struct {
	CfStacks []string `yaml:"cf_stacks"`
	File     string   `yaml:"file"`
	Name     string   `yaml:"name"`
	Sha256   string   `yaml:"sha256"`
	URI      string   `yaml:"uri"`
	Version  string   `yaml:"version"`
}

type Override struct {
	Nodejs struct {
		Dependencies []OverrideDependency `yaml:"dependencies"`
	} `yaml:"nodejs"`
}

