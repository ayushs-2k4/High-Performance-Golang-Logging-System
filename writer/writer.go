package writer

type Writer interface {
	Write(b []byte)
	Close()
}
