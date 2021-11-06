package profile

type Fields map[string][]string

type Profile interface {
	GetFields() Fields
}
