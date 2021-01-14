package main

type KVS interface {
	// Store stores a key value pair.
	Store(k, v string) error

	// Load returns the value associated with a given key.
	Load(k string) (string, error)
}

type MyKVS map[string]string

func NewMyKVS() MyKVS {
	return MyKVS{}
}

func (kvs MyKVS) Store(k, v string) error {
	_, ok := kvs[k]
	if ok {
		return MyKVSError{"already exist"}
	}
	kvs[k] = v
	return nil
}

func (kvs MyKVS) Load(k string) (string, error) {
	v, ok := kvs[k]
	if !ok {
		return "", MyKVSError{"not found"}
	}
	return v, nil
}

type MyKVSError struct{ msg string }

func (e MyKVSError) Error() string {
	return e.msg
}
