package promtail

import (
	"github.com/cortexproject/cortex/pkg/util"

	"promtail/pkg/client"
	"promtail/pkg/config"
	"promtail/pkg/positions"
	"promtail/pkg/server"
	"promtail/pkg/targets"
)

// Promtail is the root struct for Promtail...
type Promtail struct {
	client         client.Client
	positions      *positions.Positions
	targetManagers *targets.TargetManagers
	server         *server.Server
}

// New makes a new Promtail.
func New(cfg config.Config) (*Promtail, error) {
	positions, err := positions.New(cfg.PositionsConfig)
	if err != nil {
		return nil, err
	}

	if cfg.ClientConfig.URL.URL != nil {
		// if a single client config is used we add it to the multiple client config for backward compatibility
		cfg.ClientConfigs = append(cfg.ClientConfigs, cfg.ClientConfig)
	}

	client, err := client.NewMulti(util.Logger, cfg.ClientConfigs...)
	if err != nil {
		return nil, err
	}

	tms, err := targets.NewTargetManagers(util.Logger, positions, client, cfg.ScrapeConfig, &cfg.TargetConfig)
	if err != nil {
		return nil, err
	}

	server, err := server.New(cfg.ServerConfig, tms)
	if err != nil {
		return nil, err
	}

	return &Promtail{
		client:         client,
		positions:      positions,
		targetManagers: tms,
		server:         server,
	}, nil
}

// Run the promtail; will block until a signal is received.
func (p *Promtail) Run() error {
	return p.server.Run()
}

// Shutdown the promtail.
func (p *Promtail) Shutdown() {
	p.server.Shutdown()
	p.targetManagers.Stop()
	p.positions.Stop()
	p.client.Stop()
}
