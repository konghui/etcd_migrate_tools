package main

import (
	"flag"
	"log"

	"github.com/coreos/etcd/pkg/transport"
	"github.com/konghui/etcd_migrate_tools/client"
	"github.com/konghui/etcd_migrate_tools/client/etcdv2"
	"github.com/konghui/etcd_migrate_tools/client/etcdv3"
	"github.com/konghui/etcd_migrate_tools/flags"
	"github.com/konghui/etcd_migrate_tools/glob"
)

func generateEtcdClient(flags *flags.EtcdFlags) (client.EtcdClient, error) {
	var tlsInfo transport.TLSInfo
	if *(flags.CertFile) != "" || *(flags.KeyFile) != "" || *(flags.CaFile) != "" {
		tlsInfo = transport.TLSInfo{
			CertFile:      *(flags.CertFile),
			KeyFile:       *(flags.KeyFile),
			TrustedCAFile: *(flags.CaFile),
		}
	}
	var etcdClient client.EtcdClient
	var err error

	if *flags.Version == 3 {
		etcdClient, err = etcdv3.New(flags, tlsInfo)
	} else {
		etcdClient, err = etcdv2.New(flags, tlsInfo)
	}
	if err != nil {
		return etcdClient, err
	}
	return etcdClient, err
}

func doMigrate(destFlags *flags.EtcdFlags, srcFlags *flags.EtcdFlags, parttern *glob.GlobalParttern) {
	destClient, err := generateEtcdClient(destFlags)
	if err != nil {
		log.Fatal("init dest client faild: %s.", err.Error())
	}
	srcClient, err := generateEtcdClient(srcFlags)
	if err != nil {
		log.Fatal("init src client faild: %s", err.Error())
	}

	srcData, err := srcClient.GetKeyValues(*srcFlags.Path)
	if err != nil {
		log.Fatal("faild to get data from src node: %s", err.Error())
	}

	total := len(srcData)
	current := 0
	log.Printf("get data success! total count = %d", total)
	for key, value := range srcData {
		current++
		if parttern.IsInBlackList(key) {
			log.Printf("key %s match black list, skip!\n", key)
			continue
		}
		log.Printf("start migrate key = %s to dest node, current %d, total %d\n", key, current, total)
		if (*destFlags.OverWrite) == true {
			log.Printf("key=%s, value=%s\n", key, value)
			//destClient.PutValue(key, value)
		} else {
			destData, err := destClient.GetKeyValues(key)
			if err != nil {
				log.Printf("error ocure when get data from dest node, key = %s, error: %s\n", key, err.Error())
				continue
			}
			if len(destData) != 0 {
				log.Printf("key=%s, in the dest node has data, and overwirte switcher is false, skip it.\n", key)
				continue
			}
			//destClient.PutValue(key, value)
		}
	}
}

func main() {

	version := flag.Int("version", 2, "use etcd version")
	path := flag.String("path", "/", "etcd migrate path, default is /")
	overwrite := flag.Bool("overwrite", false, "overwrite the data or not, default no")
	timeout := flag.Int64("timeout", 10, "etcd timeout time (second), default 10 second")

	destFlags := &flags.EtcdFlags{
		EtcdAddress: flag.String("dest-etcd-address", "", "Etcd address"),
		CertFile:    flag.String("dest-cert", "", "identify secure client using this TLS certificate file"),
		KeyFile:     flag.String("dest-key", "", "identify secure client using this TLS key file"),
		CaFile:      flag.String("dest-cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle"),
		Version:     version,
		Timeout:     timeout,
		Path:        path,
		OverWrite:   overwrite,
	}

	srcFlags := &flags.EtcdFlags{
		EtcdAddress: flag.String("src-etcd-address", "", "Etcd address"),
		CertFile:    flag.String("src-cert", "", "identify secure client using this TLS certificate file"),
		KeyFile:     flag.String("src-key", "", "identify secure client using this TLS key file"),
		CaFile:      flag.String("src-cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle"),
		Version:     version,
		Timeout:     timeout,
		Path:        path,
		OverWrite:   overwrite,
	}

	partternFile := flag.String("parttern-file", "parttern", "search the black/while parttern from file")

	flag.Parse()

	parttern, err := glob.NewGlobalParttern(*partternFile)
	log.Printf("partternFile = %s\n", *partternFile)
	if err != nil {
		log.Fatalf("error while create parrtern from file %s: %s", *partternFile, err.Error())
	}
	doMigrate(destFlags, srcFlags, parttern)

	/*etcdClient, err := generateEtcdClient(destFlags)
	if err != nil {
		fmt.Println(err.Error())
	}


	data, err := etcdClient.GetKeyValues("/")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(data)
	etcdClient.PutValue("/test/test1/222", []byte("12345677654321"))*/
}
