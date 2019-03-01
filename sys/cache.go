package sys

import (
	"github.com/muesli/cache2go"
)

var authCache *cache2go.CacheTable

const userCacheTableKey string = "user_cache"

func UserCache() *cache2go.CacheTable {
	if authCache == nil {
		authCache = cache2go.Cache(userCacheTableKey)
	}
	return authCache
}
