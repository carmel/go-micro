package addrutil_test

import (
	"net"
	"testing"

	"go-micro/internal/addrutil"

	"github.com/stretchr/testify/require"
)

func TestAddrToKey(t *testing.T) {
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10000")
	require.Nil(t, err)
	raddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10001")
	require.Nil(t, err)
	key := addrutil.AddrToKey(laddr, raddr)
	require.Equal(t, key, laddr.Network()+"_"+laddr.String()+"_"+raddr.String())
}
