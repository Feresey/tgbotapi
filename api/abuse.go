package api

type truebool struct{}

func (truebool) MarshalText() ([]byte, error) { return []byte("true"), nil }
func (*truebool) UnmarshalText([]byte) error  { return nil }
