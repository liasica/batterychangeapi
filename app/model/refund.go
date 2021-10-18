package model

import (
    "battery/app/model/internal"
)

type Refund internal.Refund

const RefundRelationTypeComboOrder = 1

const (
    RefundStateStart = 1
    RefundStateDone  = 2
)
