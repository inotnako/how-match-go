package source

type Sourcer interface {
	Get(path string) ([]byte, error)
}
