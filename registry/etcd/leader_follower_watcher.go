package etcd

import (
	"context"

	"github.com/ByronLiang/servant/registry"

	"go.etcd.io/etcd/client/v3/concurrency"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type leaderFollowerWatcher struct {
	ctx       context.Context
	cancel    context.CancelFunc
	session   *concurrency.Session
	election  *concurrency.Election
	watchChan <-chan clientv3.GetResponse
	key       string
	first     bool
}

func (l *leaderFollowerWatcher) Next() ([]*registry.ServiceInstance, error) {
	if l.first {
		l.watchChan = l.election.Observe(l.ctx)
		l.first = false
	}
	select {
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	case <-l.watchChan:
		return l.getInstance()
	}
}

func (l *leaderFollowerWatcher) getInstance() ([]*registry.ServiceInstance, error) {
	resp, err := l.election.Leader(l.ctx)
	if err != nil {
		return nil, err
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

func (l *leaderFollowerWatcher) Stop() error {
	l.cancel()
	return l.session.Close()
}

func newLeaderFollowerWatcher(cctx context.Context, key string, session *concurrency.Session) (*leaderFollowerWatcher, error) {
	ctx, cancel := context.WithCancel(cctx)
	w := &leaderFollowerWatcher{
		ctx:    ctx,
		cancel: cancel,
		key:    key,
		first:  true,
	}
	prefix := w.key + "/"
	resp, err := session.Client().Get(ctx, prefix, clientv3.WithFirstCreate()...)
	if err != nil {
		return nil, err
	} else if len(resp.Kvs) == 0 {
		// no leader currently elected
		return nil, concurrency.ErrElectionNoLeader
	}
	w.session = session
	w.election = concurrency.ResumeElection(w.session, prefix,
		string(resp.Kvs[0].Key), resp.Kvs[0].CreateRevision)
	return w, nil
}
