package config

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"sync"

	results "github.com/jxm35/go-results"
	yaml "gopkg.in/yaml.v3"
)

var (
	config               map[string]any
	configOnce           sync.Once
	ErrConfigNotProvided = errors.New("config not provided")
)

//go:embed files/*
var configFiles embed.FS

func loadConfig(env string) error {
	config = make(map[string]any)
	file, err := configFiles.Open(fmt.Sprintf("files/%s.yaml", env))
	if err != nil {
		return err
	}
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(&config)
	return err
}

func GetConfigVal[T any](name string) results.Option[T] {
	configOnce.Do(func() {
		err := loadConfig(os.Getenv("env"))
		if err != nil {
			panic(err)
		}
	})
	val, ok := config[name]
	if !ok {
		return results.None[T]()
	}
	typedVal, ok := val.(T)
	if !ok {
		return results.None[T]()
	}
	return results.Some(typedVal)
}
