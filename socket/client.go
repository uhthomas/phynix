package pheonix

import "net/http"

type Client struct {
	s *Server
}

func NewClient(res http.ResponseWriter, req *http.Request, s *Server) (*Client, err) {

}

func (c *Client) Join(room string) error {
	return c.s.joinRoom(c, room)
]

func (C *Client) Leave(room string) error {
	return c.s.leaveRoom(c, room)
}
