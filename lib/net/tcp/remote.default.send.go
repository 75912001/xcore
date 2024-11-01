package tcp

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

func Send(remote IRemote, pb proto.Message, messageID uint32, sessionID uint32, key uint64) error {
	if err := remote.Send(
		&packet.Packet{
			Header: &packet.Header{
				MessageID: messageID,
				SessionID: sessionID,
				Key:       key,
			},
			PBMessage: pb,
		},
	); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

func SendError(remote IRemote, messageID uint32, sessionID uint32, resultID uint32, key uint64) error {
	if err := remote.Send(
		&packet.Packet{
			Header: &packet.Header{
				MessageID: messageID,
				SessionID: sessionID,
				ResultID:  resultID,
				Key:       key,
			},
		},
	); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}
