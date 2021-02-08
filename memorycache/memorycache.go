package memorycache

import (
	"sync"

	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/pkg/errors"
)

//MemoryCache MemoryCache
type MemoryCache struct{}

var memoryListCache = struct {
	sync.RWMutex
	memoryCache map[string]map[int][]byte
}{memoryCache: make(map[string]map[int][]byte)}
var memoryListAllCache = struct {
	sync.RWMutex
	memoryCache map[string][]string
}{memoryCache: make(map[string][]string)}
var memorySingleCache = struct {
	sync.RWMutex
	memoryCache map[string]string
}{memoryCache: make(map[string]string)}

// var memoryListSyncCache = &sync.Map{}
// var memoryListAllSyncCache = &sync.Map{}
// var memorySingleSyncCache = &sync.Map{}

//Interface Interface
type Interface interface {
	InsertMemory(key string, data []byte) (int, error)
	UpdateMemory(key string, data []byte) (string, error)
	DeleteMemory(key string) (int, error)
	DeleteMemoryOne(key string, value string) (int, error)
	GetMemory(key string) (string, error)
	GetMemoryByIndex(key string, index int) (string, error)
	GetMemoryEnd(key string) (string, error)
	GetMemoryLength(key string) (int, error)
	PopMemory(key string) (string, error)
	RangeMemory(key string, start int, stop int) ([]string, error)
	SetMemory(key string, index int, value []byte) (string, error)
	RemoveMemory(key string, count int, value string) (int, error)
	SaddMemory(key string, member string) (int, error)
	SremMemory(key string, member string) (int, error)
	FindAllMemoryKeys(key string) ([]string, error)
	SetMemorySet(key string, value string) (string, error)
	GetMemoryGet(key string) (string, error)
	GetMemorySize() (int, error)
}

//InsertMemory InsertMemory
func (MemoryCache) InsertMemory(key string, data []byte) (int, error) {
	//往队列尾部添加数据
	valueCache := make(map[int][]byte)
	memoryListCache.Lock()
	if len(memoryListCache.memoryCache[key]) > 0 {
		valueCache = memoryListCache.memoryCache[key]
		valueCache[len(memoryListCache.memoryCache[key])+1] = data
	} else {
		valueCache[1] = data
	}
	memoryListCache.memoryCache[key] = valueCache
	lenTemp := len(valueCache)
	memoryListCache.Unlock()
	return lenTemp, nil

	// valueCache := make(map[int][]byte)
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	valueCache = v.(map[int][]byte)
	// 	valueCache[len(valueCache)+1] = data
	// } else {
	// 	valueCache[1] = data
	// }
	// memoryListSyncCache.Store(key, valueCache)
	// return len(valueCache), nil
}

//UpdateMemory UpdateMemory
func (MemoryCache) UpdateMemory(key string, data []byte) (string, error) {
	valueCache := make(map[int][]byte)
	valueCache[1] = data
	memoryListCache.Lock()
	delete(memoryListCache.memoryCache, key)
	memoryListCache.memoryCache[key] = valueCache
	memoryListCache.Unlock()
	return string(data), nil

	// valueCache := make(map[int][]byte)
	// valueCache[1] = data
	// memoryListSyncCache.Delete(key)
	// memoryListSyncCache.Store(key, valueCache)
	// return string(data), nil
}

//DeleteMemory DeleteMemory
func (MemoryCache) DeleteMemory(key string) (int, error) {
	memoryListCache.Lock()
	delete(memoryListCache.memoryCache, key)
	delete(memorySingleCache.memoryCache, key)
	memoryListCache.Unlock()
	return 1, nil

	// memoryListSyncCache.Delete(key)
	// memorySingleSyncCache.Delete(key)
	// return 1, nil
}

//DeleteMemoryOne DeleteMemoryOne
func (MemoryCache) DeleteMemoryOne(key string, value string) (int, error) {
	memoryListCache.Lock()
	valueLen := len(memoryListCache.memoryCache[key])
	delete(memoryListCache.memoryCache, key)
	memoryListCache.Unlock()
	return valueLen, nil

	// var valueLen int
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	valueLen = len(v.(map[int][]byte))
	// }
	// memoryListSyncCache.Delete(key)
	// return valueLen, nil
}

//GetMemory GetMemory
func (MemoryCache) GetMemory(key string) (string, error) {
	memoryListCache.RLock()
	value := memoryListCache.memoryCache[key][1]
	memoryListCache.RUnlock()
	return string(value), nil

	// var value []byte
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	value = v.(map[int][]byte)[1]
	// }
	// return string(value), nil
}

//GetMemoryByIndex GetMemoryByIndex
func (MemoryCache) GetMemoryByIndex(key string, index int) (string, error) {
	memoryListCache.RLock()
	value := memoryListCache.memoryCache[key][index+1]
	memoryListCache.RUnlock()
	return string(value), nil

	// var value []byte
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	value = v.(map[int][]byte)[index+1]
	// }
	// return string(value), nil
}

//GetMemoryEnd GetMemoryEnd
func (MemoryCache) GetMemoryEnd(key string) (string, error) {
	memoryListCache.RLock()
	value := memoryListCache.memoryCache[key][len(memoryListCache.memoryCache[key])]
	memoryListCache.RUnlock()
	return string(value), nil

	// var value []byte
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	value = v.(map[int][]byte)[len(v.(map[int][]byte))]
	// }
	// return string(value), nil
}

//GetMemoryLength GetMemoryLength
func (MemoryCache) GetMemoryLength(key string) (int, error) {
	memoryListCache.RLock()
	lenTemp := len(memoryListCache.memoryCache[key])
	memoryListCache.RUnlock()
	return lenTemp, nil

	// var valueLen int
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	valueLen = len(v.(map[int][]byte))
	// }
	// return valueLen, nil
}

//PopMemory PopMemory
func (MemoryCache) PopMemory(key string) (string, error) {
	memoryListCache.Lock()
	if len(memoryListCache.memoryCache[key]) == 0 {
		return "", nil
	}
	for i := 0; i < len(memoryListCache.memoryCache[key]); i++ {
		memoryListCache.memoryCache[key][i] = memoryListCache.memoryCache[key][i+1]
	}
	value := memoryListCache.memoryCache[key][0]
	delete(memoryListCache.memoryCache[key], 0)
	memoryListCache.Unlock()
	return string(value), nil

	// var value []byte
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	valueCache := make(map[int][]byte)
	// 	for i := 0; i < len(v.(map[int][]byte)); i++ {
	// 		valueCache[i] = v.(map[int][]byte)[i+1]
	// 	}
	// 	value = valueCache[0]
	// 	delete(valueCache, 0)
	// 	memoryListSyncCache.Store(key, valueCache)
	// }
	// return string(value), nil
}

//RangeMemory RangeMemory
func (MemoryCache) RangeMemory(key string, start int, stop int) ([]string, error) {
	var valueList []string = make([]string, stop-start)
	memoryListCache.Lock()
	for i := 0; i < stop; i++ {
		valueList[i] = string(memoryListCache.memoryCache[key][i+1])
	}
	memoryListCache.Unlock()
	return valueList, nil

	// var valueList []string = make([]string, stop-start)
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	for i := 0; i < stop; i++ {
	// 		valueList[i] = string(v.(map[int][]byte)[i+1])
	// 	}
	// }
	// return valueList, nil
}

//SetMemory SetMemory
func (MemoryCache) SetMemory(key string, index int, value []byte) (string, error) {
	memoryListCache.Lock()
	if index > len(memoryListCache.memoryCache[key]) {
		return "", errors.Errorf("index out of bounds exception")
	}
	if len(memoryListCache.memoryCache[key]) == 0 {
		valueCache := make(map[int][]byte)
		valueCache[index+1] = value
		memoryListCache.memoryCache[key] = valueCache
	} else {
		memoryListCache.memoryCache[key][index+1] = value
	}
	memoryListCache.Unlock()
	return "OK", nil

	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	valueCache := make(map[int][]byte)
	// 	valueCache = v.(map[int][]byte)
	// 	if index > len(valueCache) {
	// 		return "", errors.Errorf("index out of bounds exception")
	// 	}
	// 	valueCache[index+1] = value
	// 	memoryListSyncCache.Store(key, valueCache)
	// }
	// return "OK", nil
}

//RemoveMemory RemoveMemory
func (MemoryCache) RemoveMemory(key string, count int, value string) (int, error) {
	var countTemp int = 0
	memoryListCache.Lock()
	for i := 0; i < len(memoryListCache.memoryCache[key]); i++ {
		if string(memoryListCache.memoryCache[key][i]) == value && countTemp < count {
			countTemp++
			delete(memoryListCache.memoryCache[key], i)
		}
	}
	for k := 0; k < count; k++ {
		for j := 1; j < len(memoryListCache.memoryCache[key]); j++ {
			if len(memoryListCache.memoryCache[key][j]) == 0 {
				memoryListCache.memoryCache[key][j] = memoryListCache.memoryCache[key][j+1]
			}
		}
	}
	memoryListCache.Unlock()
	return count, nil

	// var countTemp int = 0
	// if v, ok := memoryListSyncCache.Load(key); ok {
	// 	valueCache := make(map[int][]byte)
	// 	valueCache = v.(map[int][]byte)
	// 	for i := 0; i < len(valueCache); i++ {
	// 		if string(valueCache[i]) == value && countTemp < count {
	// 			countTemp++
	// 			delete(valueCache, i)
	// 		}
	// 	}
	// 	for k := 0; k < count; k++ {
	// 		for j := 1; j < len(valueCache); j++ {
	// 			if len(valueCache[j]) == 0 {
	// 				valueCache[j] = valueCache[j+1]
	// 			}
	// 		}
	// 	}
	// 	memoryListSyncCache.Store(key, valueCache)
	// }
	// return count, nil
}

//SaddMemory SaddMemory
func (MemoryCache) SaddMemory(key string, member string) (int, error) {
	memoryListAllCache.Lock()
	memoryListAllCache.memoryCache[key] = append(memoryListAllCache.memoryCache[key], member)
	lenTemp := len(memoryListAllCache.memoryCache[key])
	memoryListAllCache.Unlock()
	return lenTemp, nil

	// var valueCache []string
	// if v, ok := memoryListAllSyncCache.Load(key); ok {
	// 	valueCache = v.([]string)
	// }
	// valueCache = append(valueCache, member)
	// memoryListAllSyncCache.Store(key, valueCache)
	// return len(valueCache), nil
}

//SremMemory SremMemory
func (MemoryCache) SremMemory(key string, member string) (int, error) {
	memoryListAllCache.Lock()
	for i := 0; i < len(memoryListAllCache.memoryCache[key]); i++ {
		if memoryListAllCache.memoryCache[key][i] == member {
			memoryListAllCache.memoryCache[key][i] = ""
		}
	}
	lenTemp := len(memoryListAllCache.memoryCache[key])
	memoryListAllCache.Unlock()
	return lenTemp, nil

	// var valueCache []string
	// if v, ok := memoryListAllSyncCache.Load(key); ok {
	// 	valueCache = v.([]string)
	// 	for i := 0; i < len(valueCache); i++ {
	// 		if valueCache[i] == member {
	// 			valueCache[i] = ""
	// 		}
	// 	}
	// 	memoryListAllSyncCache.Store(key, valueCache)
	// }
	// return len(valueCache), nil
}

//FindAllMemoryKeys FindAllMemoryKeys
func (MemoryCache) FindAllMemoryKeys(key string) ([]string, error) {
	memoryListAllCache.RLock()
	value := memoryListAllCache.memoryCache[key]
	memoryListAllCache.RUnlock()
	return value, nil

	// var valueCache []string
	// if v, ok := memoryListAllSyncCache.Load(key); ok {
	// 	valueCache = v.([]string)
	// }
	// return valueCache, nil
}

//SetMemorySet SetMemorySet
func (MemoryCache) SetMemorySet(key string, value string) (string, error) {
	memorySingleCache.Lock()
	memorySingleCache.memoryCache[key] = value
	memorySingleCache.Unlock()
	return "OK", nil

	// memorySingleSyncCache.Store(key, value)
	// return "OK", nil
}

//GetMemoryGet GetMemoryGet
func (MemoryCache) GetMemoryGet(key string) (string, error) {
	memorySingleCache.RLock()
	value := memorySingleCache.memoryCache[key]
	memorySingleCache.RUnlock()
	return value, nil

	// var value string
	// if v, ok := memorySingleSyncCache.Load(key); ok {
	// 	value = v.(string)
	// }
	// return value, nil
}

//GetMemorySize  GetMemorySize
func (MemoryCache) GetMemorySize() (int, error) {
	memoryListCache.RLock()
	var sizeTotal int
	for memoryCacheKey, memoryCacheValue := range memoryListCache.memoryCache {
		var size int
		for _, byteValue := range memoryCacheValue {
			size += len(byteValue)
		}
		globallogger.Log.Errorln("[GetMemorySize]: key: ", memoryCacheKey, " size: ", size, " byte")
		globallogger.Log.Errorf("[GetMemorySize]: key: %s memoryCache: %+v", memoryCacheKey, memoryCacheValue)
		sizeTotal += size
	}
	globallogger.Log.Errorln("[GetMemorySize]: memoryListCache total size: ", sizeTotal, " byte")
	memoryListCache.RUnlock()

	memoryListAllCache.RLock()
	var sizeTotal2 int
	for memoryCacheKey, memoryCacheValue := range memoryListAllCache.memoryCache {
		var size int
		for _, stringValue := range memoryCacheValue {
			size += len(stringValue)
		}
		globallogger.Log.Errorln("[GetMemorySize]: key: ", memoryCacheKey, " size: ", size, " byte")
		globallogger.Log.Errorf("[GetMemorySize]: key: %s memoryCache: %+v", memoryCacheKey, memoryCacheValue)
		sizeTotal2 += size
	}
	globallogger.Log.Errorln("[GetMemorySize]: memoryListAllCache total size: ", sizeTotal2, " byte")
	memoryListAllCache.RUnlock()

	memorySingleCache.RLock()
	var sizeTotal3 int
	for memoryCacheKey, memoryCacheValue := range memorySingleCache.memoryCache {
		globallogger.Log.Errorln("[GetMemorySize]: key: ", memoryCacheKey, " size: ", len(memoryCacheValue), " byte")
		globallogger.Log.Errorf("[GetMemorySize]: key: %s memoryCache: %+v", memoryCacheKey, memoryCacheValue)
		sizeTotal3 += len(memoryCacheValue)
	}
	globallogger.Log.Errorln("[GetMemorySize]: memorySingleCache total size: ", sizeTotal3, " byte")
	memorySingleCache.RUnlock()
	return sizeTotal + sizeTotal2 + sizeTotal3, nil
}
