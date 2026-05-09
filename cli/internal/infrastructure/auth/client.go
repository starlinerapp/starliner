package auth

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Login() error {
	return nil
}
