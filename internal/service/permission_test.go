package service

import (
	"testing"

	"github.com/WuKongIM/WuKongIM/pkg/wkdb"
	wkproto "github.com/WuKongIM/WuKongIMGoProto"
)

func TestExceedsMaxSendWithoutReply(t *testing.T) {
	tests := []struct {
		name     string
		messages []wkdb.Message
		fromUID  string
		limit    int
		want     bool
	}{
		{
			name:     "fewer than limit",
			messages: makeMessages("u1", "u1", "u1"),
			fromUID:  "u1",
			limit:    10,
			want:     false,
		},
		{
			name:     "exact limit all from sender",
			messages: makeMessages("u1", "u1", "u1", "u1", "u1", "u1", "u1", "u1", "u1", "u1"),
			fromUID:  "u1",
			limit:    10,
			want:     true,
		},
		{
			name:     "mixed senders",
			messages: makeMessages("u1", "u1", "u1", "u1", "u2", "u1", "u1", "u1", "u1", "u1"),
			fromUID:  "u1",
			limit:    10,
			want:     false,
		},
		{
			name:     "last messages all from other user",
			messages: makeMessages("u2", "u2", "u2", "u2", "u2", "u2", "u2", "u2", "u2", "u2"),
			fromUID:  "u1",
			limit:    10,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exceedsMaxSendWithoutReply(tt.messages, tt.fromUID, tt.limit)
			if got != tt.want {
				t.Fatalf("exceedsMaxSendWithoutReply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeMessages(fromUIDs ...string) []wkdb.Message {
	messages := make([]wkdb.Message, 0, len(fromUIDs))
	for index, fromUID := range fromUIDs {
		messages = append(messages, wkdb.Message{
			RecvPacket: wkproto.RecvPacket{
				FromUID:    fromUID,
				MessageSeq: uint32(index + 1),
			},
		})
	}
	return messages
}
