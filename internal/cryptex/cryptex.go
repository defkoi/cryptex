package cryptex

import (
	"fmt"
	"iter"
	"maps"
)

type Cryptex struct {
	data map[string]string
	iter uint32
}

func New(iter uint) (*Cryptex, error) {
	c := &Cryptex{data: make(map[string]string)}

	if err := c.SetIter(iter); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cryptex) SetIter(iter uint) error {
	if iter > MaxIter {
		return fmt.Errorf("too much iterations (max %d)", MaxIter)
	}

	c.iter = uint32(iter)

	return nil
}

func (c *Cryptex) Iter() uint {
	return uint(c.iter)
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
