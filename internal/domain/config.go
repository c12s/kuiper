package domain

import (
	"context"
	"fmt"
	"time"
)

type Node string

type Namespace string

type Org string

type PlacementReqStatus int8

const (
	PlacementReqStatusAccepted PlacementReqStatus = iota
	PlacementReqStatusPlaced
	PlacementReqStatusFailed
)

func (s PlacementReqStatus) String() string {
	switch s {
	case PlacementReqStatusAccepted:
		return "Accepted"
	case PlacementReqStatusPlaced:
		return "Placed"
	case PlacementReqStatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

type Config interface {
	Org() Org
	Name() string
	Version() string
	CreatedAtUnixSec() int64
	CreatedAtUTC() time.Time
}

type ConfigBase struct {
	org       Org
	version   string
	createdAt int64
}

func (c *ConfigBase) Org() Org {
	return c.org
}

func (c *ConfigBase) Version() string {
	return c.version
}

func (c *ConfigBase) CreatedAtUnixSec() int64 {
	return c.createdAt
}

func (c *ConfigBase) CreatedAtUTC() time.Time {
	return time.Unix(c.createdAt, 0).UTC()
}

type NamedParamSet struct {
	name   string
	params map[string]string
}

func NewParamSet(name string, params map[string]string) *NamedParamSet {
	return &NamedParamSet{
		name:   name,
		params: params,
	}
}

func (ps NamedParamSet) Name() string {
	return ps.name
}

func (ps NamedParamSet) ParamSet() map[string]string {
	return ps.params
}

func (ps NamedParamSet) Diff(cmp NamedParamSet) []Diff {
	// todo: ubaciti kod koji racuna diff
	return nil
}

type StandaloneConfig struct {
	ConfigBase
	paramSet NamedParamSet
}

func NewStandaloneConfig(org Org, version string, createdAt int64, paramSet NamedParamSet) *StandaloneConfig {
	return &StandaloneConfig{
		ConfigBase: ConfigBase{
			org:       org,
			version:   version,
			createdAt: createdAt,
		},
		paramSet: paramSet,
	}
}

func (c *StandaloneConfig) Name() string {
	return c.paramSet.name
}

func (c *StandaloneConfig) ParamSet() map[string]string {
	return c.paramSet.params
}

type ConfigGroup struct {
	ConfigBase
	name      string
	paramSets []NamedParamSet
}

func NewConfigGroup(org Org, name, version string, createdAt int64, paramSets []NamedParamSet) *ConfigGroup {
	return &ConfigGroup{
		ConfigBase: ConfigBase{
			org:       org,
			version:   version,
			createdAt: createdAt,
		},
		name:      name,
		paramSets: paramSets,
	}
}

func (c *ConfigGroup) Name() string {
	return c.name
}

func (c *ConfigGroup) ParamSets() []NamedParamSet {
	return c.paramSets
}

func (c *ConfigGroup) ParamSet(name string) (map[string]string, *Error) {
	for _, ps := range c.paramSets {
		if ps.name == name {
			return ps.params, nil
		}
	}
	return nil, NewError(ErrTypeNotFound, fmt.Sprintf("param set (name: %s) not found", name))
}

type StandaloneConfigStore interface {
	Put(ctx context.Context, config *StandaloneConfig) *Error
	Get(ctx context.Context, org Org, name, version string) (*StandaloneConfig, *Error)
	List(ctx context.Context, org Org) ([]*StandaloneConfig, *Error)
	Delete(ctx context.Context, org Org, name, version string) (*StandaloneConfig, *Error)
}

type ConfigGroupStore interface {
	Put(ctx context.Context, config *ConfigGroup) *Error
	Get(ctx context.Context, org Org, name, version string) (*ConfigGroup, *Error)
	List(ctx context.Context, org Org) ([]*ConfigGroup, *Error)
	Delete(ctx context.Context, org Org, name, version string) (*ConfigGroup, *Error)
}
