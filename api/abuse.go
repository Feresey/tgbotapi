package api

type truebool struct{}

func (t truebool) MarshalText() ([]byte, error) { return []byte("true"), nil }
func (t *truebool) UnmarshalText([]byte) error  { return nil }
