package main

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/linkedservices"
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/pipeline/config"
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwhartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwerror"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwmetrics"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwtracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	"github.com/rs/zerolog/log"
	"gitlab.alm.poste.it/go/configuration"
	"os"
	"strings"
)

type MiddlewareConfig struct {
	MwErrorConfig      *mwerror.ErrorHandlerConfig             `yaml:"mw-error,omitempty" mapstructure:"mw-error,omitempty" json:"mw-error,omitempty"`
	MwTracingConfig    *mwtracing.TracingHandlerConfig         `yaml:"mw-tracing,omitempty" mapstructure:"mw-tracing,omitempty" json:"mw-tracing,omitempty"`
	MwHarTracingConfig *mwhartracing.HarTracingHandlerConfig   `yaml:"mw-har-tracing,omitempty" mapstructure:"mw-har-tracing,omitempty" json:"mw-har-tracing,omitempty"`
	MwMetricsConfig    *mwmetrics.PromHttpMetricsHandlerConfig `yaml:"mw-metrics,omitempty" mapstructure:"mw-metrics,omitempty" json:"mw-metrics,omitempty"`
}

func (mwCfg *MiddlewareConfig) ToHandlerCatalogConfig() mwregistry.HandlerCatalogConfig {

	r := make(mwregistry.HandlerCatalogConfig)
	if mwCfg.MwHarTracingConfig != nil {
		r[mwhartracing.HarTracingHandlerId] = mwCfg.MwHarTracingConfig
	}

	if mwCfg.MwErrorConfig != nil {
		r[mwerror.ErrorHandlerId] = mwCfg.MwErrorConfig
	}

	if mwCfg.MwTracingConfig != nil {
		r[mwtracing.TracingHandlerId] = mwCfg.MwTracingConfig
	}

	if mwCfg.MwMetricsConfig != nil {
		r[mwmetrics.MetricsHandlerId] = mwCfg.MwMetricsConfig
	}

	return r
}

type AppConfig struct {
	Hostname         string                 `yaml:"host-name" mapstructure:"host-name" json:"host-name"`
	Http             httpsrv.Config         `yaml:"http" mapstructure:"http" json:"http"`
	MwRegistry       *MiddlewareConfig      `yaml:"mw-handler-registry" mapstructure:"mw-handler-registry" json:"mw-handler-registry"`
	Services         *linkedservices.Config `yaml:"linked-services" mapstructure:"linked-services" json:"linked-services"`
	BlobRouterLoader *config.PipelineLoader `yaml:"blob-router-loader" mapstructure:"blob-router-loader" json:"blob-router-loader"`
	BlobRouter       []config.Pipeline      `yaml:"-" mapstructure:"-" json:"-"`
}

// Default config file.
//
//go:embed config.yml
var projectConfigFile []byte

const ConfigFileEnvVar = "TPM_BLOB_ROUTER_CFG_FILE_PATH"
const ConfigurationName = "tpm_blob_router"

func ReadConfig() (*AppConfig, error) {

	configPath := os.Getenv(ConfigFileEnvVar)
	var cfgFileReader *strings.Reader
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			log.Info().Str("cfg-file-name", configPath).Msg("reading config")
			cfgContent, rerr := util.ReadFileAndResolveEnvVars(configPath)
			if rerr != nil {
				return nil, err
			} else {
				cfgFileReader = strings.NewReader(string(cfgContent))
			}

		} else {
			return nil, fmt.Errorf("the %s env variable has been set but no file cannot be found at %s", ConfigFileEnvVar, configPath)
		}
	} else {
		log.Warn().Msgf("The config path variable %s has not been set. Reverting to bundled configuration", ConfigFileEnvVar)
		cfgFileReader = strings.NewReader(util.ResolveConfigValueToString(string(projectConfigFile)))

		// return nil, fmt.Errorf("the config path variable %s has not been set; please set", ConfigFileEnvVar)
	}

	appCfg := &AppConfig{}
	_, err := configuration.NewConfiguration(
		configuration.WithType("yaml"),
		configuration.WithName(ConfigurationName),
		configuration.WithReader(cfgFileReader),
		configuration.WithData(appCfg))

	if err != nil {
		return nil, err
	}

	if appCfg.BlobRouterLoader != nil {
		appCfg.BlobRouter, err = appCfg.BlobRouterLoader.Load()
		if err != nil {
			return nil, err
		}
	}
	return appCfg, nil
}

func (m *AppConfig) GetDefaults() []configuration.VarDefinition {
	vd := make([]configuration.VarDefinition, 0, 20)
	vd = append(vd, GetHttpSrvConfigDefaults()...)
	vd = append(vd, GetMiddlewareConfigDefaults("config.mw-handler-registry")...)
	vd = append(vd, linkedservices.GetDefaults()...)
	return vd
}

func GetHttpSrvConfigDefaults() []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{"config.http.bind-address", httpsrv.DefaultBindAddress, "host reference"},
		{"config.http.server-context.path", httpsrv.DefaultContextPath, "context-path"},
		{"config.http.port", httpsrv.DefaultListenPort, "port"},
		{"config.http.shutdown-timeout", httpsrv.DefaultShutdownTimeout, "shutdown timeout"},
		{"config.http.server-mode", httpsrv.DefaultServerMode, "modalita' di lavoro server gin"},
	}
}

func GetMiddlewareConfigDefaults(contextPath string) []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{strings.Join([]string{contextPath, mwerror.ErrorHandlerId, "with-cause"}, "."), mwerror.ErrorHandlerDefaultWithCause, "error is in clear"},
		{strings.Join([]string{contextPath, mwerror.ErrorHandlerId, "alphabet"}, "."), mwerror.ErrorHandlerDefaultAlphabet, "alphabet"},
		{strings.Join([]string{contextPath, mwerror.ErrorHandlerId, "span-tag"}, "."), mwerror.ErrorHandlerDefaultSpanTag, "span-tag"},
		{strings.Join([]string{contextPath, mwerror.ErrorHandlerId, "header"}, "."), mwerror.ErrorHandlerDefaultHeader, "header"},
	}
}

func (m *AppConfig) PostProcess() error {
	return nil
}
