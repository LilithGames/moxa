package main

import (
	"fmt"
	"log"

	"github.com/LilithGames/moxa/proxy"
	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/utils"
)

func main() {
	target := "moxa:8001"
	simpleClient, err := service.NewSimpleClient(target)
	if err != nil {
		log.Fatalln("[FATAL]", fmt.Errorf("NewSimpleClient err: %w", err))
	}
	memberSyncSource := proxy.NewMemberSyncClientSource(simpleClient)
	memberSyncClient := utils.NewSyncClient(memberSyncSource)
	memberReader := proxy.NewMemberStateReader(memberSyncClient)

	server := proxy.NewServer(8001, "dragonboat://:8001", memberReader)
	if err := server.Run(); err != nil {
		log.Fatalln("[FATAL]", fmt.Errorf("server.Run err: %w", err))
	}
}
