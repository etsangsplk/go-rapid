package paxos

import (
	"context"
	"errors"

	"github.com/Workiva/go-datastructures/trie/ctrie"
	"github.com/sasha-s/go-deadlock"

	"google.golang.org/protobuf/encoding/protojson"

	"google.golang.org/protobuf/proto"

	"github.com/rs/zerolog"

	"github.com/casualjim/go-rapid/api"

	"github.com/casualjim/go-rapid/remoting"
)

type ConsensusRegistry struct {
	//data *hashmap.HashMap
	data *ctrie.Ctrie
	//data map[*remoting.Endpoint]*Fast
	//lock deadlock.Mutex
}

func (c *ConsensusRegistry) key(ep *remoting.Endpoint) []byte {
	b, err := proto.Marshal(ep)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *ConsensusRegistry) Get(ep *remoting.Endpoint) *Fast {
	v, ok := c.data.Lookup(c.key(ep))
	if !ok {
		return nil
	}
	return v.(*Fast)
}

func (c *ConsensusRegistry) Each(handle func(*remoting.Endpoint, *Fast)) {

	for entry := range c.data.Iterator(nil) {
		var k remoting.Endpoint
		if err := proto.Unmarshal(entry.Key, &k); err != nil {
			panic(err)
		}
		handle(&k, entry.Value.(*Fast))
	}
}

func (c *ConsensusRegistry) GetOrSet(ep *remoting.Endpoint, loader func(ep *remoting.Endpoint) *Fast) *Fast {
	v, ok := c.data.Lookup(c.key(ep))
	if ok {
		return v.(*Fast)
	}
	f := loader(ep)
	c.Set(ep, f)
	return f
}

func (c *ConsensusRegistry) Set(ep *remoting.Endpoint, cc *Fast) {
	c.data.Insert(c.key(ep), cc)
}

func (c *ConsensusRegistry) EachKey(handle func(*remoting.Endpoint)) {
	for entry := range c.data.Iterator(nil) {
		var k remoting.Endpoint
		if err := proto.Unmarshal(entry.Key, &k); err != nil {
			panic(err)
		}
		handle(&k)
	}
}

func (c *ConsensusRegistry) First() *Fast {
	var f *Fast
	for v := range c.data.Iterator(nil) {
		if f == nil {
			f = v.Value.(*Fast)
		}
	}
	return f
}

type DirectBroadcaster struct {
	log            zerolog.Logger
	mtypes         *typeRegistry
	paxosInstances *ConsensusRegistry
	client         api.Client
}

func (d *DirectBroadcaster) Broadcast(ctx context.Context, req *remoting.RapidRequest) {
	if d.mtypes.Get(req.Content) {
		d.log.Info().Msg("exiting broadcast becuase of types filter")
		return
	}

	d.paxosInstances.EachKey(func(k *remoting.Endpoint) {
		d.log.Info().Str("to", endpointStr(k)).Str("message", protojson.Format(req)).Msg("broadcasting")
		go func() {
			_, err := d.client.Do(ctx, k, req)
			if err != nil {
				return
			}
		}()
	})
}

func (d *DirectBroadcaster) SetMembership([]*remoting.Endpoint) { panic("not supported") }

func (d *DirectBroadcaster) Start() {}

func (d *DirectBroadcaster) Stop() {}

type DirectClient struct {
	paxosInstances *ConsensusRegistry
	lock           deadlock.Mutex
}

func (d *DirectClient) Do(ctx context.Context, target *remoting.Endpoint, in *remoting.RapidRequest) (*remoting.RapidResponse, error) {
	go func(in *remoting.RapidRequest) {
		inst := d.paxosInstances.Get(target)
		d.lock.Lock()
		_, _ = inst.Handle(ctx, in)
		d.lock.Unlock()
	}(in)
	return &remoting.RapidResponse{}, nil
}
func (d *DirectClient) DoBestEffort(ctx context.Context, target *remoting.Endpoint, in *remoting.RapidRequest) (*remoting.RapidResponse, error) {
	return nil, errors.New("not supported")
}

func (d *DirectClient) Close() error {
	return errors.New("not supported")
}

type NoopClient struct{}

func (n *NoopClient) Do(ctx context.Context, target *remoting.Endpoint, in *remoting.RapidRequest) (*remoting.RapidResponse, error) {
	return nil, nil
}
func (n *NoopClient) DoBestEffort(ctx context.Context, target *remoting.Endpoint, in *remoting.RapidRequest) (*remoting.RapidResponse, error) {
	return nil, nil
}

func (n *NoopClient) Close() error {
	return nil
}

type NoopBroadcaster struct{}

func (n *NoopBroadcaster) Broadcast(ctx context.Context, req *remoting.RapidRequest) {
}

func (n *NoopBroadcaster) SetMembership([]*remoting.Endpoint) {}
func (n *NoopBroadcaster) Start()                             {}

func (n *NoopBroadcaster) Stop() {}
