package etcdv3

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/konghui/etcd_migrate_tools/flags"
)

type ETCDV3Client struct {
	client  *clientv3.Client
	timeout time.Duration
}

type etcdFlags struct {
	etcdAddress *string
	certFile    *string
	keyFile     *string
	caFile      *string
	version     *int
	timeout     *int64
}

func New(flags *flags.EtcdFlags, tlsInfo transport.TLSInfo) (*ETCDV3Client, error) {

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}
	var client *clientv3.Client

	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{*(flags.EtcdAddress)},
		DialTimeout: time.Duration(*(flags.Timeout)) * time.Second,
		TLS:         tlsConfig,
	})

	if err != nil {
		return nil, err
	}

	ETCDV3Client := &ETCDV3Client{
		client:  client,
		timeout: time.Duration(*(flags.Timeout)) * time.Second,
	}
	return ETCDV3Client, err
}

func (this *ETCDV3Client) GetKeyValues(keys string) (map[string][]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	resp, err := this.client.Get(ctx, keys, clientv3.WithPrefix())
	cancel()

	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte, 100)
	for _, ev := range resp.Kvs {
		result[string(ev.Key)] = ev.Value
	}

	return result, err
}

func (this *ETCDV3Client) PutValue(key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	_, err := this.client.Put(ctx, key, string(value))
	cancel()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}
