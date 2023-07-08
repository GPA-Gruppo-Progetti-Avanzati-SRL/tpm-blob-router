package config_test

import (
	"GPA-Gruppo-Progetti-Avanzati-SRL/tpm-blob-router/pipeline/config"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

const (
	EchoActivityName          = "echo-activity"
	OrchestrationYAMLFileName = "tpm-blob-router-pipeline.yml"
)

var cfgOrc config.Pipeline

func SetUpPipeline(t *testing.T) {
	ea := config.NewEchoActivity().WithName(EchoActivityName).WithDescription("test echo activity").WithMessage("hello echo activity")

	cfgOrc = config.Pipeline{
		Id:          "sample-pipeline",
		Description: "sample pipeline",
		Activities: []config.Configurable{
			ea,
		},
	}
}

func TestConfig(t *testing.T) {

	SetUpPipeline(t)

	var pln config.Pipeline

	t.Log("JSON SerDe --------------------------")
	b, err := cfgOrc.ToJSON()
	require.NoError(t, err)
	t.Log(string(b))

	// Deserialization
	pln, err = config.NewPipelineFromJSON(b)
	require.NoError(t, err)

	b, err = pln.ToJSON()
	require.NoError(t, err)
	t.Log(string(b))

	t.Log("YAML SerDe --------------------------")
	b, err = cfgOrc.ToYAML()
	require.NoError(t, err)

	err = os.WriteFile(OrchestrationYAMLFileName, b, fs.ModePerm)
	require.NoError(t, err)

	// Should remove... at the moment is good this way....
	// defer os.Remove(OrchestrationYAMLFileName)

	pln, err = config.NewPipelineFromYAML(b)
	require.NoError(t, err)

	b, err = pln.ToYAML()
	require.NoError(t, err)
	t.Log(string(b))
}

var serde = []byte(`
activities: 
  - activity:
      name: start-name
      type: start-activity
    property: a-start-property
  - activity:
      name: echo-name
      type: echo-activity
    message: a-message
  - activity:
      name: end-name
      type: end-activity
`)

func TestConfigSerde(t *testing.T) {
	deserOrch, err := config.NewPipelineFromYAML(serde)
	require.NoError(t, err)

	b, err := deserOrch.ToYAML()
	require.NoError(t, err)
	t.Log(string(b))
}
