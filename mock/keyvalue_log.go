package mock

import (
	"context"
	"time"

	"github.com/influxdata/influxdb"
)

var _ influxdb.KeyValueLog = &KeyValueLog{}

// KeyValueLog is a mocked influxdb.KeyValueLog
type KeyValueLog struct {
	entries           map[string][]logEntry
	ForEachLogEntryFn func(ctx context.Context, k []byte, opts influxdb.FindOptions, fn func(v []byte, t time.Time) error) error
}

type logEntry struct {
	t time.Time
	b []byte
}

// Empty will clean up the map
func (log *KeyValueLog) Empty() {
	for k := range log.entries {
		delete(log.entries, k)
	}
}

// AddLogEntry adds an entry (v,t) to the log defined for the key k.
func (log *KeyValueLog) AddLogEntry(ctx context.Context, k []byte, v []byte, t time.Time) error {
	if log.entries == nil {
		log.entries = make(map[string][]logEntry)
	}
	key := string(k)
	_, ok := log.entries[key]
	if !ok {
		log.entries[key] = make([]logEntry, 0)
	}
	log.entries[key] = append(log.entries[key], logEntry{
		b: v,
		t: t,
	})
	return nil
}

// FirstLogEntry is used to retrieve the first entry in the log at key k.
func (log *KeyValueLog) FirstLogEntry(ctx context.Context, k []byte) ([]byte, time.Time, error) {
	if log.entries == nil {
		return []byte{}, time.Time{}, &influxdb.Error{Code: influxdb.ENotFound, Msg: "log not found"}
	}
	list, ok := log.entries[string(k)]
	if !ok || len(list) == 0 {
		return []byte{}, time.Time{}, &influxdb.Error{Code: influxdb.ENotFound, Msg: "log not found"}
	}
	return list[0].b, list[0].t, nil
}

// LastLogEntry is used to retrieve the last entry in the log at key k.
func (log *KeyValueLog) LastLogEntry(ctx context.Context, k []byte) ([]byte, time.Time, error) {
	if log.entries == nil {
		return []byte{}, time.Time{}, &influxdb.Error{Code: influxdb.ENotFound, Msg: "log not found"}
	}
	list, ok := log.entries[string(k)]
	if !ok || len(list) == 0 {
		return []byte{}, time.Time{}, &influxdb.Error{Code: influxdb.ENotFound, Msg: "log not found"}
	}
	return list[len(list)-1].b, list[len(list)-1].t, nil
}

// ForEachLogEntry iterates through all the log entries at key k and applies the function fn for each record.
func (log *KeyValueLog) ForEachLogEntry(ctx context.Context, k []byte, opts influxdb.FindOptions, fn func(v []byte, t time.Time) error) error {
	return log.ForEachLogEntryFn(ctx, k, opts, fn)
}
