package cryptex

const Version = 0

type Cryptex struct {
	data map[string]string

	/*! avoid type conversion (uint -> uint32 -> int) errors !*/
	iter uint
}

func New(iter uint) *Cryptex {
	return &Cryptex{
		data: make(map[string]string),
		iter: iter,
	}
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
