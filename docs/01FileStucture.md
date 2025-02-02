```
raikuv/                  # Project root
├── cmd/                 # Entrypoints (CLI/client and server)
│   ├── node/            # Start a KV node (raft + storage)
│   │   └── main.go      # Initializes raft, storage, and cluster
│   └── client/          # CLI to interact with the cluster
│       └── main.go      # Sends PUT/GET commands
│
├── internal/            # Core implementation (not importable externally)
│   ├── storage/         # Storage engines (memory, LSM)
│   │   ├── memory.go    # In-memory + WAL
│   │   ├── wal.go       # Write-ahead log
│   │   └── lsm/         # LSM Tree (SSTables, compaction)
│   │       ├── sstable.go
│   │       └── compaction.go
│   │
│   ├── raft/            # Raft consensus
│   │   ├── node.go      # Leader/follower logic
│   │   ├── transport.go # RPCs (gRPC/HTTP)
│   │   └── state.go     # Persistent logs and metadata
│   │
│   ├── cluster/         # Cluster management
│   │   ├── membership.go # Node discovery (static peers or gossip)
│   │   └── sharding.go  # Consistent hashing for data distribution
│   │
│   └── api/             # Client-facing API (minimal)
│       ├── grpc/        # gRPC service (PUT/GET)
│       └── cli/         # CLI client (like `redis-cli`)
│
├── pkg/                 # Reusable utilities (optional)
│   ├── concurrency/     # Thread-safe primitives
│   └── serialization/   # Encoding/decoding (protobuf, JSON)
│
├── config/              # Configuration
│   └── config.yaml      # Node settings (ports, peers, storage path)
│
├── scripts/             # Helper scripts
│   ├── start-cluster.sh # Launch a 3-node cluster
│   └── benchmark.sh     # Stress-test the system
│
├── docker-compose.yaml  # Spin up nodes for local testing
├── Dockerfile           # Containerize a single node
├── Makefile             # Build/test shortcuts
├── go.mod
├── go.sum
└── README.md            # Docs, design notes, and API examples
```