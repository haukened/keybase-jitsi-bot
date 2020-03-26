package main

// writeKV is an internal func that ensures KVStore values get written consistently
func (b *bot) writeKV(key string, value string) error {
	_, err := b.k.KVPut(&b.config.KVStoreTeam, b.k.Username, key, value)
	if err != nil {
		return err
	}
	return nil
}

// getGV is an internal function that ensures KVStore values are retreived consistently
func (b *bot) getKV(key string) (value string, revision int, err error) {
	res, err := b.k.KVGet(&b.config.KVStoreTeam, b.k.Username, key)
	return res.EntryValue, res.Revision, err
}
