package update_getter

import (
	"github.com/darkenery/gobot/api"
	"github.com/darkenery/gobot/api/model"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
)

const lastUpdateIdRedisKey = "GoBot.LastUpdateId"

type UpdateGetter struct {
	botApi          *api.BotApi
	updateHandlerCh chan []*model.Update
	redis           *redis.ClusterClient
	limit           int
	timeout         int
	allowedUpdates  []string
	logger          log.Logger
}

func NewUpdateGetter(botApi *api.BotApi, ch chan []*model.Update, redis *redis.ClusterClient, limit, timeout int, allowedUpdates []string, logger log.Logger) *UpdateGetter {
	return &UpdateGetter{
		botApi:          botApi,
		updateHandlerCh: ch,
		redis:           redis,
		limit:           limit,
		timeout:         timeout,
		allowedUpdates:  allowedUpdates,
		logger:          logger,
	}
}

func (u *UpdateGetter) GetUpdates() {
	lastUpdateIdInt64, err := u.redis.Get(lastUpdateIdRedisKey).Int64()
	if err != nil && err != redis.Nil {
		u.logger.Log("err", err)
		return
	}

	lastUpdateId := int(lastUpdateIdInt64)

	for {
		updates, err := u.botApi.GetUpdates(
			lastUpdateId+1,
			u.limit,
			u.timeout,
			u.allowedUpdates,
		)

		if err != nil {
			u.logger.Log("err", err)
			continue
		}

		if len(updates) == 0 {
			continue
		}

		lastUpdateId = updates[len(updates) - 1].UpdateId
		err = u.redis.Set(lastUpdateIdRedisKey, lastUpdateId, 0).Err()
		if err != nil {
			u.logger.Log("err", err)
		}

		u.updateHandlerCh <- updates
	}
}
