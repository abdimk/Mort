# Mort

> **Real-Time Distributed Cache**

*This is for experimental purposes*

---

## About

Mort is a real-time distributed cache built for high-performance systems. Written in Go, it provides:

- Ultra-low latency data access
- Automatic distribution across nodes
- Efficient memory management

Mort is designed to handle large-scale workloads with features like intelligent caching, horizontal scalability, and high availability, making it ideal for modern applications, microservices, and data-intensive platforms that need fast and reliable data access.

---

## Design Goals

> **Note:** This is experimental. The core idea is to implement all the features that Redis has and improve on things Redis lacks.

1. **Implement core Redis features** (must-have)
2. **Add features Redis lacks or does poorly** (differentiation)
3. **Leverage Go strengths** (concurrency, networking)

---

## Architecture Layers

```
┌─────────────────────────────────────┐
│     Advanced Features Layer         │
│  (Stampede protection, WASM, etc.)  │
├─────────────────────────────────────┤
│     Distributed System Layer        │
│  (Sharding, Replication, Gossip)    │
├─────────────────────────────────────┤
│     Core Cache Layer                │
│  (KV Storage, TTL, Data Structures) │
└─────────────────────────────────────┘
```

---

## 1. Core Features (Redis-like Basics)

### Key–Value Storage

| Command | Description |
|---------|-------------|
| `SET key value` | Store a value |
| `GET key` | Retrieve a value |
| `DEL key` | Delete a key |
| `EXISTS key` | Check key existence |
| `MGET / MSET` | Multi-key operations |

**Implementation:**
- `map[string]Value` + `sync.RWMutex`
- Or **sharded maps** for better performance

### TTL / Expiration

Redis-style expiration with `SET key value EX 60`:

- Per-key TTL
- Lazy expiration
- Background expiration worker

**Implementation options:**
- Min heap
- Time wheel
- Periodic scan

### Data Structures

| Structure | Use Case |
|-----------|----------|
| String | Standard cache |
| Hash | Objects |
| List | Queues |
| Set | Uniqueness |
| Sorted Set | Leaderboards |
| Bitmap | Analytics |
| HyperLogLog | Cardinality |

**Minimum required:** `SET`, `HASH`, `SET`, `SORTED SET`

### Atomic Operations

```
INCR, DECR, INCRBY
```

Critical for:
- Rate limiting
- Counters
- Analytics

### Pipelining

Clients send multiple commands without waiting:

```
SET a 1
SET b 2
SET c 3
```

Huge performance improvement.

### Transactions (Optional)

Redis supports: `MULTI`, `EXEC`, `WATCH`

---

## 2. Distributed System Features

### Cluster Sharding

Distribute keys across nodes using **consistent hashing** (hash ring).

### Replication

| Mode | Behavior |
|------|----------|
| Sync | Safe but slower |
| Async | Redis default |

**Failover:**
1. Detect failure
2. Elect new primary
3. Algorithms: Raft, Gossip + leader election

### Gossip Protocol

Nodes share cluster state:
```
node A -> node B
node B -> node C
```

Used by: Cassandra, Redis Cluster, Scylla

### Distributed Locking

Better than Redis `SET key value NX PX`:
- Redlock alternative
- Fencing tokens
- Lease locks

---

## 3. Redis Limitations (Opportunities to Improve)

| Limitation | Mort's Solution |
|------------|-----------------|
| LRU/LFU eviction | **TinyLFU** (Caffeine, Ristretto) |
| Single-threaded | **Goroutines + sharded storage** |
| Cache stampede | **Request coalescing / singleflight** |
| No read-through cache | `cache.GetOrFetch(key, dbCall)` |
| Single-region | **Multi-region replication** |
| No rate limiting | **Built-in rate limiter** |
| No event streaming | **Streaming / Event Bus** |
| Limited observability | **Native Prometheus/OpenTelemetry** |
| No hot key handling | **Hot key detection & replication** |
| No compute near data | **WASM functions** |

---

## 4. Network Protocol

Options:
- **Redis Protocol (RESP):** Best for existing client compatibility
- **gRPC / HTTP:** Better for microservices

---

## 5. Persistence

| Type | Description |
|------|-------------|
| Snapshot | Every N minutes |
| Append Only Log | SET / DEL log |

Redis uses both.

---

## 6. High Performance Internals

**Techniques:**
- Arena allocation
- Memory pooling
- Zero-copy

**Go tools:**
- `sync.Pool`
- Byte buffers
- Sharded storage (e.g., 256 shards with map + mutex)

---

## 7. Suggested Architecture

```
Client
   │
   ▼
TCP Server (RESP)
   │
   ▼
Command Router
   │
   ▼
Shard Manager
   │
   ▼
Local Storage
   │
   ▼
Replication Layer
   │
   ▼
Cluster Manager
```

---

## 8. Go Libraries Worth Studying

| Library | Purpose |
|---------|---------|
| [Ristretto](https://github.com/dgraph-io/ristretto) | Best cache library |
| [groupcache](https://github.com/golang/groupcache) | Distributed caching |
| [Olric](https://github.com/buraksezer/olric) | Distributed cache |
| [DragonflyDB](https://github.com/dragonflydb/dragonfly) | Redis alternative |

---

## 9. Killer Features to Beat Redis

- [ ] TinyLFU eviction
- [ ] Hot key replication
- [ ] Built-in rate limiting
- [ ] Cache stampede protection
- [ ] Multi-region replication
- [ ] WASM compute
- [ ] Native observability
- [ ] Auto-scaling cluster

---

## 10. Current Implementation Status

### ✅ Implemented

- [x] Basic SET/GET commands
- [x] TTL-based expiration
- [x] Leader-follower replication
- [x] Thread-safe cache operations

### 📋 Planned

See sections above for detailed design.

---

## Project Structure

```
Mort/
├── main.go          # Server bootstrap
├── server.go        # TCP server & replication
├── command.go       # Protocol parser
├── go.mod           # Go module definition
├── Makefile         # Build automation
├── cache/
│   ├── cache.go     # In-memory cache implementation
│   └── cacher.go    # Cache interface
└── bin/
    └── Mort         # Compiled binary
```

---

## Usage

### Start as Leader

```bash
./bin/Mort -listenaddr :3000
```

### Start as Follower

```bash
./bin/Mort -listenaddr :3001 -leaderaddr localhost:3000
```

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25.2 |
| Networking | TCP sockets |
| Concurrency | Goroutines & channels |
| Synchronization | `sync.RWMutex` |

---

## License

Experimental project
