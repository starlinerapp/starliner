package notification

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"go.uber.org/fx"
	apiport "starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

const envNotificationChannel = "notifications:environment"

type envNotificationEnvelope struct {
	CorrelationId string                         `json:"correlationId"`
	EnvironmentId int64                          `json:"environmentId"`
	Notification  *value.EnvironmentNotification `json:"notification"`
}

type environmentKey struct {
	correlationId string
	environmentId int64
}

type EnvironmentNotificationHub struct {
	pubsub port.PubSub

	mu      sync.RWMutex
	clients map[environmentKey]map[chan *value.EnvironmentNotification]struct{}

	cancel context.CancelFunc
}

func NewEnvironmentNotificationHub(lc fx.Lifecycle, ps port.PubSub) *EnvironmentNotificationHub {
	h := &EnvironmentNotificationHub{
		pubsub:  ps,
		clients: make(map[environmentKey]map[chan *value.EnvironmentNotification]struct{}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			sub, err := ps.Subscribe(ctx, envNotificationChannel)
			if err != nil {
				cancel()
				return err
			}
			go h.run(ctx, sub)
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})

	return h
}

func (h *EnvironmentNotificationHub) Subscribe(correlationId string, environmentId int64) apiport.EnvironmentNotificationSubscription {
	key := environmentKey{correlationId: correlationId, environmentId: environmentId}
	ch := make(chan *value.EnvironmentNotification, 16)
	h.mu.Lock()
	if h.clients[key] == nil {
		h.clients[key] = make(map[chan *value.EnvironmentNotification]struct{})
	}
	h.clients[key][ch] = struct{}{}
	h.mu.Unlock()
	return &environmentSubscription{hub: h, key: key, ch: ch}
}

func (h *EnvironmentNotificationHub) unsubscribe(key environmentKey, ch chan *value.EnvironmentNotification) {
	h.mu.Lock()
	if subs, ok := h.clients[key]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(h.clients, key)
		}
	}
	h.mu.Unlock()
}

func (h *EnvironmentNotificationHub) Broadcast(correlationId string, environmentId int64, notification *value.EnvironmentNotification) {
	payload, err := json.Marshal(envNotificationEnvelope{
		CorrelationId: correlationId,
		EnvironmentId: environmentId,
		Notification:  notification,
	})
	if err != nil {
		log.Printf("failed to marshal environment notification envelope: %v", err)
		return
	}

	if err := h.pubsub.Publish(context.Background(), envNotificationChannel, payload); err != nil {
		log.Printf("failed to publish environment notification: %v", err)
	}
}

func (h *EnvironmentNotificationHub) run(ctx context.Context, sub port.Subscription) {
	defer func() {
		if err := sub.Close(); err != nil {
			log.Printf("failed to close environment notification subscription: %v", err)
		}
	}()

	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case raw, ok := <-ch:
			if !ok {
				return
			}
			var env envNotificationEnvelope
			if err := json.Unmarshal(raw, &env); err != nil {
				log.Printf("failed to unmarshal environment notification envelope: %v", err)
				continue
			}
			h.deliverLocal(env.CorrelationId, env.EnvironmentId, env.Notification)
		}
	}
}

func (h *EnvironmentNotificationHub) deliverLocal(correlationId string, environmentId int64, notification *value.EnvironmentNotification) {
	key := environmentKey{correlationId: correlationId, environmentId: environmentId}
	h.mu.RLock()
	subs := h.clients[key]

	targets := make([]chan *value.EnvironmentNotification, 0, len(subs))
	for ch := range subs {
		targets = append(targets, ch)
	}
	h.mu.RUnlock()

	for _, ch := range targets {
		select {
		case ch <- notification:
		default:
			log.Printf("dropping notification for slow client on environment %d correlationId %s", environmentId, correlationId)
		}
	}
}
