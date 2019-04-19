package http

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb"
)

func logKey(r influxdb.ResourceType, id influxdb.ID) []byte {
	return []byte(fmt.Sprintf("%s_%s_auditlog", r, id))
}

func addLog(ctx context.Context, svc influxdb.KeyValueLog, r influxdb.ResourceType, id influxdb.ID, desc string, now time.Time) error {
	return svc.AddLogEntry(ctx, logKey(r, id), []byte(desc), now)
}

func timeToStr(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func getLogCreatedUpdated(ctx context.Context, svc influxdb.KeyValueLog, id influxdb.ID, r influxdb.ResourceType) (createdAt, updatedAt time.Time, err error) {
	key := logKey(r, id)
	_, createdAt, err = svc.FirstLogEntry(ctx, key)
	if err != nil {
		return createdAt, updatedAt, err
	}
	_, updatedAt, err = svc.LastLogEntry(ctx, key)
	if err != nil {
		return createdAt, updatedAt, err
	}
	return createdAt, updatedAt, err
}

// timeGenerator can be easily embed to any backend service,
// to call Now() to return the real time or Set UseFake to true,
// to use fakeValue.
type timeGenerator struct {
	UseFake   bool
	FakeValue time.Time
}

func (g *timeGenerator) Now() time.Time {
	if g.UseFake {
		return g.FakeValue
	}
	return time.Now()
}
