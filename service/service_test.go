package service

import "sample-project/storage"

type fakeStoragePostgres struct {
	storage.Storager
}

type fakeStoragePostgresErr struct {
	storage.Storager
}

type fakeStorageRedis struct {
	storage.RedisStorager
}

type fakeStorageRedisErr struct {
	storage.RedisStorager
}
