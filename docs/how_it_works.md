Assumptions:
    Every node runs a criu page server to receive remote checkpoints
    There's a centralized TrackerService for checkpoints
        TrackerService has access to the kubernetes api to read nodes, workloads, etc (TBD)
    The node "agent" will interface directly with CRIU/CRI-O.

Creation:
    Pods gets submitted through the API server
    if an annotation called renitens.io/$container/restore=true is present the webhook will call the checkpoint TrackerService passing the address ($namespace/$pod/$container) and ask where that specific checkpoint is stored:
        if no checkpoint is found, the image address will remain the same
        if a checkpoint is found, the image address will change to the name of the node where the checkpoint is present
Running:
    Once the pod is created, if the annotation "renitens.io/$container/checkpoint=true" is present, the controller will notice the new pod, calculate the appropriate neighbours to store its checkpoints, expose these information via the TrackerService. E.g.:
        namespace/pod -> {checkpoint to nodes: 1,2,3 - every 5mins - refresh every 10s, checkpoints: [], latestCheckpoint: -1}

    The agent running on the workload node, will call the above TrackerServer API every x seconds (based on the configuration) and initiate checkpoints to the neighboring nodes.
    The node agent will also call the tracker service and register or unregister checkpoints. (not sure about this bit, might be better to have a system that avoids RPCs).
    Checkpoints will be executed using criu/cri-o and pointing to the neighbours page servers (authentication/authorization here is yet to be understood).
    Checkpoints will be taken on a scheduled/regularly basis.

Deletion
    Moving node
    Terminating
        When terminating 

Drain



Components:

- Controller
- Mutating webhook
- Node Agent
    - Page Server
    - 
