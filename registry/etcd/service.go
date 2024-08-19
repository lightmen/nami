package etcd

import (
	"encoding/json"

	"github.com/lightmen/nami/registry"
)

func marshal(inst *registry.Instance) (string, error) {
	data, err := json.Marshal(inst)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (inst *registry.Instance, err error) {
	err = json.Unmarshal(data, &inst)
	return
}
