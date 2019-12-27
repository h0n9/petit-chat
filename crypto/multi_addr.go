package crypto

import (
	"strings"

	maddr "github.com/multiformats/go-multiaddr"
)

type Addrs []maddr.Multiaddr

func (addrs *Addrs) String() string {
	strs := make([]string, len(*addrs))
	for i, addr := range *addrs {
		strs[i] = addr.String()
	}

	return strings.Join(strs, ",")
}

func (addrs *Addrs) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}

	*addrs = append(*addrs, addr)

	return nil
}

func (addrs *Addrs) ToMultiAddr() []maddr.Multiaddr {
	return []maddr.Multiaddr(*addrs)
}

func NewMultiAddr(value string) (maddr.Multiaddr, error) {
	return maddr.NewMultiaddr(value)
}
