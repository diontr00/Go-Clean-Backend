package mongo

import "fmt"

type ErrClientDisconnected struct {
	err error
}

func (e ErrClientDisconnected) Error() string {
	return fmt.Sprintf("[Client-Disconnect], %v", e.err)
}

type ErrEmptyParams struct {
	err error
}

func (e ErrEmptyParams) Error() string {
	return fmt.Sprintf("[Empty-Param] , %v", e.err)
}

type ErrTimeout struct {
	err error
}

func (e ErrTimeout) Error() string {
	return fmt.Sprintf("[Timeout], %v", e.err)
}

type ErrNoDocuments struct {
	err error
}

func (e ErrNoDocuments) Error() string {
	return fmt.Sprintf("[No-Document], %v", e.err)
}

type ErrKeyExist struct {
	err error
}

func (e ErrKeyExist) Error() string {
	return fmt.Sprintf("[Key-Exist], %v", e.err)
}

type ErrInternal struct {
	err error
}

func (e ErrInternal) Error() string {
	return fmt.Sprintf("[Internal-Error], %v", e.err)
}
