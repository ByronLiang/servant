package etcd

import (
	"context"
	"fmt"

	"github.com/ByronLiang/servant/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type leaderFollowerRegistry struct {
	ctx      context.Context
	cancel   context.CancelFunc
	client   *clientv3.Client
	session  *concurrency.Session
	election *concurrency.Election
	opts     *LeaderFollowerOptions
}

type LeaderFollowerOptions struct {
	HeartBeatTTL int
	Prefix       string
	Value        string
}

func NewLeaderFollowerRegistry(client *clientv3.Client, opts *LeaderFollowerOptions) (*leaderFollowerRegistry, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &leaderFollowerRegistry{
		ctx:    ctx,
		client: client,
		cancel: cancel,
		opts:   opts,
	}, nil
}

func (r *leaderFollowerRegistry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	key := fmt.Sprintf("%s/%s", r.opts.Value, service.Name)
	r.opts.Prefix = key
	marshalStr, err := marshal(service)
	if err != nil {
		return err
	}
	r.opts.Value = marshalStr
	// 若正常退出, 触发resign  而 ttl 针对非resign 下故障容忍时长触发重新选举
	session, err := concurrency.NewSession(r.client, concurrency.WithTTL(r.opts.HeartBeatTTL))
	if err != nil {
		return err
	}
	election := concurrency.NewElection(session, r.opts.Prefix)
	r.session = session
	r.election = election
	campaignRes := r.campaign()
	err = <-campaignRes
	return err
}

func (r *leaderFollowerRegistry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	// 查看当前leader
	res, err := r.election.Leader(ctx)
	if err == nil {
		leaderInstance, err := unmarshal(res.Kvs[0].Value)
		if err != nil {
			return err
		}
		currentInstance, err := unmarshal([]byte(r.opts.Value))
		if err != nil {
			return err
		}
		// compare register endpoint
		if leaderInstance.Endpoints[0] == currentInstance.Endpoints[0] {
			// 优雅停机: 从备份节点重新选举出 leader
			err = r.election.Resign(ctx)
			if err != nil {
				return err
			}
		} else {
			// 解除参与选举leader
			r.cancel()
		}
	}
	return r.session.Close()
}

func (r *leaderFollowerRegistry) campaign() chan error {
	campaignRes := make(chan error)
	go func() {
		err := r.election.Campaign(r.ctx, r.opts.Value)
		campaignRes <- err
	}()
	return campaignRes
}

func (r *leaderFollowerRegistry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.Prefix, name)
	return newLeaderFollowerWatcher(ctx, key, r.client)
}

func (r *leaderFollowerRegistry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	prefix := fmt.Sprintf("%s/%s/", r.opts.Prefix, name)
	resp, err := r.client.Get(ctx, prefix, clientv3.WithFirstCreate()...)
	if err != nil {
		return nil, err
	} else if len(resp.Kvs) == 0 {
		// no leader currently elected
		return nil, concurrency.ErrElectionNoLeader
	}
	items := make([]*registry.ServiceInstance, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		items[i] = si
	}
	return items, nil
}
