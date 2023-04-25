package command

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/nomad/api"
)

type NodePoolApplyCommand struct {
	Meta
}

func (c *NodePoolApplyCommand) Help() string {
	return `
can I haz help?`
}

func (c *NodePoolApplyCommand) Synopsis() string {
	return "Create or update a node pool"
}

func (c *NodePoolApplyCommand) Name() string { return "node-pool apply" }

func (c *NodePoolApplyCommand) Run(args []string) int {
	file := args[0]

	var config nodePoolConfig
	hclsimple.DecodeFile(file, nil, &config)

	// Get the HTTP client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	_, err = client.NodePools().Upsert(config.NodePool, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error applying namespace: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Successfully upserted node pool %q!", config.NodePool.Name))

	return 0
}

type nodePoolConfig struct {
	NodePool *api.NodePool `hcl:"node_pool,block"`
}
