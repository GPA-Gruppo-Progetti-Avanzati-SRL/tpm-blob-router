package config

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"os"
)

const pipelineConfigFilePattern = "^pl[a-z0-9_-]*\\.yml"

type PipelineLoader struct {
	Typ        Type   `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	FolderPath string `yaml:"folder-path,omitempty" mapstructure:"folder-path,omitempty" json:"folder-path,omitempty"`
}

func (pll *PipelineLoader) Load() ([]Pipeline, error) {

	const semLogContext = "blob-router-loader::load"
	log.Info().Interface("type", pll.Typ).Str("folder", pll.FolderPath).Msg(semLogContext)

	files, err := util.FindFiles(pll.FolderPath, util.WithFindFileType(util.FileTypeFile), util.WithFindOptionIncludeList([]string{pipelineConfigFilePattern}))
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	log.Info().Int("num-files", len(files)).Msg(semLogContext + " found pipeline files")

	var pls []Pipeline
	for _, f := range files {
		log.Info().Str("filename", f).Msg(semLogContext + " loading pipeline")

		b, err := os.ReadFile(f)
		if err != nil {
			log.Error().Err(err).Str("filename", f).Msg(semLogContext)
			return nil, err
		}

		pl, err := NewPipelineFromYAML(b)
		if err != nil {
			log.Error().Err(err).Str("filename", f).Msg(semLogContext)
			return nil, err
		}

		log.Info().Str("id", pl.Id).Int("num-activities", len(pl.Activities)).Msg(semLogContext + " loaded pipeline")
		for i, a := range pl.Activities {
			log.Info().Int("act-ndx", i).Str("id", pl.Id).Interface("type", a.Type()).Str("name", a.Name()).Msg(semLogContext + " loaded pipeline")
		}
		pls = append(pls, pl)
	}

	return pls, nil
}
