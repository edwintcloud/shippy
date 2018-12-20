package main

// Authable is the jwt token interface
type Authable interface {
	Decode(token string) (interface{}, error)
	Encode(data interface{}) (string, error)
}

// TokenService is the struct that holds our repo
type TokenService struct {
	repo Repository
}

// Decode decodes a jwt token
func (s *TokenService) Decode(token string) (interface{}, error) {
	// TODO: implement
	return "", nil
}

// Encode encodes data to a jwt token
func (s *TokenService) Encode(data interface{}) (string, error) {
	// TODO: implement
	return "", nil
}
