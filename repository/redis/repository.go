package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/nikas-lebedenko/urlshortener/shortener"
	"github.com/pkg/errors"
)

type repository struct {
	client *redis.Client
}

func newClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *repository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *repository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *repository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

func NewRepository(url string) (shortener.RedirectRepository, error) {
	r := &repository{}
	client, err := newClient(url)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRepository")
	}
	r.client = client
	return r, nil
}
