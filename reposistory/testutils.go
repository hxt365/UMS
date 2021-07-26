package reposistory

import "github.com/go-redis/redis/v8"

type RedisMock struct {
	store  int
	load   int
	delete int
}

func (r *RedisMock) StoreJson(key string, v interface{}) error {
	r.store++
	return nil
}

func (r *RedisMock) CalledStore() bool {
	return r.store > 0
}

func (r *RedisMock) CalledStoreTwice() bool {
	return r.store == 2
}

func (r *RedisMock) LoadJson(key string, v interface{}) error {
	r.load++
	if r.store == 0 {
		return redis.Nil
	}
	return nil
}

func (r *RedisMock) CalledLoad() bool {
	return r.load > 0
}

func (r *RedisMock) CalledLoadTwice() bool {
	return r.load == 2
}

func (r *RedisMock) Delete(key string) error {
	r.delete++
	return nil
}

func (r *RedisMock) CalledDelete() bool {
	return r.delete > 0
}

func (r *RedisMock) FlushAll() error {
	return nil
}
