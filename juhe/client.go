package juhe

type Client struct {
	key string
}

var (
	Default *Client
)

func SetDefault(key string) {
	Default = NewClient(key)
}

func NewClient(key string) (c *Client) {
	c = &Client{key: key}
	return
}
