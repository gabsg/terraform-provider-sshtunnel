package tunnel

type Addr struct {
	Proto string
	IP    string
}

func (a Addr) Network() string {
	return a.Proto
}

func (a Addr) String() string {
	return a.IP
}
