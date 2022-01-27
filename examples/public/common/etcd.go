package common

import (
	"strings"

	"github.com/ByronLiang/servant/examples/public/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var EtcdClusterClient *clientv3.Client

func InitEtcdClusterClient() error {
	var err error
	etcdCluster := strings.Split(config.EtcdHost, ",")
	EtcdClusterClient, err = clientv3.New(clientv3.Config{
		Endpoints: etcdCluster,
	})
	if err != nil {
		return err
	}
	return nil
}

func CloseEtcdClusterClient() {
	EtcdClusterClient.Close()
}
