package main

import (
	"context"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

func router(ctx context.Context, chNats <-chan commoninterfaces.ChannelRequester) {
	for {
		select {
		case <-ctx.Done():

		case msg := <-chNats:
		}
	}
}
