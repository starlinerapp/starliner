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

const userNotificationChannel = "notifications:user"

type userNotificationEnvelope struct {
	UserId       int64                      `json:"userId"`
	Notification *value.ClusterNotification `json:"notification"`
}

type UserNotificationHub struct {
	pubsub port.PubSub

	mu      sync.RWMutex
	clients map[int64]map[chan *value.ClusterNotification]struct{}

	cancel context.CancelFunc
}

func NewUserNotificationHub(lc fx.Lifecycle, ps port.PubSub) *UserNotificationHub {
	h := &UserNotificationHub{
		pubsub:  ps,
		clients: make(map[int64]map[chan *value.ClusterNotification]struct{}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			sub, err := ps.Subscribe(ctx, userNotificationChannel)
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

func (h *UserNotificationHub) Subscribe(userId int64) apiport.UserNotificationSubscription {
	ch := make(chan *value.ClusterNotification, 16)
	h.mu.Lock()
	if h.clients[userId] == nil {
		h.clients[userId] = make(map[chan *value.ClusterNotification]struct{})
	}
	h.clients[userId][ch] = struct{}{}
	h.mu.Unlock()
	return &userSubscription{hub: h, userId: userId, ch: ch}
}

func (h *UserNotificationHub) unsubscribe(userId int64, ch chan *value.ClusterNotification) {
	h.mu.Lock()
	if subs, ok := h.clients[userId]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(h.clients, userId)
		}
	}
	h.mu.Unlock()
}

func (h *UserNotificationHub) Broadcast(userId int64, notification *value.ClusterNotification) {
	payload, err := json.Marshal(userNotificationEnvelope{
		UserId:       userId,
		Notification: notification,
	})
	if err != nil {
		log.Printf("failed to marshal user notification envelope: %v", err)
		return
	}

	if err := h.pubsub.Publish(context.Background(), userNotificationChannel, payload); err != nil {
		log.Printf("failed to publish user notification: %v", err)
	}
}

func (h *UserNotificationHub) run(ctx context.Context, sub port.Subscription) {
	defer func() {
		if err := sub.Close(); err != nil {
			log.Printf("failed to close user notification subscription: %v", err)
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
			var env userNotificationEnvelope
			if err := json.Unmarshal(raw, &env); err != nil {
				log.Printf("failed to unmarshal user notification envelope: %v", err)
				continue
			}
			h.deliverLocal(env.UserId, env.Notification)
		}
	}
}

func (h *UserNotificationHub) deliverLocal(userId int64, notification *value.ClusterNotification) {
	h.mu.RLock()
	subs := h.clients[userId]
	targets := make([]chan *value.ClusterNotification, 0, len(subs))
	for ch := range subs {
		targets = append(targets, ch)
	}
	h.mu.RUnlock()

	for _, ch := range targets {
		select {
		case ch <- notification:
		default:
			log.Printf("dropping notification for slow client on user %d", userId)
		}
	}
}
