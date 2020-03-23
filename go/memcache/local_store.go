package mc

import "context"

type MapStore struct {
	store map[string]string
}

var LocalStore = &MapStore{
	store: make(map[string]string, 1000),
}

func (ms *MapStore) GetWithErr(ctx context.Context, key string) (string, error) {
	return ms.store[key], nil
}

func (ms *MapStore) SetWithErr(ctx context.Context, key, value string, expire int) error {
	ms.store[key] = value
	return nil
}

func (ms *MapStore) Delete(ctx context.Context, key string) error {
	delete(ms.store, key)
	return nil
}
