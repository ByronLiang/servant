package etcd

import etcdclient "go.etcd.io/etcd/client/v3"

func NewEtcdCli() {
	etcdclient.New()
}
