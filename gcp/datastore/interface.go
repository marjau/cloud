package datastore

type DataFilter interface {
	GetField() string
	GetCondition() string
	GetValue() interface{}
}
