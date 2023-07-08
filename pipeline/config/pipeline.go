package config

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Pipeline struct {
	Id            string            `mapstructure:"id,omitempty" yaml:"id,omitempty" json:"id,omitempty"`
	Description   string            `mapstructure:"description,omitempty" yaml:"description,omitempty" json:"description,omitempty"`
	Activities    []Configurable    `json:"-" yaml:"activities"`
	RawActivities []json.RawMessage `json:"activities" yaml:"-"`
}

func NewPipelineFromJSON(data []byte) (Pipeline, error) {
	o := Pipeline{}
	err := json.Unmarshal(data, &o)

	return o, err
}

func NewPipelineFromYAML(data []byte) (Pipeline, error) {
	o := Pipeline{}
	err := yaml.Unmarshal(data, &o)

	return o, err
}

func (o *Pipeline) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Pipeline) ToYAML() ([]byte, error) {
	return yaml.Marshal(o)
}

func (o *Pipeline) FindActivityByName(n string) Configurable {
	for _, a := range o.Activities {
		if a.Name() == n {
			return a
		}
	}

	return nil
}

func (o *Pipeline) AddActivity(a Configurable) error {

	if o.FindActivityByName(a.Name()) != nil {
		return fmt.Errorf("activity with the same id already present (id: %s)", a.Name())
	}

	o.Activities = append(o.Activities, a)
	return nil
}

func (o *Pipeline) UnmarshalJSON(b []byte) error {

	// Clear the state....
	o.Activities = nil

	type pipeline Pipeline
	err := json.Unmarshal(b, (*pipeline)(o))
	if err != nil {
		return err
	}

	for _, raw := range o.RawActivities {
		var v Activity
		err = json.Unmarshal(raw, &v)
		if err != nil {
			return err
		}

		i, err := NewActivityFromJSON(v.Type(), raw)
		if err != nil {
			return err
		}

		o.AddActivity(i)
	}
	return nil
}

func (o *Pipeline) MarshalJSON() ([]byte, error) {

	// Clear the state....
	o.RawActivities = nil

	type pipeline Pipeline
	if o.Activities != nil {
		for _, v := range o.Activities {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			o.RawActivities = append(o.RawActivities, b)
		}
	}
	return json.Marshal((*pipeline)(o))
}

func (o *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type pipeline Pipeline

	var m struct {
		Id          string        `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
		Description string        `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
		Activities  []interface{} `json:"activities" yaml:"activities"`
	}
	m.Activities = make([]interface{}, 0)
	err := unmarshal(&m)
	if err != nil {
		return err
	}

	o.Id = m.Id
	o.Description = m.Description

	for _, a := range m.Activities {
		var wa struct {
			Activity Activity
		}
		err := mapstructure.Decode(a, &wa)
		if err != nil {
			return err
		}

		i, err := NewActivityFromYAML(wa.Activity.Type(), a)
		if err != nil {
			return err
		}

		o.Activities = append(o.Activities, i)
	}

	return nil
}
