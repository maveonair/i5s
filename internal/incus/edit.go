package incus

import (
	"context"

	"github.com/lxc/incus/v7/shared/api"
	clicmd "github.com/lxc/incus/v7/shared/cmd"
	"go.yaml.in/yaml/v4"
)

const instanceConfigEditHelp = `### This is a YAML representation of the instance configuration.
### Any line starting with a '#' will be ignored.`

// EditInstanceConfig opens the instance config in the user's editor and applies the result.
func (s *IncusService) EditInstanceConfig(ctx context.Context, name string) error {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return err
	}

	inst, etag, err := server.GetInstance(name)
	if err != nil {
		return err
	}

	data, err := instanceConfigEditContent(inst)
	if err != nil {
		return err
	}

	content, err := clicmd.TextEditor("", data)
	if err != nil {
		return err
	}

	updated, err := parseInstanceConfigEdit(content)
	if err != nil {
		return err
	}

	op, err := server.UpdateInstance(name, updated, etag)
	if err != nil {
		return err
	}
	return op.WaitContext(ctx)
}

func instanceConfigEditContent(inst *api.Instance) ([]byte, error) {
	copy := *inst
	copy.ExpandedConfig = nil
	copy.ExpandedDevices = nil

	data, err := yaml.Dump(&copy, yaml.V2)
	if err != nil {
		return nil, err
	}

	content := instanceConfigEditHelp + "\n\n" + string(data)
	return []byte(content), nil
}

func parseInstanceConfigEdit(content []byte) (api.InstancePut, error) {
	var updated api.InstancePut
	if err := yaml.Load(content, &updated); err != nil {
		return api.InstancePut{}, err
	}
	return updated, nil
}
