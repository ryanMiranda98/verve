package mq

type MQ interface {
	Connect() error
	Close() error
	Write([]byte) (int, error)
}
