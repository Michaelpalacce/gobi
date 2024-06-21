package processor_v1

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/redis"
)

func (p *Processor) subscribeToRedis() {
	chanName := p.WebsocketClient.User.Username + "-" + p.WebsocketClient.Client.VaultName
	redisChan := redis.Subscribe(chanName).Channel()
	slog.Info("Subscribed to Redis channel", "channel", chanName)

	for {
		msg := <-redisChan
		fmt.Println(msg.Payload)
	}
}
