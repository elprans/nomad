package command

import (
	"fmt"
	"sort"
)

type NodePoolListCommand struct {
	Meta
}

func (c *NodePoolListCommand) Help() string {
	return `
can I haz help?`
}

func (c *NodePoolListCommand) Synopsis() string {
	return "List node pools"
}

func (c *NodePoolListCommand) Name() string { return "node-pool list" }

func (c *NodePoolListCommand) Run(args []string) int {
	// Get the HTTP client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	pools, _, err := client.NodePools().List(nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error retrieving node pools: %s", err))
		return 1
	}

	// Sort the output by namespace name
	sort.Slice(pools, func(i, j int) bool { return pools[i].Name < pools[j].Name })

	rows := make([]string, len(pools)+1)
	rows[0] = "Name|Path|Description|Meta"
	for i, p := range pools {
		rows[i+1] = fmt.Sprintf("%s|%s|%s|%v",
			p.Name,
			p.Path,
			p.Description,
			p.Meta,
		)
	}

	c.Ui.Output(formatList(rows))
	return 0
}
