package main

import (
	"context"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

func router(
	ctx context.Context,
	handlers map[string]func(commoninterfaces.ChannelRequester),
	chNatsIn <-chan commoninterfaces.ChannelRequester) {

	for {
		select {
		case <-ctx.Done():

		case msg := <-chNatsIn:
			if msg.GetCommand() == "send_command" {
				if f, ok := handlers[msg.GetOrder()]; ok {
					go f(msg)
				}
			}
		}
	}
}
