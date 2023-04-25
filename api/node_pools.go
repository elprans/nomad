package api

type NodePools struct {
	client *Client
}

func (c *Client) NodePools() *NodePools {
	return &NodePools{client: c}
}

func (n *NodePools) List(q *QueryOptions) ([]*NodePool, *QueryMeta, error) {
	var resp []*NodePool
	qm, err := n.client.query("/v1/node_pools", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

func (n *NodePools) Upsert(pool *NodePool, q *WriteOptions) (*WriteMeta, error) {
	wm, err := n.client.put("/v1/node_pool", pool, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

type NodePool struct {
	Name        string            `hcl:"name,label"`
	Description string            `hcl:"description,optional"`
	Meta        map[string]string `hcl:"meta,optional"`
	CreateIndex uint64
	ModifyIndex uint64
}
