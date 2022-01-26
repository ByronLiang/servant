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
	session  *concurrency.Session
	election *concurrency.Election
	opts     *leaderFollowerOptions
}

type LeaderFollowerOption func(o *leaderFollowerOptions)

type leaderFollowerOptions struct {
	heartBeatTTL int
	prefix       string
	value        string
}

func HeartBeatTTL(ttl int) LeaderFollowerOption {
	return func(o *leaderFollowerOptions) {
		o.heartBeatTTL = ttl
	}
}

func Prefix(prefix string) LeaderFollowerOption {
	return func(o *leaderFollowerOptions) {
		o.prefix = prefix
	}
}

func Value(value string) LeaderFollowerOption {
	return func(o *leaderFollowerOptions) {
		o.value = value
	}
}

func NewLeaderFollowerRegistry(client *clientv3.Client, options ...LeaderFollowerOption) (*leaderFollowerRegistry, error) {
	ctx, cancel := context.WithCancel(context.Background())
	opts := &leaderFollowerOptions{
		heartBeatTTL: 5,
	}
	for _, option := range options {
		option(opts)
	}
	// 若正常退出, 触发resign  而 ttl 针对非resign 下故障容忍时长触发重新选举
	session, err := concurrency.NewSession(client, concurrency.WithTTL(opts.heartBeatTTL))
	if err != nil {
		return nil, err
	}
	return &leaderFollowerRegistry{
		ctx:     ctx,
		session: session,
		cancel:  cancel,
		opts:    opts,
	}, nil
}

func (r *leaderFollowerRegistry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	key := fmt.Sprintf("%s/%s", r.opts.value, service.Name)
	r.opts.prefix = key
	marshalStr, err := marshal(service)
	if err != nil {
		return err
	}
	r.opts.value = marshalStr
	election := concurrency.NewElection(r.session, r.opts.prefix)
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
		// compare register endpoint
		if leaderInstance.Endpoints[0] == service.Endpoints[0] {
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
		err := r.election.Campaign(r.ctx, r.opts.value)
		campaignRes <- err
	}()
	return campaignRes
}

func (r *leaderFollowerRegistry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.prefix, name)
	return newLeaderFollowerWatcher(ctx, key, r.session)
}

func (r *leaderFollowerRegistry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	prefix := fmt.Sprintf("%s/%s/", r.opts.prefix, name)
	resp, err := r.session.Client().Get(ctx, prefix, clientv3.WithFirstCreate()...)
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
