package chunk

type Iterator interface {
	Current() string
	Next() bool
	Prev() bool
}
