package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/validation"
	"strings"
)

// Keys of custom validators
const (
	IP_IN_CIDR  = "ipInCidr"
	LB_REQUIRED = "lbRequired"
	VALID_HOST  = "validHost"
	VALID_POOL  = "validPool"
)

type Config struct {
	Hosts      []Host     `yaml:"hosts"`
	Cluster    Cluster    `yaml:"cluster"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
	Addons     Addons     `yaml:"addons,omitempty"`
	Kubitect   Kubitect   `yaml:"kubitect,omitempty"`
}

func (c Config) Validate() error {
	defer validation.ClearCustomValidators()

	validation.RegisterCustomValidator(IP_IN_CIDR, c.ipInCidrValidator())
	validation.RegisterCustomValidator(VALID_HOST, c.hostNameValidator())

	return validation.Struct(&c,
		validation.Field(&c.Hosts,
			validation.MinLen(1).Error("At least {.Param} host must be configured."),
			validation.UniqueField("Name"),
			c.singleDefaultHostValidator(),
		),
		validation.Field(&c.Cluster, validation.NotEmpty().Error("Configuration must contain '{.Field}' section.")),
		validation.Field(&c.Kubernetes, validation.NotEmpty().Error("Configuration must contain '{.Field}' section.")),
		validation.Field(&c.Addons),
		validation.Field(&c.Kubitect),
	)
}

func (c *Config) SetDefaults() {
	// If no host is set as the default host,
	// then set the first one as the default host.
	if len(c.Hosts) > 0 {
		for _, h := range c.Hosts {
			if h.Default {
				return
			}
		}

		c.Hosts[0].Default = true
	}
}

// singleDefaultHostValidator returns a validator that triggers an error
// if multiple hosts are configured as default.
func (c Config) singleDefaultHostValidator() validation.Validator {
	var defs int

	for _, h := range c.Hosts {
		if h.Default {
			defs++
		}
	}

	if defs > 1 {
		return validation.Fail().Errorf("Only one host can be configured as default.")
	}

	return validation.None
}

// ipInCidrValidator registers a custom validator that checks whether
// an IP address is within the configured network CIDR.
func (c Config) ipInCidrValidator() validation.Validator {
	return validation.IPInRange(string(c.Cluster.Network.CIDR))
}

// hostNameValidator returns a custom cross-validator that checks whether
// a host with a given name has been configured.
func (c Config) hostNameValidator() validation.Validator {
	var names []string

	for _, h := range c.Hosts {
		names = append(names, h.Name)
	}

	return validation.OneOf(names...).Errorf("Field '{.Field}' must point to one of the configured hosts: [%v] (actual: {.Value})", strings.Join(names, "|"))
}

// poolNameValidator returns a custom cross-validator that checks whether
// a given pool name is valid for a matching host.
func poolNameValidator(hostName string) validation.Validator {
	c, ok := validation.TopParent().(*Config)

	if !ok || c == nil || len(c.Hosts) == 0 {
		return validation.None
	}

	// By default, the first host in a list is a default host.
	host := (c.Hosts)[0]

	for _, h := range c.Hosts {
		if h.Default {
			host = h
		}

		if hostName == "" || h.Name == "" {
			continue
		}

		if hostName == h.Name {
			host = h
			break
		}
	}

	if len(host.DataResourcePools) == 0 {
		return validation.Fail().Errorf("Field '{.Field}' points to a data resource pool, but matching host '%v' has none configured.", host.Name)
	}

	var pools []string

	for _, p := range host.DataResourcePools {
		pools = append(pools, p.Name)
	}

	return validation.OneOf(pools...).Errorf("Field '{.Field}' must point to one of the pools configured on a matching host '%s': [%s] (actual: {.Value})", host.Name, strings.Join(pools, "|"))
}
