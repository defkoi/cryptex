package cryptex

import (
	"fmt"
	"iter"
	"maps"
)

type Cryptex struct {
	data map[string]string

	iter uint32
	mode Mode
}

func New(iter uint, mode Mode) (*Cryptex, error) {
	if iter > MaxIter {
		return nil, fmt.Errorf("too much iterations (max %d)", MaxIter)
	}

	return &Cryptex{
		data: make(map[string]string),
		iter: uint32(iter),
		mode: ModeGCM,
	}, nil
}

func (c *Cryptex) Store(k, v string) {
	if v == "" {
		delete(c.data, k)
	} else {
		c.data[k] = v
	}
}

func (c *Cryptex) Load(k string) string {
	return c.data[k]
}

func (c *Cryptex) Has(k string) bool {
	_, ok := c.data[k]
	return ok
}

func (c *Cryptex) Keys() iter.Seq[string] {
	return maps.Keys(c.data)
}
