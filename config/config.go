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
	LastFetched			string	`yaml:"LastFetched"`
	GMAP				bool	`yaml:"GMAP"`
	DisplayNumber		int		`yaml:"DisplayNumber"`
	HeaderText			string	`yaml:"HeaderText"`
	TagSync				bool	`yaml:"TagSync"`
	GMAPHost			string	`yaml:"GMAPHost"`
	GMAPPort			string	`yaml:"GMAPPort"`
	// Add other configuration fields as needed
}

type SiteConfig struct {
	ExposeSettings	bool `yaml:"ExposeSettings"`
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

func Set(username string, fieldName string, value interface{}) error {
	yamlFile := filepath.Join("../users", username, "config.yaml")
	yamlContent, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	var configMap map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &configMap); err != nil {
		return fmt.Errorf("failed to unmarshal config data: %w", err)
	}
	configMap[fieldName] = value
	updatedYAML, err := yaml.Marshal(&configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal config data: %w", err)
	}
	if err := ioutil.WriteFile(yamlFile, updatedYAML, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func ExportSiteConfig() (*SiteConfig, error) {
	// Construct the path to the YAML file based on the username
    yamlFile := filepath.Join("../site_config.yaml")
    // Read the contents of the YAML file
    yamlContent, err := ioutil.ReadFile(yamlFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    // Unmarshal the YAML content into the Config struct
    var config SiteConfig
    if err := yaml.UnmarshalStrict(yamlContent, &config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
    }
    return &config, nil
}
