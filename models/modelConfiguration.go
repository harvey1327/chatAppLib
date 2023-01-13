package models

type ModelConfiguration[T any] interface {
	GetQueueName() string
}
