package utils

type WriterFunc func(p []byte) (n int, err error)

func (w WriterFunc) Write(p []byte) (n int, err error) {
	return w(p)
}
