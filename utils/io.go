package utils

// WriterFunc is a function which satisfies the io.Writer interface
type WriterFunc func(p []byte) (n int, err error)

// Write will execute the WriterFunc
func (w WriterFunc) Write(p []byte) (n int, err error) {
	return w(p)
}
