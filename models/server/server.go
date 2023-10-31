package server

type Server struct {
	id    int
	table map[string]string
}

func New(id int) Server {
	return Server{id, make(map[string]string)}
}

func (server Server) Id() int {
	return server.id
}

func (server Server) Map() map[string]string {
	return server.table
}

func (server *Server) Put(key string, val string) {
	server.table[key] = val
}

func (server Server) Get(key string) string {
	return server.table[key]
}

func (server *Server) Del(key string) {
	delete(server.table, key)
}

func (server Server) Keys() []string {
	keys := make([]string, 0)
	for k := range server.table {
		keys = append(keys, k)
	}
	return keys
}
