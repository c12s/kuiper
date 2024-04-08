package domain

import (
	"context"
	"time"
)

type PlacementReq struct {
	id         string
	node       Node
	namespace  Namespace
	status     PlacementReqStatus
	acceptedAt int64
	resolvedAt int64
}

func NewPlacementReq(id string, node Node, namepsace Namespace, status PlacementReqStatus, acceptedAt, resolvedAt int64) *PlacementReq {
	return &PlacementReq{
		id:         id,
		node:       node,
		namespace:  namepsace,
		status:     status,
		acceptedAt: acceptedAt,
		resolvedAt: resolvedAt,
	}
}

func (p *PlacementReq) Id() string {
	return p.id
}

func (p *PlacementReq) Node() Node {
	return p.node
}

func (p *PlacementReq) Namespace() Namespace {
	return p.namespace
}

func (p *PlacementReq) AcceptedAtUnixSec() int64 {
	return p.acceptedAt
}

func (p *PlacementReq) AcceptedAtUTC() time.Time {
	return time.Unix(p.acceptedAt, 0).UTC()
}

func (p *PlacementReq) ResolvedAtUnixSec() int64 {
	return p.resolvedAt
}

func (p *PlacementReq) ResolveddAtUTC() time.Time {
	return time.Unix(p.resolvedAt, 0).UTC()
}

func (p *PlacementReq) Resolved() bool {
	return p.status != PlacementReqStatusAccepted
}

func (p *PlacementReq) Status() PlacementReqStatus {
	return p.status
}

type PlacementStore interface {
	Place(ctx context.Context, config Config, req *PlacementReq) *Error
	GetPlacement(ctx context.Context, org Org, name, version string) ([]*PlacementReq, *Error)
}
