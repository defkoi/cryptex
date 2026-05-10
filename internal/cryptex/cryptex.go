package cryptex

const Version = 0

type Cryptex struct {
	data map[string]string
}

func New() *Cryptex {
	return &Cryptex{
		data: make(map[string]string),
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
