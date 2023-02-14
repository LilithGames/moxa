package main

import (
	"fmt"
	"log"

	"github.com/LilithGames/moxa/proxy"
	"github.com/LilithGames/moxa/utils"
)

func main() {
	target := "moxa:8001"
	memberSyncSource := proxy.NewMemberSyncClientSource(target)
	memberSyncClient := utils.NewSyncClient(memberSyncSource)
	memberReader := proxy.NewMemberStateReader(memberSyncClient)

	server := proxy.NewServer(8001, "dragonboat://:8001", memberReader)
	if err := server.Run(); err != nil {
		log.Fatalln("[FATAL]", fmt.Errorf("server.Run err: %w", err))
	}
}
