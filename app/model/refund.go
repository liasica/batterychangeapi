package model

import (
	"battery/app/model/internal"
)

type Refund internal.Refund

const RefundRelationTypePackagesOrder = 1

const (
	RefundStateStart = 1
	RefundStateDone  = 2
)
