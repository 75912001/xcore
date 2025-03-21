package tcp

import (
	xcommon "github.com/75912001/xcore/lib/net/common"
	xpacket "github.com/75912001/xcore/lib/packet"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func Send(remote xcommon.IRemote, pb proto.Message, messageID uint32, sessionID uint32, key uint64) error {
	if err := remote.Send(
		&xpacket.Packet{
			Header: &xpacket.Header{
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

func SendError(remote xcommon.IRemote, messageID uint32, sessionID uint32, resultID uint32, key uint64) error {
	if err := remote.Send(
		&xpacket.Packet{
			Header: &xpacket.Header{
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
