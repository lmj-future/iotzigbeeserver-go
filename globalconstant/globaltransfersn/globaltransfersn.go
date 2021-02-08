package globaltransfersn

import "sync"

// TransferSN TransferSN
var TransferSN = struct {
	sync.RWMutex
	SN map[string]int
}{SN: make(map[string]int)}
