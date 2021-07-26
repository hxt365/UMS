package reposistory

type Cache interface {
	StoreJson(key string, v interface{}) error
	LoadJson(key string, v interface{}) error
	Delete(key string) error
	FlushAll() error
}
