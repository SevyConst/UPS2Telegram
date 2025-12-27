package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"UPS2Telegram/internal/config"
	"UPS2Telegram/internal/telegram"
)
const timerMinutes int8 = 15

func main() {
	timestamp := time.Now().Format("02.01.2006 15:04:05")

	env := flag.String("env", "", "Environment")
	action := flag.String("action", "", "From UPS")

	flag.Parse()
	if *env == "" {
		log.Fatalf("Input parameter -env is neсessary")
	}
	if *env != "local" && *env != "prod" {
		log.Fatalf("invalid env: '%s'", *env)
	}
	
	cfg, err := config.Load(*env)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var msg string
	switch *action {
	case "earlyshutdown":
		msg = fmt.Sprintf("⚠️ ВЫКЛЮЧЕНИЕ RASPBERRY PI. На питании от аккумулятора уже больше %d минут. %s",
			timerMinutes,
			timestamp,
		)
	case "onbatt":
		msg = fmt.Sprintf("🔋 ПИТАНИЕ ОТ АККУМУЛЯТОРА ИБП. Таймер выключения запущен (%d мин). %s",
			timerMinutes,
			timestamp,
		)
	case "online":
		msg = fmt.Sprintf("✅ ПИТАНИЕ ОТ СЕТИ ВОССТАНОВЛЕНО. Таймер отменен. %s", timestamp)
	default:
		log.Fatalf("Unknown action: '%s'", *action)
	}

	if err := telegram.SendToMultipleChats(cfg.Telegram.Token, cfg.Telegram.ChatIDs, msg); err != nil {
		log.Fatalf("failed to send to Telegram: %v", err) 
	}
	
}