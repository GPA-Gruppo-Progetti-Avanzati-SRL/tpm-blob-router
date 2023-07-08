package pipeline

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/pipeline/config"
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/pipeline/executable"
)

type Pipeline struct {
	Cfg         *config.Pipeline
	Executables map[string]executable.Executable
}
