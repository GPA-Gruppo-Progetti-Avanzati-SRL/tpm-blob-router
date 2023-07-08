package main

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/linkedservices"
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/filetracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/logzerotracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-kafka-har/kafkahartracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv/resource/health"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-srv/httpsrv/resource/metrics"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-kafka-common/kafkalks"
	"github.com/rs/zerolog/log"
	"gitlab.alm.poste.it/go/observability/tracing"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

//go:embed sha.txt
var sha string

//go:embed VERSION
var version string

// appLogo contains the ASCII splash screen
//
//go:embed app-logo.txt
var appLogo []byte

func main() {

	const semLogContext = "tpm-blob-router::main"

	fmt.Println(string(appLogo))
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Sha: %s\n", sha)

	appCfg, err := ReadConfig()
	if nil != err {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	log.Info().Interface("config", appCfg).Msg("configuration loaded")

	jc, err := InitGlobalTracer()
	if nil != err {
		log.Fatal().Err(err).Msg(semLogContext)
	}
	defer jc.Close()

	err = linkedservices.InitRegistry(appCfg.Services)
	if nil != err {
		log.Fatal().Err(err).Msg(semLogContext + " linked services initialization error")
	}

	// Har Tracing is not enabled in blob processor
	/*
		hc, err := InitHarTracing()
		if nil != err {
			log.Fatal().Err(err).Msg(semLogContext)
		}
		if hc != nil {
			defer hc.Close()
		}
	*/

	if appCfg.MwRegistry != nil {
		if err := mwregistry.InitializeHandlerRegistry(appCfg.MwRegistry.ToHandlerCatalogConfig(), appCfg.Http.MwUse); err != nil {
			log.Fatal().Err(err).Msg(semLogContext)
		}
	}

	// shutdownChannel := make(chan os.Signal, 1)
	// signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info().Msg("Enabling SIGINT e SIGTERM")
	shutdownChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		shutdownChannel <- fmt.Errorf("signal received: %v", <-c)
	}()

	s, err := httpsrv.NewServer(appCfg.Http /* , httpsrv.WithListenPort(9090), httpsrv.WithDocumentRoot("/www", "/tmp", false) */)
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	if err := s.Start(); err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}
	defer s.Stop()

	var wg sync.WaitGroup

	for !s.IsReady() {
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	sig := <-shutdownChannel
	log.Debug().Interface("signal", sig).Msg("got termination signal")

	wg.Wait()
	log.Info().Msg("terminated...")
}

func InitHarTracing() (io.Closer, error) {

	const semLogContext = "main::init-har-tracing"
	const semLogLabelTracerType = "tracer-type"
	var trc hartracing.Tracer
	var closer io.Closer
	var err error

	trcType := os.Getenv(hartracing.HARTracerTypeEnvName)
	switch strings.ToLower(trcType) {
	case filetracer.HarFileTracerType:
		trc, closer, err = filetracer.NewTracer()
		if err != nil {
			return nil, err
		}

		log.Info().Str(semLogLabelTracerType, trcType).Msg(semLogContext + " file har tracer initialized")

	case logzerotracer.HarLogZeroTracerType:
		trc, closer, err = logzerotracer.NewTracer()
		if err != nil {
			return nil, err
		}
		log.Info().Str(semLogLabelTracerType, trcType).Msg(semLogContext + " logzero har tracer initialized")

	case kafkahartracer.HarKafkaTracerType:
		brokerName := os.Getenv(kafkahartracer.BrokerNameEnvVar)
		if brokerName == "" {
			err := fmt.Errorf("broker name environment variable %s not set", kafkahartracer.BrokerNameEnvVar)
			log.Error().Err(err).Str(semLogLabelTracerType, trcType).Msgf(semLogContext)
			return nil, err
		}

		topic := os.Getenv(kafkahartracer.TopicNameEnvVar)
		if topic == "" {
			err := fmt.Errorf("topic name environment variable %s not set", kafkahartracer.TopicNameEnvVar)
			log.Error().Err(err).Str(semLogLabelTracerType, trcType).Msgf(semLogContext)
			return nil, err
		}

		lks, err := kafkalks.GetKafkaLinkedService(brokerName)
		if err != nil {
			return nil, err
		}

		trc, closer, err = kafkahartracer.NewTracer(kafkahartracer.WithKafkaLinkedService(lks), kafkahartracer.WithTopic(topic))
		if err != nil {
			return nil, err
		}

		log.Info().Str(semLogLabelTracerType, trcType).Str("broker-name", brokerName).Str("topic", topic).Msg(semLogContext + " kafka har tracer initialized")
	default:
		log.Info().Str(semLogLabelTracerType, trcType).Msgf(semLogContext+" env var %s not set or unrecognized tracer type", "HAR_TRACER_TYPE")
	}

	if trc != nil {
		hartracing.SetGlobalTracer(trc)
	}

	return closer, nil
}

func InitGlobalTracer() (*tracing.Tracer, error) {
	tracer, err := tracing.NewTracer()
	if err != nil {
		return nil, err
	}

	return tracer, err
}
