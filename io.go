package util

type ZeroReader struct{}

func (zeroReader ZeroReader) Read(p []byte) (n int, err error) {
	if p == nil {
		panic("ZeroReader.Read: nil input")
	}
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}
