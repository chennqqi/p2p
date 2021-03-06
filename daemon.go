package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/ccding/go-stun/stun"
	ptp "github.com/subutai-io/p2p/lib"
)

// ExecDaemon starts P2P daemon
func ExecDaemon(port int, sFile, profiling, syslog string) {
	if syslog != "" {
		ptp.SetSyslogSocket(syslog)
	}
	StartProfiling(profiling)
	go ptp.InitPlatform()
	ptp.InitErrors()

	if !ptp.CheckPermissions() {
		os.Exit(1)
	}
	StartTime = time.Now()

	ReadyToServe = false
	proc := new(Daemon)
	proc.Initialize(sFile)
	setupRESTHandlers(port, proc)

	ptp.Log(ptp.Info, "Determining outbound IP")
	nat, host, err := stun.NewClient().Discover()
	if err != nil {
		ptp.Log(ptp.Error, "Failed to discover outbound IP: %s", err)
		OutboundIP = nil
	} else {
		OutboundIP = net.ParseIP(host.IP())
		ptp.Log(ptp.Info, "Public IP is %s. %s", OutboundIP.String(), nat)
	}

	if sFile != "" {
		ptp.Log(ptp.Info, "Restore file provided")
		// Try to restore from provided file
		instances, err := proc.Instances.LoadInstances(proc.SaveFile)
		if err != nil {
			ptp.Log(ptp.Error, "Failed to load instances: %v", err)
		} else {
			ptp.Log(ptp.Info, "%d instances were loaded from file", len(instances))
			for _, inst := range instances {
				proc.run(&inst, new(Response))
			}
		}
	}

	ReadyToServe = true

	SignalChannel = make(chan os.Signal, 1)
	signal.Notify(SignalChannel, os.Interrupt)

	go func() {
		for sig := range SignalChannel {
			fmt.Println("Received signal: ", sig)
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
	select {}
}
