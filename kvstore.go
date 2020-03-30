package main

import (
	"reflect"

	"samhofi.us/x/keybase/types/chat1"
)

// mashals an interface to JSON and sends to kvstore
func (b *bot) KVStorePutStruct(convIDstr chat1.ConvIDStr, v interface{}) error {
	// marshal the struct to JSON
	kvstoreDataString, err := encodeStructToJSONString(v)
	if err != nil {
		return err
	}
	// put the string in kvstore
	err = b.KVStorePut(string(convIDstr), getTypeName(v), kvstoreDataString)
	if err != nil {
		return err
	}
	return nil
}

func (b *bot) KVStoreGetStruct(convIDstr chat1.ConvIDStr, v interface{}) error {
	// get the string from kvstore
	result, err := b.KVStoreGet(string(convIDstr), getTypeName(v))
	if err != nil {
		return err
	}
	// if there was no result just return and the struct is unmodified
	if result == "" {
		return nil
	}
	// unmarshal the string into JSON
	err = decodeJSONStringToStruct(v, result)
	if err != nil {
		return err
	}
	return nil
}

func (b *bot) KVStorePut(namespace string, key string, value string) error {
	_, err := b.k.KVPut(&b.config.KVStoreTeam, namespace, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (b *bot) KVStoreGet(namespace string, key string) (string, error) {
	kvResult, err := b.k.KVGet(&b.config.KVStoreTeam, namespace, key)
	if err != nil {
		return "", err
	}
	return kvResult.EntryValue, nil
}

// getTypeName returns the name of a type, regardless of if its a pointer or not
func getTypeName(v interface{}) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
