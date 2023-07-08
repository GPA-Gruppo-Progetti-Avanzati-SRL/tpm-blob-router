package linkedservices

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-az-common/cosmosdb/coslks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-az-common/storage/azbloblks"
	"github.com/rs/zerolog/log"
)

type ServiceRegistry struct {
}

var registry ServiceRegistry

func InitRegistry(cfg *Config) error {

	registry = ServiceRegistry{}
	log.Info().Msg("initialize services registry")

	var err error

	_, err = coslks.Initialize(cfg.CosmosDb)
	if err != nil {
		return err
	}

	_, err = azbloblks.Initialize(cfg.Storage)
	if err != nil {
		return err
	}

	return nil
}
