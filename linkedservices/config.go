package linkedservices

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-az-common/cosmosdb/coslks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-az-common/storage/azstoragecfg"
	"gitlab.alm.poste.it/go/configuration"
)

type Config struct {
	CosmosDb []coslks.Config       `mapstructure:"cosmos-db,omitempty"  json:"cosmos-db,omitempty" yaml:"cosmos-db,omitempty"`
	Storage  []azstoragecfg.Config `mapstructure:"blob-storage,omitempty" json:"blob-storage,omitempty" yaml:"blob-storage,omitempty"`
}

func (c *Config) PostProcess() error {

	var err error

	if len(c.CosmosDb) > 0 {
		for i := range c.CosmosDb {
			err = c.CosmosDb[i].PostProcess()
			if err != nil {
				return err
			}
		}
	}

	if len(c.Storage) > 0 {
		for i := range c.Storage {
			err = c.Storage[i].PostProcess()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetDefaults() []configuration.VarDefinition {
	vd := make([]configuration.VarDefinition, 0)
	return vd
}
