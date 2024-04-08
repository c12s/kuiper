package domain

import (
	"context"
	"fmt"
	"slices"
	"time"
)

type Node string

type Namespace string

type Org string

type Config interface {
	Org() Org
	Name() string
	Version() string
	CreatedAtUnixSec() int64
	CreatedAtUTC() time.Time
	Revision() int64
	Nodes() []Node
	AddToNode(node Node) bool
	Namespaces() []Namespace
	AddToNamespace(namespace Namespace) bool
}

type ConfigBase struct {
	org        Org
	version    string
	createdAt  int64
	revision   int64
	nodes      []Node
	namespaces []Namespace
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

func (c *ConfigBase) Revision() int64 {
	return c.revision
}

func (c *ConfigBase) Nodes() []Node {
	return c.nodes
}

func (c *ConfigBase) AddToNode(node Node) bool {
	if slices.Contains(c.nodes, node) {
		return false
	}
	c.nodes = append(c.nodes, node)
	return true
}

func (c *ConfigBase) Namespaces() []Namespace {
	return c.namespaces
}

func (c *ConfigBase) AddToNamespace(namespace Namespace) bool {
	if slices.Contains(c.namespaces, namespace) {
		return false
	}
	c.namespaces = append(c.namespaces, namespace)
	return true
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

func (ps NamedParamSet) Diff(cmp NamedParamSet) []Diff {
	// todo: ubaciti kod koji racuna diff
	return nil
}

type StandaloneConfig struct {
	ConfigBase
	paramSet NamedParamSet
}

func NewStandaloneConfigWithoutRevision(org Org, version string, createdAt int64, paramSet NamedParamSet) *StandaloneConfig {
	return &StandaloneConfig{
		ConfigBase: ConfigBase{
			org:        org,
			version:    version,
			createdAt:  createdAt,
			nodes:      make([]Node, 0),
			namespaces: make([]Namespace, 0),
		},
		paramSet: paramSet,
	}
}

func NewStandaloneConfig(org Org, version string, createdAt, revision int64, paramSet NamedParamSet) *StandaloneConfig {
	return &StandaloneConfig{
		ConfigBase: ConfigBase{
			org:        org,
			version:    version,
			createdAt:  createdAt,
			revision:   revision,
			nodes:      make([]Node, 0),
			namespaces: make([]Namespace, 0),
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

func NewConfigGroupWithoutRevision(org Org, name, version string, createdAt int64, paramSets []NamedParamSet) *ConfigGroup {
	return &ConfigGroup{
		ConfigBase: ConfigBase{
			org:        org,
			version:    version,
			createdAt:  createdAt,
			nodes:      make([]Node, 0),
			namespaces: make([]Namespace, 0),
		},
		name:      name,
		paramSets: paramSets,
	}
}

func NewConfigGroup(org Org, name, version string, createdAt, revision int64, paramSets []NamedParamSet) *ConfigGroup {
	return &ConfigGroup{
		ConfigBase: ConfigBase{
			org:        org,
			version:    version,
			createdAt:  createdAt,
			revision:   revision,
			nodes:      make([]Node, 0),
			namespaces: make([]Namespace, 0),
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
	Get(ctx context.Context, Org, name, version string) (*StandaloneConfig, *Error)
	GetHistory(ctx context.Context, Org, name, version string) ([]*StandaloneConfig, *Error)
	List(ctx context.Context) ([]*StandaloneConfig, *Error)
	Delete(ctx context.Context, Org, name, version string) (*StandaloneConfig, *Error)
	AddToNodes(ctx context.Context, config *StandaloneConfig) *Error
	AddToNamespaces(ctx context.Context, config *StandaloneConfig) *Error
	ListNode(ctx context.Context, node Node, org Org) ([]*StandaloneConfig, *Error)
	ListNamespace(ctx context.Context, namespace Namespace, org Org) ([]*StandaloneConfig, *Error)
}

type ConfigGroupStore interface {
	Put(ctx context.Context, config *ConfigGroup) *Error
	Get(ctx context.Context, Org, name, version string) (*ConfigGroup, *Error)
	GetHistory(ctx context.Context, Org, name, version string) ([]*ConfigGroup, *Error)
	List(ctx context.Context) ([]*ConfigGroup, *Error)
	Delete(ctx context.Context, Org, name, version string) (*ConfigGroup, *Error)
	AddToNodes(ctx context.Context, config *ConfigGroup) *Error
	AddToNamespaces(ctx context.Context, config *ConfigGroup) *Error
	ListNode(ctx context.Context, node Node, org Org) ([]*ConfigGroup, *Error)
	ListNamespace(ctx context.Context, namespace Namespace, org Org) ([]*ConfigGroup, *Error)
}
