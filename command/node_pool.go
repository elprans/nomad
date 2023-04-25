package command

import (
	"github.com/mitchellh/cli"
)

type NodePoolCommand struct {
	Meta
}

func (c *NodePoolCommand) Help() string {
	return `
can I haz help?`
}

func (c *NodePoolCommand) Synopsis() string {
	return "Node pool"
}

func (c *NodePoolCommand) Name() string { return "node_pool" }

func (c *NodePoolCommand) Run(args []string) int {
	return cli.RunResultHelp
}
