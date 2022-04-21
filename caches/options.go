package caches

type Options struct {
	//MaxEntrySize a threshold for the total nums of key value pairs
	MaxEntrySize int64
	//MaxGcCount Garbage Collection threshold, when the key pairs collected reach the threshold will stop
	MaxGcCount int
	// GcDuration, the interval of garbage collection
	GcDuration int64
}

func DefaultOptions() Options {
	return Options{
		MaxEntrySize: int64(4),
		MaxGcCount:   1000,
		GcDuration:   60,
	}
}
