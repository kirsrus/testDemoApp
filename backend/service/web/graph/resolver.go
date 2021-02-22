package graph

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver создаёмся через конструкто NewResolver
type Resolver struct {
	log *logrus.Entry
}

// ConfigResolver конфигурация структуры Resolver
type ConfigResolver struct {
	Log *logrus.Logger
}

// NewResolver конструктор структуры Resolver
func NewResolver(config ConfigResolver) (*Resolver, error) {
	// region Проверка входящих данных

	log := logrus.New()
	log.Out = ioutil.Discard
	if config.Log != nil {
		log = config.Log
	}

	// endregion
	// region Настройка Resolver
	resolver := Resolver{
		log: log.WithFields(map[string]interface{}{
			"module": "graphql",
			"scope":  "service",
		}),
	}

	// endregion

	return &resolver, nil
}
