package incus

import (
	"context"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/lxc/incus/v7/shared/api"
)

// ListInstances returns instance rows for the selected runtime remote/project.
func (s *IncusService) ListInstances(ctx context.Context) ([]InstanceRow, error) {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return nil, err
	}

	instances, err := server.GetInstances(api.InstanceTypeAny)
	if err != nil {
		return nil, err
	}

	rows := make([]InstanceRow, 0, len(instances))
	for _, inst := range instances {
		rows = append(rows, instanceToRow(inst))
	}

	s.enrichInstanceRows(ctx, server, rows)

	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows, nil
}

func (s *IncusService) enrichInstanceRows(ctx context.Context, server interface {
	GetInstanceSnapshotNames(name string) ([]string, error)
	GetInstanceState(name string) (*api.InstanceState, string, error)
}, rows []InstanceRow) {
	sem := make(chan struct{}, instanceEnrichmentConcurrency)
	var wg sync.WaitGroup

launch:
	for i := range rows {
		select {
		case <-ctx.Done():
			break launch
		default:
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return
			}

			name := rows[i].Name
			if snapshots, err := server.GetInstanceSnapshotNames(name); err == nil {
				rows[i].Snapshots = len(snapshots)
			} else {
				log.Printf("failed to enrich snapshots for %s: %v", name, err)
			}

			if rows[i].StatusCode == api.Running {
				if state, _, err := server.GetInstanceState(name); err == nil {
					rows[i].IPv4, rows[i].IPv6 = instanceIPs(state)
				} else {
					log.Printf("failed to enrich state for %s: %v", name, err)
				}
			}
		}(i)
	}

	wg.Wait()
}

// GetInstanceState returns runtime state for one instance.
func (s *IncusService) GetInstanceState(ctx context.Context, name string) (*api.InstanceState, error) {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return nil, err
	}

	state, _, err := server.GetInstanceState(name)
	return state, err
}

// StartInstance starts an instance and waits for the Incus operation to finish.
func (s *IncusService) StartInstance(ctx context.Context, name string) error {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return err
	}

	op, err := server.UpdateInstanceState(name, api.InstanceStatePut{Action: "start", Timeout: -1}, "")
	if err != nil {
		return err
	}
	return op.WaitContext(ctx)
}

// StopInstance stops an instance and waits for the Incus operation to finish.
func (s *IncusService) StopInstance(ctx context.Context, name string) error {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return err
	}

	op, err := server.UpdateInstanceState(name, api.InstanceStatePut{Action: "stop", Timeout: 30, Force: false}, "")
	if err != nil {
		return err
	}
	return op.WaitContext(ctx)
}

// DeleteInstance deletes an instance and waits for the Incus operation to finish.
func (s *IncusService) DeleteInstance(ctx context.Context, name string) error {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return err
	}

	op, err := server.DeleteInstance(name)
	if err != nil {
		return err
	}
	return op.WaitContext(ctx)
}

func instanceToRow(inst api.Instance) InstanceRow {
	image := firstNonEmpty(inst.Config["image.description"], inst.Config["volatile.base_image"], inst.Description)
	return InstanceRow{
		Name:        inst.Name,
		Type:        inst.Type,
		Status:      strings.ToUpper(inst.Status),
		StatusCode:  inst.StatusCode,
		IPv4:        "-",
		IPv6:        "-",
		Image:       image,
		Description: inst.Description,
		Profiles:    inst.Profiles,
		Location:    dash(inst.Location),
		CreatedAt:   inst.CreatedAt,
		Raw:         inst,
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "-"
}

func dash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func instanceIPs(state *api.InstanceState) (string, string) {
	if state == nil {
		return "-", "-"
	}

	var ipv4 []string
	var ipv6 []string
	for name, network := range state.Network {
		if name == "lo" || network.Type == "loopback" {
			continue
		}

		for _, address := range network.Addresses {
			if address.Scope != "global" {
				continue
			}

			switch address.Family {
			case "inet":
				ipv4 = append(ipv4, address.Address)
			case "inet6":
				ipv6 = append(ipv6, address.Address)
			}
		}
	}
	sort.Strings(ipv4)
	sort.Strings(ipv6)

	return joinIPs(ipv4), joinIPs(ipv6)
}

func joinIPs(ips []string) string {
	if len(ips) == 0 {
		return "-"
	}
	return strings.Join(ips, ",")
}
