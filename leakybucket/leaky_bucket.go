package leakybucket

type LeakyBucket struct {
	capacity uint64
	queue    chan []byte
	rate     uint32
	storage  [][]byte
}
