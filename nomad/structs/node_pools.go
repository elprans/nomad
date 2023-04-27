// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package structs

import "fmt"

type NodePool struct {
	Name        string
	Description string
	Meta        map[string]string
	Children    []*NodePool
	Path        string

	CreateIndex uint64
	ModifyIndex uint64
}

func (n *NodePool) Canonicalize() {
	if n.Name != "all" && n.Path == "" {
		n.Path = fmt.Sprintf("/all/%s", n.Name)
	}

	for _, p := range n.Children {
		p.Path = fmt.Sprintf("%s/%s", n.Path, p.Name)
		p.Canonicalize()
	}
}

type NodePoolListRequest struct {
	QueryOptions
}

type NodePoolListResponse struct {
	NodePools []*NodePool
	QueryMeta
}

type NodePoolListNodesRequest struct {
	Name string
	QueryOptions
}

type NodePoolListNodesResponse struct {
	Nodes []*NodeListStub
	QueryMeta
}

type NodePoolUpsertRequest struct {
	NodePool *NodePool
	WriteRequest
}

type NodePoolDeleteRequest struct {
	Name string
	WriteRequest
}
