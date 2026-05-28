package cryptex

import "fmt"

const Version = 0

type Cryptex struct {
	data map[string]string

	iter uint32
}

func New(iter uint) (*Cryptex, error) {
	if iter > MaxIter {
		return nil, fmt.Errorf("too much iterations (max %d)", MaxIter)
	}

	return &Cryptex{
		data: make(map[string]string),
		iter: uint32(iter),
	}, nil
}

func (c *Cryptex) Store(k, v string) {
	c.data[k] = v
}

func (c *Cryptex) Load(k string) string {
	return c.data[k]
}

func (c *Cryptex) Has(k string) bool {
	_, ok := c.data[k]
	return ok
}
