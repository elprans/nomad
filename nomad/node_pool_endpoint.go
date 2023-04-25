package nomad

import (
	"fmt"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/nomad/nomad/state"
	"github.com/hashicorp/nomad/nomad/structs"
)

type NodePool struct {
	srv *Server
	ctx *RPCContext
}

func NewNodePoolEndpoint(srv *Server, ctx *RPCContext) *NodePool {
	return &NodePool{srv: srv, ctx: ctx}
}

func (n *NodePool) ListNodePools(args *structs.NodePoolListRequest, reply *structs.NodePoolListResponse) error {
	authErr := n.srv.Authenticate(n.ctx, args)
	if done, err := n.srv.forward("NodePool.ListNodePools", args, args, reply); done {
		return err
	}
	n.srv.MeasureRPCRate("node_pool", structs.RateMetricList, args)
	if authErr != nil {
		return structs.ErrPermissionDenied
	}
	defer metrics.MeasureSince([]string{"nomad", "node_pool", "list"}, time.Now())

	// Setup the blocking query
	opts := blockingOptions{
		queryOpts: &args.QueryOptions,
		queryMeta: &reply.QueryMeta,
		run: func(ws memdb.WatchSet, s *state.StateStore) error {
			// Iterate over all the namespaces
			var err error
			var iter memdb.ResultIterator

			iter, err = s.NodePools(ws)
			if err != nil {
				return err
			}

			reply.NodePools = nil
			for {
				raw := iter.Next()
				if raw == nil {
					break
				}
				pool := raw.(*structs.NodePool)
				reply.NodePools = append(reply.NodePools, pool)
			}

			// Use the last index that affected the namespace table
			index, err := s.Index(state.TableNodePools)
			if err != nil {
				return err
			}

			// Ensure we never set the index to zero, otherwise a blocking query cannot be used.
			// We floor the index at one, since realistically the first write must have a higher index.
			if index == 0 {
				index = 1
			}
			reply.Index = index
			return nil
		}}

	return n.srv.blockingRPC(&opts)
}

func (n *NodePool) UpsertNodePool(args *structs.NodePoolUpsertRequest, reply *structs.GenericResponse) error {
	authErr := n.srv.Authenticate(n.ctx, args)
	if done, err := n.srv.forward("NodePool.Upsert", args, args, reply); done {
		return err
	}
	n.srv.MeasureRPCRate("node_pool", structs.RateMetricWrite, args)
	if authErr != nil {
		return structs.ErrPermissionDenied
	}
	defer metrics.MeasureSince([]string{"nomad", "job", "register"}, time.Now())

	// Check management permissions.
	if aclObj, err := n.srv.ResolveACL(args); err != nil {
		return err
	} else if aclObj != nil && !aclObj.IsManagement() {
		return structs.ErrPermissionDenied
	}

	// Validate request.
	if args.NodePool == nil {
		return fmt.Errorf("missing node pool")
	}

	// Update via Raft
	_, index, err := n.srv.raftApply(structs.NodePoolUpsertRequestType, args)
	if err != nil {
		return err
	}

	// Update the index
	reply.Index = index
	return nil
}
