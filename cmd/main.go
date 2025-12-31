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

var validEnvs = map[string]struct{}{
    "local": {},
    "test":  {},
    "prod":  {},
}

func main() {
	timestamp := time.Now().Format("02.01.2006 15:04:05")

	env := flag.String("env", "", "Environment")
	action := flag.String("action", "", "From UPS")

	flag.Parse()
	if *env == "" {
		log.Fatal("Input parameter -env is necessary")
	}
	if _, exists := validEnvs[*env]; !exists {
		log.Fatalf("invalid env: '%s'", *env)
	}
	
	cfg, err := config.Load(*env)
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
		if err := telegram.SendToMultipleChats(cfg.Telegram.Token, cfg.Telegram.ChatIDs, msg); err != nil {
			log.Fatalf("failed to send to Telegram: %v", err) 
		}

	case "online":
		msg := fmt.Sprintf("✅ ЭЛЕКТРИЧЕСТВО ВКЛЮЧИЛИ. %s", timestamp)
		if err := telegram.SendToMultipleChats(cfg.Telegram.Token, cfg.Telegram.ChatIDs, msg); err != nil {
			log.Fatalf("failed to send to Telegram: %v", err) 
		}

	default:
		log.Fatalf("Unknown action: '%s'", *action)
	}
	
}