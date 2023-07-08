package executable

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/pipeline/config"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/promutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"time"
)

type Executable interface {
	Execute() error
}

type Activity struct {
	Cfg config.Configurable
}

func (a *Activity) Name() string {
	return a.Cfg.Name()
}

func (a *Activity) Type() config.Type {
	return a.Cfg.Type()
}

func (a *Activity) IsDisabled() bool {
	return a.IsDisabled()
}

func (a *Activity) IsValid() bool {

	rc := true
	switch a.Cfg.Type() {
	case config.EchoActivityType:
	}

	return rc
}

func (a *Activity) MetricsGroup() (promutil.Group, bool, error) {
	mCfg := a.Cfg.MetricsConfig()

	var g promutil.Group
	var err error
	var ok bool
	if mCfg.IsEnabled() {
		g, err = promutil.GetGroup(mCfg.GId)
		if err == nil {
			ok = true
		}
	}

	return g, ok, err
}

func (a *Activity) SetMetrics(begin time.Time, lbls prometheus.Labels) error {

	const semLogContext = "executable::set-metrics"
	cfg := a.Cfg.MetricsConfig()
	if cfg.IsEnabled() {
		g, _, err := a.MetricsGroup()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return err
		}

		if cfg.IsCounterEnabled() {
			g.SetMetricValueById(cfg.CounterId, 1, lbls)
		}

		if cfg.IsHistogramEnabled() {
			g.SetMetricValueById(cfg.HistogramId, time.Since(begin).Seconds(), lbls)
		}
	}

	return nil
}
