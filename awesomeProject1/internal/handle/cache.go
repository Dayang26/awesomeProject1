package handle

import (
	g "awesomeProject1/internal/global"
	"awesomeProject1/internal/model"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

// redis context
var rctx = context.Background()

// 将页面列表缓存在Redis中
func addPageCache(rdb *redis.Client, pages []model.Page) error {
	data, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	return rdb.Set(rctx, g.PAGE, string(data), 0).Err()
}

func removePageCache(rdb *redis.Client) error { return rdb.Del(rctx, g.PAGE).Err() }

// get page cache
func getPageCache(rdb *redis.Client) (cache []model.Page, err error) {
	s, err := rdb.Get(rctx, g.PAGE).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(s), &cache); err != nil {
		return nil, err
	}
	return cache, nil
}

func addConfigCache(rdb *redis.Client, config map[string]string) error {
	return rdb.HMSet(rctx, g.CONFIG, config).Err()
}

func removeConfigCache(rdb *redis.Client) error {
	return rdb.Del(rctx, g.CONFIG).Err()
}

func getConfigCache(rdb *redis.Client) (cache map[string]string, err error) {
	return rdb.HGetAll(rctx, g.CONFIG).Result()
}
