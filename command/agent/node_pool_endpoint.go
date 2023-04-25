// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package agent

import (
	"net/http"

	"github.com/hashicorp/nomad/nomad/structs"
)

func (s *HTTPServer) NodePoolsRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	args := structs.NodePoolListRequest{}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var out structs.NodePoolListResponse
	if err := s.agent.RPC("NodePool.ListNodePools", &args, &out); err != nil {
		return nil, err
	}

	setMeta(resp, &out.QueryMeta)
	if out.NodePools == nil {
		out.NodePools = make([]*structs.NodePool, 0)
	}
	return out.NodePools, nil
}

func (s *HTTPServer) NodePoolCreateRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "PUT" && req.Method != "POST" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	var pool structs.NodePool
	if err := decodeBody(req, &pool); err != nil {
		return nil, CodedError(http.StatusBadRequest, err.Error())
	}

	args := structs.NodePoolUpsertRequest{
		NodePool: &pool,
	}
	s.parseWriteRequest(req, &args.WriteRequest)

	var out structs.GenericResponse
	if err := s.agent.RPC("NodePool.UpsertNodePool", &args, &out); err != nil {
		return nil, err
	}
	setIndex(resp, out.Index)
	return nil, nil
}
