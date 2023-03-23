package format

type Formatter interface {
	Match([]byte) error

	Format(args ...interface{}) ([]byte, error)
}
