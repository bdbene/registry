package storage

import "fmt"

// StorageError is used to indicate an occurance of an error during
// the storage of data.
type StorageError struct {
	cause string
}

func (err *StorageError) Error() string {
	return fmt.Sprintf("Storage error: %s", err.cause)
}
