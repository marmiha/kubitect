package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type WorkerDefault struct {
	CPU          *CpuSize            `yaml:"cpu,omitempty"`
	Labels       *map[LabelKey]Label `yaml:"labels,omitempty"`
	MainDiskSize *MB                 `yaml:"mainDiskSize,omitempty"`
	RAM          *MB                 `yaml:"ram,omitempty"`
	Taints       *[]Taint            `yaml:"taints,omitempty"`
}

func (d WorkerDefault) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Labels),
		validation.Field(&d.Taints),
		validation.Field(&d.RAM),
		validation.Field(&d.MainDiskSize),
		validation.Field(&d.CPU),
	)
}

type LabelKey string // TODO: Check if correct type
type Label string    // TODO: Check if correct type
type Taint string    // TODO: Check if correct type
