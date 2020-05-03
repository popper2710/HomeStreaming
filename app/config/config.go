package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Secret struct {
	Database struct {
		Image    string `yaml:"image"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}

	Redis struct {
		Size     int    `yaml:"size"`
		Network  string `yaml:"network"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		KeyPair  string `yaml:"keyPair"`
	}
}

func (s *Secret) Init() {
	buf, err := ioutil.ReadFile("./config/secret.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(buf, &s)
	if err != nil {
		panic(err)
	}
}
