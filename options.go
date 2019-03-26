package server

// Option allows setting variadic options on server.
type Option func(*Server) error

// SetVersion does just that
func SetVersion(v string) Option {
	return func(s *Server) error {
		if v != "" {
			s.meta.Version = v
		}
		return nil
	}
}

// SetDescription sets the description
// Not that useful, just another example for you all
func SetDescription(d string) Option {
	return func(s *Server) error {
		s.meta.Description = d
		return nil
	}
}
