package network

import "time"

type ServerOption func(s *Server)

func WithConnEvent(event IConnEvent) ServerOption {
	return func(s *Server) {
		s.connEvent = event
	}
}

type ClientOption func(c *Client)

func WithClientTimeOut(time time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = time
	}
}


