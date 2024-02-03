package config

import (
	"fmt"
    "io/ioutil"
    "path/filepath"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GeminiFetchAddress	string	`yaml:"GeminiFetchAddress"`
	CertFilename    	string	`yaml:"CertFilename"`
	AutoFetch			bool	`yaml:"AutoFetch"`
	Timezone			string	`yaml:"Timezone"`
	// Add other configuration fields as needed
}

func ExportConfig(username string) (*Config, error) {
	// Construct the path to the YAML file based on the username
    yamlFile := filepath.Join("../users", username, "config.yaml")
    // Read the contents of the YAML file
    yamlContent, err := ioutil.ReadFile(yamlFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    // Unmarshal the YAML content into the Config struct
    var config Config
    if err := yaml.UnmarshalStrict(yamlContent, &config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
    }
    return &config, nil
}
