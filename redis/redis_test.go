package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	data := struct {
		A string  `json:"a"`
		B *string `json:"b"`
	}{
		A: "",
		B: nil,
	}
	redis := NewRedis(time.Second * 5)
	if err := redis.Set("this", data); err != nil {
		t.FailNow()
	}
	this, err := redis.Get("this")
	assert.NoError(t, err)
	assert.Equal(t, data.A, this["a"])
	redis.Delete("this")
	_, err = redis.Get("this")
	assert.Error(t, err)
}
