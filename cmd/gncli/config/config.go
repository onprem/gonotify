package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config represnts cli config
type Config struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Token    string `json:"token"`
	BaseURL  string `json:"baseURL"`
	Path     string `json:"-"`
}

// LoadDefault loads the default config
func (c *Config) LoadDefault() {
	c.BaseURL = "https://gonotify.xyz"
}

// Save saves the config at given path
func (c *Config) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Path, data, 0644)
}

// Load loads the config from given path
func (c *Config) Load() error {
	file, err := os.Open(c.Path)
	if err != nil {
		c.LoadDefault()

		err := c.Save()
		if err != nil {
			return err
		}

		return nil
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}
