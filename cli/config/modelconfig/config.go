package modelconfig

import (
	v "cli/validation"
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
	Hosts      *[]Host     `yaml:"hosts"`
	Cluster    *Cluster    `yaml:"cluster"`
	Kubernetes *Kubernetes `yaml:"kubernetes"`
	Addons     *Addons     `yaml:"addons"`
}

func (c Config) Validate() error {
	defer v.ClearCustomValidators()

	v.RegisterCustomValidator(IP_IN_CIDR, c.ipInCidrValidator())
	v.RegisterCustomValidator(VALID_HOST, c.hostNameValidator())

	return v.Struct(&c,
		v.Field(&c.Hosts,
			v.MinLen(1).Error("At least {.Param} {.Field} must be configured."),
			v.UniqueField("Name"),
		),
		v.Field(&c.Cluster, v.Required().Error("Configuration must contain '{.Field}' section.")),
		v.Field(&c.Kubernetes, v.Required().Error("Configuration must contain '{.Field}' section.")),
		v.Field(&c.Addons, v.OmitEmpty()),
	)
}

// ipInCidrValidator registers a custom validator that checks whether
// an IP address is within the configured network CIDR.
func (c Config) ipInCidrValidator() v.Validator {
	if c.Cluster != nil && c.Cluster.Network != nil && c.Cluster.Network.CIDR != nil {
		return v.None
	}

	return v.IPInRange(string(*c.Cluster.Network.CIDR))
}

// hostNameValidator returns a custom cross-validator that checks whether
// a host with a given name has been configured.
func (c Config) hostNameValidator() v.Validator {
	if c.Hosts == nil {
		return v.None
	}

	var hostNames []string

	for _, h := range *c.Hosts {
		if h.Name != nil {
			hostNames = append(hostNames, *h.Name)
		}
	}

	return v.OneOf(hostNames).Errorf("Field '{.Field}' must point to one of the configured host: [%s]", strings.Join(hostNames, "|"))
}

// poolNameValidator returns a custom cross-validator that checks whether
// a given pool name is valid for a matching host.
func poolNameValidator(hostName *string) v.Validator {

	c, ok := v.TopParent().(*Config)
	if !ok || c == nil {
		return v.None
	}

	if c.Hosts == nil || len(*c.Hosts) == 0 {
		return v.None
	}

	// By default, the first host in a list is a default host.
	host := (*c.Hosts)[0]

	for _, h := range *c.Hosts {
		if h.Default != nil && *h.Default {
			host = h
		}

		if hostName == nil || h.Name == nil {
			continue
		}

		if *h.Name == *hostName {
			host = h
			break
		}
	}

	if host.Name == nil {
		// Ignore, because in such case an error is already triggered for a host.
		return v.None
	}

	pools := host.DataResourcePools

	if pools == nil || len(*pools) == 0 {
		return v.Fail().Errorf("Field '{.Field}' points to a data resource pool, but matching host '%s' has none configured.", host)
	}

	var poolNames []string

	for _, p := range *host.DataResourcePools {
		if p.Name != nil {
			poolNames = append(poolNames, *p.Name)
		}
	}

	return v.OneOf(poolNames).Errorf("Field '{.Field}' must point to one of the pools configured on a matching host '%s': [%s]", *host.Name, strings.Join(poolNames, "|"))
}
