package profile

type Profile interface {
	Hash() string
	GetFields() map[string][]string
}
