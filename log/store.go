package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

/* 레코드 길이를 저장하는 바이트 개수를 정의 */
const (
	lenWidth = 8
)

var (
	enc = binary.BigEndian
)

/* 파일의 단순한 래퍼(wrapper) */
type store struct {
	mu sync.Mutex

	*os.File
	size uint64
	buf  *bufio.Writer
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())

	if err != nil {
		return nil, err
	}

	size := uint64(fi.Size())

	return &store{
		File: f,
		buf:  bufio.NewWriter(f),
		size: size,
	}, nil
}

// Append
// 저장 파일에 직접 쓰지는 않고, 버퍼를 거쳐 저장하는데, 시스템 호출 회수를 줄여 성능을 개선해준다.
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos = s.size

	/* store buffer에 p 길이만큼의 공간확보 */
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}

	w, err := s.buf.Write(p)

	if err != nil {
		return 0, 0, err
	}

	w += lenWidth
	s.size += uint64(w)

	return uint64(w), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	/* 혹시 모를 버퍼에 있을 데이터를 flush */
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}

	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(p, off)
}
