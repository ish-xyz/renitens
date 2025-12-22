package csi

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
	NodeService should host an "always running" PageServer, TrackerService
		PageServer can receive memory pages and store them to disk
		TrackerService exposes the existing checkpoints, its status, the checkpoint origin (where they come from)
			This service can be used by a scheduler plugin to improve scheduling decision

	StageVolume
		Should initiate a continous routine that checkpoints to X neighbors nodes
		If a checkpoint for this container is available locally restore it
		If a checkpoint is available to any of the neighbours (retrieve this info via trackerService),
			restore from them (in the future I should put this behaviour behind a flag - maybe a storage class?)

	UnstageVolume
		NodeService should initiate a delta-checkpoint to a neighbor node's PageServer
		to do so NodeService should contact a neighbor in a hashring or call a central directory server


**/

type NodeService struct {
	nodeID                      string
	csi.UnimplementedNodeServer // note to myself: this is needed for forward compatibility
}

var _ csi.NodeServer = &NodeService{}

var (
	volumeCaps = []csi.VolumeCapability_AccessMode{
		{
			Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
		{
			Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		},
		{
			Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
		},
	}
)

func NewNodeService(nodeID string) *NodeService {
	return &NodeService{
		nodeID: nodeID,
	}
}

func (n *NodeService) NodeStageVolume(ctx context.Context, request *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	// setup lazy restore server
	request.
	return nil, status.Error(codes.Unimplemented, "")
}

// NodeUnstageVolume is called by the CO when a workload that was using the specified volume is being moved to a different node.
func (n *NodeService) NodeUnstageVolume(ctx context.Context, request *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	request.
	return nil, status.Error(codes.Unimplemented, "")
}

// NodePublishVolume mounts the volume on the node.
func (n *NodeService) NodePublishVolume(ctx context.Context, request *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	// restore checkpoint as container image

	// volumeID := request.GetVolumeId()
	// if len(volumeID) == 0 {
	// 	return nil, status.Error(codes.InvalidArgument, "Volume id not provided")
	// }

	// target := request.GetTargetPath()
	// if len(target) == 0 {
	// 	return nil, status.Error(codes.InvalidArgument, "Target path not provided")
	// }

	// volCap := request.GetVolumeCapability()
	// if volCap == nil {
	// 	return nil, status.Error(codes.InvalidArgument, "Volume capability not provided")
	// }

	// if !isValidVolumeCapabilities([]*csi.VolumeCapability{volCap}) {
	// 	return nil, status.Error(codes.InvalidArgument, "Volume capability not supported")
	// }

	// readOnly := false
	// if request.GetReadonly() || request.VolumeCapability.AccessMode.GetMode() == csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY {
	// 	readOnly = true
	// }

	// options := make(map[string]string)
	// if m := volCap.GetMount(); m != nil {
	// 	for _, f := range m.MountFlags {
	// 		// get mountOptions from PV.spec.mountOptions
	// 		options[f] = ""
	// 	}
	// }

	// if readOnly {
	// 	// Todo add readonly in your mount options
	// }

	// TODO modify your volume mount logic here

	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmount the volume from the target path
func (n *NodeService) NodeUnpublishVolume(ctx context.Context, request *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	target := request.GetTargetPath()
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path not provided")
	}

	// TODO modify your volume umount logic here

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats get the volume stats
func (n *NodeService) NodeGetVolumeStats(ctx context.Context, request *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// NodeExpandVolume expand the volume
func (n *NodeService) NodeExpandVolume(ctx context.Context, request *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// NodeGetCapabilities get the node capabilities
func (n *NodeService) NodeGetCapabilities(ctx context.Context, request *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{}, nil
}

// NodeGetInfo get the node info
func (n *NodeService) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{NodeId: n.nodeID}, nil
}

// func isValidVolumeCapabilities(volCaps []*csi.VolumeCapability) bool {
// 	hasSupport := func(cap *csi.VolumeCapability) bool {
// 		for _, c := range volumeCaps {
// 			if c.GetMode() == cap.AccessMode.GetMode() {
// 				return true
// 			}
// 		}
// 		return false
// 	}

// 	foundAll := true
// 	for _, c := range volCaps {
// 		if !hasSupport(c) {
// 			foundAll = false
// 		}
// 	}
// 	return foundAll
// }
