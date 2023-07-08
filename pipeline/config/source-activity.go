package config

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"regexp"
)

type Mode string

const (
	ModeTag Mode = "tag"
)

type TagValueId string

const (
	TagValueReady   TagValueId = "ready"
	TagValueWorking TagValueId = "working"
	TagValueDone    TagValueId = "done"
)

type Tag struct {
	Name   string     `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Values []TagValue `mapstructure:"values,omitempty" yaml:"values,omitempty" json:"values,omitempty"`
}

type TagValue struct {
	Id    TagValueId `mapstructure:"id,omitempty" yaml:"id,omitempty" json:"id,omitempty"`
	Value string     `mapstructure:"value,omitempty" yaml:"value,omitempty" json:"value,omitempty"`
}

type Path struct {
	Container   string         `mapstructure:"container,omitempty" yaml:"container,omitempty" json:"container,omitempty"`
	NamePattern string         `mapstructure:"pattern,omitempty" yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Id          string         `mapstructure:"id,omitempty" yaml:"id,omitempty" json:"id,omitempty"`
	Regexp      *regexp.Regexp `mapstructure:"-" yaml:"-" json:"-"`
}

type SourceActivity struct {
	Activity
	StorageName  string `mapstructure:"storage-name,omitempty" yaml:"storage-name,omitempty" json:"storage-name,omitempty"`
	Mode         Mode   `mapstructure:"mode,omitempty" yaml:"mode,omitempty" json:"mode,omitempty"`
	TagName      string `mapstructure:"tag-name,omitempty" yaml:"tag-name,omitempty" json:"tag-name,omitempty"`
	Tag          Tag    `mapstructure:"tag,omitempty" yaml:"tag,omitempty" json:"tag,omitempty"`
	Paths        []Path `mapstructure:"paths,omitempty" yaml:"paths,omitempty" json:"paths,omitempty"`
	TickInterval string `mapstructure:"tick-interval" yaml:"tick-interval" json:"tick-interval"`
	DownloadPath string `mapstructure:"download-path" yaml:"download-path" json:"download-path"`
}

func (c *SourceActivity) WithName(n string) *SourceActivity {
	c.Nm = n
	return c
}

func (c *SourceActivity) WithDescription(n string) *SourceActivity {
	c.Cm = n
	return c
}

func NewSourceActivity() *SourceActivity {
	s := SourceActivity{}
	s.Tp = SourceActivityType
	return &s
}

func NewSourceActivityFromJSON(message json.RawMessage) (Configurable, error) {
	i := NewSourceActivity()
	err := json.Unmarshal(message, i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func NewSourceActivityFromYAML(mp interface{}) (Configurable, error) {
	sa := NewSourceActivity()
	err := mapstructure.Decode(mp, sa)
	if err != nil {
		return nil, err
	}

	return sa, nil
}
