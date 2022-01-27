package etcd

import (
	"context"
	"fmt"

	"github.com/ByronLiang/servant/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type LeaderFollowerRegistry struct {
	ctx         context.Context
	cancel      context.CancelFunc
	session     *concurrency.Session
	election    *concurrency.Election
	opts        *leaderFollowerOptions
	CampaignRes chan error
}

type LeaderFollowerOption func(o *leaderFollowerOptions)

type leaderFollowerOptions struct {
	heartBeatTTL int
	namespace    string
	prefix       string
}

func HeartBeatTTL(ttl int) LeaderFollowerOption {
	return func(o *leaderFollowerOptions) {
		o.heartBeatTTL = ttl
	}
}

func NameSpace(namespace string) LeaderFollowerOption {
	return func(o *leaderFollowerOptions) {
		o.namespace = namespace
	}
}

func NewLeaderFollowerRegistry(client *clientv3.Client, options ...LeaderFollowerOption) (*LeaderFollowerRegistry, error) {
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
	return &LeaderFollowerRegistry{
		ctx:     ctx,
		session: session,
		cancel:  cancel,
		opts:    opts,
	}, nil
}

func (r *LeaderFollowerRegistry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	r.opts.prefix = fmt.Sprintf("%s/%s", r.opts.namespace, service.Name)
	marshalStr, err := marshal(service)
	if err != nil {
		return err
	}
	election := concurrency.NewElection(r.session, r.opts.prefix)
	r.election = election
	r.campaign(marshalStr)
	return nil
}

func (r *LeaderFollowerRegistry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	r.election.Resign(ctx)
	r.cancel()
	return r.session.Close()
}

func (r *LeaderFollowerRegistry) campaign(value string) {
	r.CampaignRes = make(chan error)
	go func() {
		err := r.election.Campaign(r.ctx, value)
		r.CampaignRes <- err
		close(r.CampaignRes)
	}()
}

func (r *LeaderFollowerRegistry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	r.opts.prefix = fmt.Sprintf("%s/%s", r.opts.namespace, name)
	return newLeaderFollowerWatcher(ctx, r.opts.prefix, r.session)
}

func (r *LeaderFollowerRegistry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
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
