package config

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/promutil"
)

type Type string

const (
	SourceActivityType  Type = "source-activity"
	EchoActivityType    Type = "echo-activity"
	CopyToActivityType  Type = "copy-to-activity"
	DeleteActivityType  Type = "delete-activity"
	SetTagsActivityType Type = "set-tags-activity"
)

type ActivityTypeRegistryEntry struct {
	Tp                 Type
	UnmarshallFromJSON func(raw json.RawMessage) (Configurable, error)
	UnmarshalFromYAML  func(mp interface{}) (Configurable, error)
}

type Configurable interface {
	Name() string
	Type() Type
	Disabled() bool
	Description() string
	MetricsConfig() promutil.MetricsConfigReference
}

var activityTypeRegistry = map[Type]ActivityTypeRegistryEntry{
	EchoActivityType: {Tp: EchoActivityType, UnmarshallFromJSON: NewEchoActivityFromJSON, UnmarshalFromYAML: NewEchoActivityFromYAML},
}

func NewActivityFromJSON(t Type, message json.RawMessage) (Configurable, error) {

	if e, ok := activityTypeRegistry[t]; ok {
		c, err := e.UnmarshallFromJSON(message)
		return c, err
	}

	return nil, fmt.Errorf("unknown activity type %s", t)
}

func NewActivityFromYAML(t Type, m interface{}) (Configurable, error) {

	if e, ok := activityTypeRegistry[t]; ok {
		c, err := e.UnmarshalFromYAML(m)
		return c, err
	}

	return nil, fmt.Errorf("unknown activity type %s", t)
}

const (
	DefaultMetricsGroupId = "activity"
	DefaultCounterId      = "activity-counter"
	DefaultHistogramId    = "activity-duration"
)

var DefaultMetricsCfg = promutil.MetricsConfigReference{
	GId:         DefaultMetricsGroupId,
	CounterId:   DefaultCounterId,
	HistogramId: DefaultHistogramId,
}

type Activity struct {
	Nm         string                          `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Tp         Type                            `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Cm         string                          `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Dis        bool                            `yaml:"disabled,omitempty" mapstructure:"disabled,omitempty" json:"disabled,omitempty"`
	MetricsCfg promutil.MetricsConfigReference `yaml:"ref-metrics,omitempty" mapstructure:"ref-metrics,omitempty" json:"ref-metrics,omitempty"`
}

func (c *Activity) WithName(n string) *Activity {
	c.Nm = n
	return c
}

func (c *Activity) WithDescription(n string) *Activity {
	c.Cm = n
	return c
}

func (c *Activity) Name() string {
	return c.Nm
}

func (c *Activity) Type() Type {
	return c.Tp
}

func (c *Activity) Description() string {
	return c.Cm
}

func (c *Activity) Disabled() bool {
	return c.Dis
}

func (c *Activity) MetricsConfig() promutil.MetricsConfigReference {
	r := promutil.CoalesceMetricsConfig(c.MetricsCfg, DefaultMetricsCfg)
	return r
}
