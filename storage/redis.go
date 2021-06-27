package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"sample-project/structs"
)

type (
	RedisStorager interface {
		FetchAuth(authD *structs.AccessDetails) error
		RegisterAuth(userUUID string, token *structs.AuthToken) error
		DeleteAuth(givenUuid string) (int64, error)
	}
	TokenStorage struct {
		client *redis.Client
	}
)

func NewRedis(redisDSN, password string) (*TokenStorage, error) {
	if len(redisDSN) == 0 {
		redisDSN = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{Addr: redisDSN, Password: password})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &TokenStorage{client: client}, nil
}

func (s *TokenStorage) FetchAuth(authD *structs.AccessDetails) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := s.client.Get(ctx, authD.AccessUuid).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *TokenStorage) RegisterAuth(userUUID string, token *structs.AuthToken) error {
	atExpires := time.Unix(token.AccessExpiresAt, 0)
	rtExpires := time.Unix(token.RefreshExpiresAt, 0)

	now := time.Now()

	ctxAT, cancelAT := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelAT()
	err := s.client.Set(ctxAT, token.AccessUuid, userUUID, atExpires.Sub(now)).Err()
	if err != nil {
		return fmt.Errorf("access token store: %w", err)
	}

	ctxRT, cancelRT := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelRT()
	err = s.client.Set(ctxRT, token.RefreshUuid, userUUID, rtExpires.Sub(now)).Err()
	if err != nil {
		return fmt.Errorf("refresh token store: %w", err)
	}

	return err
}

func (s *TokenStorage) DeleteAuth(givenUuid string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	deleted, err := s.client.Del(ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
