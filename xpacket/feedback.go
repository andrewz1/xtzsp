package xpacket

type fback bool

func (f *fback) SetTruncated() {
	*f = true
}

func (f fback) isTruncated() bool {
	return bool(f)
}
