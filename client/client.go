package client

type EtcdClient interface {
	GetKeyValues(string) (map[string][]byte, error)
	PutValue(string, []byte) error
}
