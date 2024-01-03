package data

import "google.golang.org/protobuf/types/known/timestamppb"

var (
	PrefixSeparator = []byte{byte('.')}
)

type StorableData interface {
	Key() []byte
	Value() []byte
	GetCreatedAt() *timestamppb.Timestamp
}
