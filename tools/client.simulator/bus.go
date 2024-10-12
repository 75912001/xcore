package main

import (
	xnettcp "xcore/lib/net/tcp"
)

func Handle(busChannel chan interface{}) error {
	for {
		select {
		case value := <-busChannel:
			var err error
			switch event := value.(type) {
			//tcp
			case *xnettcp.Connect:
				err = event.IHandler.OnConnect(event.IRemote)
			case *xnettcp.Packet:
				err = event.IHandler.OnPacket(event.IRemote, event.IPacket)
			case *xnettcp.Disconnect:
				err = event.IHandler.OnDisconnect(event.IRemote)
				if event.IRemote.IsConnect() {
					event.IRemote.Stop()
				}
			default:
			}
			if err != nil {
			}

		}
	}
}
