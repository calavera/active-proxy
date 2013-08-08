package raft

// Join command interface
type JoinCommand interface {
	CommandName() string
	Apply(server *Server) (interface{}, error)
	NodeName() string
}

// Join command
type DefaultJoinCommand struct {
	Name string `json:"name"`
}

// The name of the Join command in the log
func (c *DefaultJoinCommand) CommandName() string {
	return "raft:join"
}

func (c *DefaultJoinCommand) Apply(server *Server) (interface{}, error) {
	err := server.AddPeer(c.Name)

	return []byte("join"), err
}

func (c *DefaultJoinCommand) NodeName() string {
	return c.Name
}
