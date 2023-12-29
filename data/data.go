package data

const (
	PrefixSeparator = byte('.')
)

type StorableData interface {
	Key() []byte
	Value() []byte
}
