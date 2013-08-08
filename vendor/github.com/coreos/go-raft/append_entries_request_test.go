package raft

import (
	"bytes"
	"testing"
)

func BenchmarkAppendEntriesRequestEncoding(b *testing.B) {
	req, tmp := createTestAppendEntriesRequest(2000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		req.encode(&buf)
	}
	b.SetBytes(int64(len(tmp)))
}

func BenchmarkAppendEntriesRequestDecoding(b *testing.B) {
	req, buf := createTestAppendEntriesRequest(2000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.decode(bytes.NewReader(buf))
	}
	b.SetBytes(int64(len(buf)))
}

func createTestAppendEntriesRequest(entryCount int) (*AppendEntriesRequest, []byte) {
	entries := make([]*LogEntry, 0)
	for i := 0; i < entryCount; i++ {
		command := &DefaultJoinCommand{Name: "localhost:1000"}
		entry, _ := newLogEntry(nil, 1, 2, command)
		entries = append(entries, entry)
	}
	req := newAppendEntriesRequest(1, 1, 1, 1, "leader", entries)

	var buf bytes.Buffer
	req.encode(&buf)

	return req, buf.Bytes()
}
