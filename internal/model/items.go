package model

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type ItemYamlConfig struct {
	Sources map[string]struct {
		Items ItemMetas `yaml:"items"`
	} `yaml:"sources"`
}

type ItemMeta struct {
	Category    string `yaml:"category"`
	Unit        string `yaml:"unit"`
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
}

type ItemMetas map[string]ItemMeta

//const itemConfigPath = "item_meta.yaml"

func LoadItemYamlConfig(path string) (*ItemYamlConfig, error) {
	var itemConfig ItemYamlConfig

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &itemConfig)
	if err != nil {
		return nil, err
	}

	return &itemConfig, nil
}

func (c ItemYamlConfig) ConvertItemMetas() ItemMetas {
	itemMetas := make(ItemMetas)
	for _, sources := range c.Sources {
		for itemName, itemMeta := range sources.Items {
			itemMetas[itemName] = itemMeta
		}
	}
	return itemMetas
}

func LoadItemMetas(path string) (ItemMetas, error) {
	itemConfig, err := LoadItemYamlConfig(path)
	if err != nil {
		return nil, err
	}
	return itemConfig.ConvertItemMetas(), nil
}

func (m ItemMetas) GetItemMeta(itemName string) (ItemMeta, bool) {
	itemMeta, ok := m[itemName]
	return itemMeta, ok
}

func (m ItemMetas) GetItemNamesByCategory(category string) []string {
	var itemNames []string
	for itemName, itemMeta := range m {
		if itemMeta.Category == category {
			itemNames = append(itemNames, itemName)
		}
	}
	return itemNames
}
