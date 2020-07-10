package broker

import (
	"container/list"
	"context"
	"sync"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/util/log"
	"github.com/wmsx/xconf/config-srv/broadcast"
	"github.com/wmsx/xconf/proto/config"
)

var _ broadcast.Broadcast = &Broker{}

const broadcastUpdateTopic = "go.micro.xconf.broadcast.update"

type Broker struct {
	sync.RWMutex
	service   micro.Service
	publisher micro.Publisher

	watchers list.List
}

func New(service micro.Service) (broadcast.Broadcast, error) {
	broker := &Broker{
		service: service,
	}

	if err := micro.RegisterSubscriber(broadcastUpdateTopic, service.Server(), broker.subEvent); err != nil {
		return nil, err
	}
	broker.publisher = micro.NewPublisher(broadcastUpdateTopic, service.Client())

	return broker, nil
}

func (b *Broker) Send(namespace *config.ConfigResponse) error {
	return b.publisher.Publish(context.Background(), namespace)
}

func (b *Broker) Watch() broadcast.Watcher {
	w := &Watcher{
		exit:    make(chan interface{}),
		updates: make(chan *config.ConfigResponse, 2), // TODO 1 ?? 2 ?? or config
	}

	b.Lock()
	e := b.watchers.PushBack(w)
	b.Unlock()

	go func() {
		<-w.exit
		b.Lock()
		b.watchers.Remove(e)
		b.Unlock()
	}()

	return w
}

func (b *Broker) subEvent(ctx context.Context, event *config.ConfigResponse) error {
	md, _ := metadata.FromContext(ctx)
	log.Infof("[pubsub.2] Received event %+v with metadata %+v\n", event, md)

	watchers := make([]*Watcher, 0, b.watchers.Len())
	b.RLock()
	for e := b.watchers.Front(); e != nil; e = e.Next() {
		watchers = append(watchers, e.Value.(*Watcher))
	}
	b.RUnlock()

	for _, w := range watchers {
		select {
		case w.updates <- event:
		default:
		}
	}

	return nil
}
