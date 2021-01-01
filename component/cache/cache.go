/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/12/20
   Description :
-------------------------------------------------
*/

package cache

import (
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zlyuancn/zcache"
	memory_cache "github.com/zlyuancn/zcache/cachedb/memory-cache"
	no_cache "github.com/zlyuancn/zcache/cachedb/no-cache"
	redis_cache "github.com/zlyuancn/zcache/cachedb/redis-cache"
	no_sf "github.com/zlyuancn/zcache/single_flight/no-sf"
	single_sf "github.com/zlyuancn/zcache/single_flight/single-sf"
	"go.uber.org/zap"

	cache_codec "github.com/zlyuancn/zcache/codec"
	cache_core "github.com/zlyuancn/zcache/core"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

type Cache struct {
	caches map[string]*zcache.Cache
}

func (c *Cache) GetCache(name ...string) *zcache.Cache {
	n := consts.DefaultComponentName
	if len(name) > 0 {
		n = name[0]
	}
	cache := c.caches[n]

	if cache == nil {
		logger.Log.Panic(zap.String("name", n), zap.Error(errors.New("试图获取未注册的cache")))
	}

	return cache
}

func NewCache(app core.IApp) core.ICache {
	configs := app.GetConfig().Config().Components.Cache
	caches := make(map[string]*zcache.Cache, len(configs))
	for name, conf := range configs {
		cache := zcache.NewCache(
			zcache.WithCacheDB(makeCacheDB(&conf)),
			zcache.WithDefaultExpire(time.Duration(conf.DefaultExpire)*time.Millisecond, time.Duration(conf.DefaultExpireMax)*time.Millisecond),
			zcache.WithDirectReturnOnCacheFault(conf.DirectReturnOnCacheFault),
			zcache.WithPanicOnLoaderExists(conf.PanicOnLoaderExists),
			zcache.WithCodec(makeCodec(conf.Codec)),
			zcache.WithSingleFlight(makeSingleFlight(conf.SingleFlight)),
			zcache.WithLogger(app.GetLogger()),
		)
		caches[name] = cache
	}
	return &Cache{caches: caches}
}

// 构建编解码器
func makeCodec(codecType string) cache_core.ICodec {
	switch strings.ToLower(codecType) {
	case "", "default":
		return cache_codec.DefaultCodec
	case "byte":
		return cache_codec.Byte
	case "json":
		return cache_codec.Json
	case "jsoniter":
		return cache_codec.JsonIterator
	case "msgpack":
		return cache_codec.MsgPack
	case "proto_buffer":
		return cache_codec.ProtoBuffer
	}
	logger.Log.Fatal("不支持的Codec", zap.String("Codec", codecType))
	return nil
}

// 构建单跑模块
func makeSingleFlight(sf string) cache_core.ISingleFlight {
	switch strings.ToLower(sf) {
	case "", "default", "single":
		return single_sf.NewSingleFlight()
	case "no":
		return no_sf.NoSingleFlight()
	}
	logger.Log.Fatal("不支持的SingleFlight", zap.String("SingleFlight", sf))
	return nil
}

// 构建缓存db
func makeCacheDB(conf *core.CacheConfig) cache_core.ICacheDB {
	switch strings.ToLower(conf.CacheDB) {
	case "", "default", "memory":
		return makeMemoryCacheDB(conf)
	case "redis":
		return makeRedisCacheDB(conf)
	case "no":
		return no_cache.NoCache()
	}
	logger.Log.Fatal("不支持的CacheDB", zap.String("CacheDB", conf.CacheDB))
	return nil
}

// 构建内存缓存db
func makeMemoryCacheDB(conf *core.CacheConfig) cache_core.ICacheDB {
	return memory_cache.NewMemoryCache(
		memory_cache.WithCleanupInterval(time.Duration(conf.MemoryCacheDB.CleanupInterval) * time.Millisecond),
	)
}

// 构建redis缓存db
func makeRedisCacheDB(conf *core.CacheConfig) cache_core.ICacheDB {
	if conf.RedisCacheDB.Address == "" {
		logger.Log.Panic("redis的address为空")
	}

	var client redis.UniversalClient
	if conf.RedisCacheDB.IsCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        strings.Split(conf.RedisCacheDB.Address, ","),
			Password:     conf.RedisCacheDB.Password,
			PoolSize:     conf.RedisCacheDB.PoolSize,
			ReadTimeout:  time.Duration(conf.RedisCacheDB.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.RedisCacheDB.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.RedisCacheDB.DialTimeout) * time.Millisecond,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:         conf.RedisCacheDB.Address,
			Password:     conf.RedisCacheDB.Password,
			DB:           conf.RedisCacheDB.DB,
			PoolSize:     conf.RedisCacheDB.PoolSize,
			ReadTimeout:  time.Duration(conf.RedisCacheDB.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.RedisCacheDB.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.RedisCacheDB.DialTimeout) * time.Millisecond,
		})
	}

	return redis_cache.NewRedisCache(client,
		redis_cache.WithKeyPrefix(conf.RedisCacheDB.KeyPrefix),
	)
}
