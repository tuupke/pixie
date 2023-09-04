package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractInt(t *testing.T) {
	body := "ipp://localhost:631/jobs/795!\u0000\u0006job-id\u0000\u0004\u0000\u0000\u0003\u001B# \tjob-state \u0004   \u0004A \u0011job-state-message"
	val, found := extractInt("job-id", []byte(body))
	require.True(t, found)
	require.EqualValues(t, 795, val)

	body = "\u0001\u0001\u0003\u001B"
	fmt.Println(int32(binary.BigEndian.Uint32([]byte(body))))
}

func TestLoadTwice(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	if a, b := Load(ip, nil), Load(ip, nil); a != b {
		t.Fail()
	}
}
