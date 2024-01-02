package data

import "google.golang.org/protobuf/types/known/timestamppb"

const (
	PrefixSeparator = byte('.')
)

type StorableData interface {
	Key() []byte
	Value() []byte
	GetCreatedAt() *timestamppb.Timestamp
	GetUpdatedAt() *timestamppb.Timestamp
}
