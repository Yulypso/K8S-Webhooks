package external

import (
	"io/ioutil"
	"sync"
)

func SyncReadFile(filename string) ([]byte, error) {
	var mutex sync.Mutex
	mutex.Lock()
	configBytes, err := ioutil.ReadFile(filename)
	mutex.Unlock()
	if err != nil {
		return nil, err
	}
	return configBytes, nil
}

func SyncWriteFile(filename string, bytes []byte) error {
	var mutex sync.Mutex
	mutex.Lock()
	err := ioutil.WriteFile(filename, bytes, 0644)
	mutex.Unlock()

	if err != nil {
		return err
	}
	return nil
}
