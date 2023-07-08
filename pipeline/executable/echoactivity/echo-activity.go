package echoactivity

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/constants"
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/pipeline/config"
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/pipeline/executable"
	"github.com/rs/zerolog/log"
)

type EchoActivity struct {
	executable.Activity
}

func NewEchoActivity(item config.Configurable) (*EchoActivity, error) {
	ea := &EchoActivity{}
	ea.Cfg = item
	return ea, nil
}

func (a *EchoActivity) Execute() error {

	const semLogContext = "echo-activity::execute"
	log.Trace().Str(constants.SemLogActivity, a.Name()).Bool("enabled", !a.IsDisabled()).Str("type", "echo").Msg(semLogContext + " start")
	if !a.IsDisabled() {

	}

	_, ok := a.Cfg.(*config.EchoActivity)
	if !ok {
		log.Error().Msgf("this is weird %v is not (*config.EchoActivity)", a.Cfg)
	}

	log.Trace().Str(constants.SemLogActivity, a.Name()).Bool("enabled", !a.IsDisabled()).Str("type", "echo").Msg(semLogContext + " end")
	return nil
}
