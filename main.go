// file: main.go
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"on-air/configutil"
	"on-air/lifxutil"
	"on-air/schedule"
)

func main() {
	var (
		configPath            = flag.String("config", "config.json", "path to config file")
		credsPath             = flag.String("credentials", "", "path to OAuth client JSON")
		tokenPath             = flag.String("token", "", "path to store OAuth tokens")
		calID                 = flag.String("calendar", "", "calendar ID or 'primary'")
		days                  = flag.Int("days", 0, "how many days ahead to check")
		lifxToken             = flag.String("lifx_token", "", "Lifx API token")
		lifxLightID           = flag.String("lifx_light_id", "", "Lifx Light ID")
		lifxLightLabel        = flag.String("lifx_light_label", "", "Lifx Light Label")
		lifxBusyColor         = flag.String("lifx_busy_color", "", "Lifx Busy Color")
		lifxFreeColor         = flag.String("lifx_free_color", "", "Lifx Free Color")
		reloadIntervalSeconds = flag.Int("reload_interval_seconds", 0, "Reload interval in seconds")
	)
	flag.Parse()

	cfg, err := configutil.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Override config values with flags if provided
	if *credsPath != "" {
		cfg.CredsPath = *credsPath
	}
	if *tokenPath != "" {
		cfg.TokenPath = *tokenPath
	}
	if *calID != "" {
		cfg.CalID = *calID
	}
	if *days != 0 {
		cfg.Days = *days
	}
	if *lifxToken != "" {
		cfg.LifxToken = *lifxToken
	}
	if *lifxLightID != "" {
		cfg.LifxLightID = *lifxLightID
	}
	if *lifxLightLabel != "" {
		cfg.LifxLightLabel = *lifxLightLabel
	}
	if *lifxBusyColor != "" {
		cfg.LifxBusyColor = *lifxBusyColor
	}
	if *lifxFreeColor != "" {
		cfg.LifxFreeColor = *lifxFreeColor
	}
	if *reloadIntervalSeconds != 0 {
		cfg.ReloadIntervalSeconds = *reloadIntervalSeconds
	}

	manager := &schedule.Manager{
		CredsPath:             cfg.CredsPath,
		TokenPath:             cfg.TokenPath,
		CalID:                 cfg.CalID,
		Days:                  cfg.Days,
		LifxToken:             cfg.LifxToken,
		LifxLightID:           cfg.LifxLightID,
		LifxLightLabel:        cfg.LifxLightLabel,
		LifxBusyColor:         cfg.LifxBusyColor,
		LifxFreeColor:         cfg.LifxFreeColor,
		ReloadIntervalSeconds: cfg.ReloadIntervalSeconds,
	}
	manager.Update(manager.LoadSchedule()) // initial load

	actionCh := make(chan schedule.Action, 10) // buffered channel

	go schedule.Reloader(manager)
	go schedule.ActionWorker(actionCh, cfg.LifxToken, cfg.LifxLightID, cfg.LifxLightLabel, cfg.LifxBusyColor, cfg.LifxFreeColor)
	go schedule.Executor(manager, actionCh)

	// Signal handling for graceful shutdown
	// TODO: make this cleaner
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Received signal %v, setting light to free state and exiting...", sig)
		lifxc := lifxutil.NewClient(cfg.LifxToken)
		if err != nil {
			log.Printf("Failed to create Lifx client: %v", err)
			os.Exit(1)
		}
		light := lifxutil.Light{
			ID:    cfg.LifxLightID,
			Label: cfg.LifxLightLabel,
		}
		err = lifxc.SetFree(light, cfg.LifxFreeColor)
		if err != nil {
			log.Printf("Failed to set light to free state: %v", err)
		}
		os.Exit(0)
	}()

	select {} // block forever
}
