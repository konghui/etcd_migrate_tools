package etcdv2

import (
	"context"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/konghui/etcd_migrate_tools/flags"
)

type ETCDV2Client struct {
	client   client.Client
	kvclient client.KeysAPI
	timeout  time.Duration
	options  *client.GetOptions
}

func New(flags *flags.EtcdFlags, tlsInfo transport.TLSInfo) (*ETCDV2Client, error) {

	trans, err := transport.NewTransport(tlsInfo, 30*time.Second)
	if err != nil {
		return nil, err
	}
	config := client.Config{
		Endpoints:               []string{*(flags.EtcdAddress)},
		HeaderTimeoutPerRequest: time.Duration(*(flags.Timeout)) * time.Second,
		Transport:               trans,
	}

	v2client, err := client.New(config)
	if err != nil {
		return nil, err
	}

	ETCDV2Client := &ETCDV2Client{
		client:   v2client,
		kvclient: client.NewKeysAPI(v2client),
		timeout:  time.Duration(*(flags.Timeout)) * time.Second,
		options:  &client.GetOptions{Sort: true, Recursive: true, Quorum: true},
	}

	return ETCDV2Client, err
}

func (this *ETCDV2Client) GetKeyValues(key string) (map[string][]byte, error) {
	data := make(map[string][]byte, 100)
	err := this.KeyValuesWalker(key, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *ETCDV2Client) KeyValuesWalker(key string, data map[string][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	resp, err := this.kvclient.Get(ctx, key, this.options)
	cancel()
	if err != nil {
		return err
	}

	for _, node := range resp.Node.Nodes {
		if node.Dir == true {
			err = this.KeyValuesWalker(node.Key, data)
			if err != nil {
				return err
			}
		} else {
			data[string(node.Key)] = []byte(node.Value)
		}
	}
	return nil
}

func (this *ETCDV2Client) PutValue(key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	_, err := this.kvclient.Create(ctx, key, string(value))
	cancel()
	if err != nil {
		return err
	}
	return nil
}
