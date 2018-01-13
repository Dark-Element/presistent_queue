package adapters

type QueueInterface interface {
	Push(data []byte)
	Pop(out chan []byte, n int64, s int64)
	Peek() (int64, int64)
	CanPush(s int64, atomic bool) bool
	Close()
}