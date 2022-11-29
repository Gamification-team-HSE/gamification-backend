package cache

import (
	gocache "github.com/patrickmn/go-cache"
	"time"
)

type Client struct {
	*gocache.Cache
}

func New(expireAt time.Duration, purgeAt time.Duration) *Client {
	cache := gocache.New(expireAt, purgeAt)
	return &Client{Cache: cache}
}
