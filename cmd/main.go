package main

import (
	"os/exec"
	"flag"
	"fmt"
	"log"
	"time"

	"UPS2Telegram/internal/config"
	"UPS2Telegram/internal/telegram"
)


func main() {
	timestamp := time.Now().Format("02.01.2006 15:04:05")

	action := flag.String("action", "", "From UPS")
	configFile := flag.String("config", "", "Path to config directory")
	flag.Parse()
	
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("%v", err)
	}

	switch *action {
	case "earlyshutdown":
		cmd := exec.Command("sudo", "shutdown", "-h", "now")
		if err := cmd.Run(); err != nil {
        	log.Fatalf("Can't turn off raspberry pi: %v", err)
    	}

	case "onbatt":
		msg := fmt.Sprintf("⚠️ ЭЛЕКТРИЧЕСТВО ОТКЛЮЧИЛИ. %s", timestamp)
		telegram.SendToMultipleChats(cfg.Telegram.Token, cfg.Telegram.ChatIDs, msg)

	case "online":
		msg := fmt.Sprintf("✅ ЭЛЕКТРИЧЕСТВО ВКЛЮЧИЛИ. %s", timestamp)
		telegram.SendToMultipleChats(cfg.Telegram.Token, cfg.Telegram.ChatIDs, msg)


	default:
		log.Fatalf("Unknown action: '%s'", *action)
	}
	
}