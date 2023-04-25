// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package structs

type NodePool struct {
	Name        string
	Description string
	Meta        map[string]string

	CreateIndex uint64
	ModifyIndex uint64
}

type NodePoolListRequest struct {
	QueryOptions
}

type NodePoolListResponse struct {
	NodePools []*NodePool
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
