package sender

import "fmt"

type ShortenerSenderRepositoryMock struct {
}

func NewShortenerSenderRepositoryMock() *ShortenerSenderRepositoryMock {
	return &ShortenerSenderRepositoryMock{}
}

func (s *ShortenerSenderRepositoryMock) SendRedirectInfo() {
	fmt.Println("redirect info!")
}
