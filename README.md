# kuiper — c12s config store and placement service

The configuration plane: stores **config groups** and **standalone configs**
(versions in etcd), validates them against schemas via Quasar, and *places*
them — selecting nodes, recording a placement task per node, and delivering
the config through Gravity with an HTTP webhook for the node's reply.

## Role in c12s

```
cockpit put/place ──▶ gateway ──▶ kuiper ──▶ quasar  (validate config against schema)
                                     ├──▶ etcd   (kuiper_etcd: configs, placements)
                                     ├──▶ magnetar (select nodes: org-owned + label query)
                                     ├──▶ gravity  (DeseminateConfig per node, auth metadata forwarded)
                                     └──▶ oort     (async: namespace→config inheritance rel)
node reply ──▶ gravity webhook ──▶ kuiper webhook server (:9002 /standalone, /groups)
                                     └──▶ placement task status: Accepted → Placed / Failed
```

Authorization uses the JWT `permissions` claim (`config.put`, `config.get`,
`namespace.put`, ... — Apollo-minted via Oort).

## API (`pkg/api/proto/kuiper.proto`, reflection enabled)

| Config groups | Standalone configs |
|---|---|
| `PutConfigGroup`, `GetConfigGroup`, `ListConfigGroup`, `DeleteConfigGroup` | `PutStandaloneConfig`, `GetStandaloneConfig`, `ListStandaloneConfig`, `DeleteStandaloneConfig` |
| `PlaceConfigGroup`, `ListPlacementTaskByConfigGroup` | `PlaceStandaloneConfig`, `ListPlacementTaskByStandaloneConfig` |
| `DiffConfigGroup` | `DiffStandaloneConfig` |

Placement strategy (`PlaceReq.Strategy`):

| `name` | Selection |
|---|---|
| `default` | magnetar `QueryOrgOwnedNodes` with the given label `query` (empty query = all org nodes) |
| `gossip` | magnetar `ListOrgOwnedNodes`, then a random sample of `percentage` (rounded up) |

Namespaces themselves are owned by **Meridian** (`AddNamespace`,
`GetNamespaceHierarchy`, ...), not Kuiper; Kuiper looks namespaces up through
its Meridian client only where needed.

## Running

Compose service `kuiper`: gRPC on `${KUIPER_PORT}` (9001) and the task
webhook server on `${KUIPER_WEBHOOK_PORT}` (9002). Depends on its own etcd
(`kuiper_etcd`), magnetar, quasar, gravity, oort, NATS, meridian.

```bash
go build ./... && go test ./...   # builds; no tests exist
go run ./cmd/server               # needs the env below
```

## Configuration (env)

| Variable | Used for |
|---|---|
| `KUIPER_ADDRESS` | gRPC listen address (`:9001`) |
| `WEBHOOK_ADDRESS` / `WEBHOOK_URL` | task-reply webhook listen addr and the URL given to Gravity |
| `ETCD_ADDRESS` | kuiper etcd (`kuiper_etcd:2379`) |
| `MAGNETAR_ADDRESS`, `AGENT_QUEUE_ADDRESS`, `QUASAR_ADDRESS`, `OORT_ADDRESS`, `NATS_ADDRESS` | downstream clients |
| `MERIDIAN_ADDRESS` | read directly in startup (not via the config struct) |
| `SECRET_KEY` | JWT permission checks |

## Data and dependencies

- **etcd** (`kuiper_etcd`, volume `kuiper_etcd_data`):
  `groups/<org>/<ns>/<name>/<version>`,
  `standalone/<org>/<ns>/<name>/<version>`,
  `placements/<type>/<org>/<ns>/<name>/<version>/<taskId>`.
  Config versions are immutable-ish: `Put` fails with *version exists* for a
  duplicate `(org, ns, name, version)`; a new version is a new key.
- **Placement task statuses**: `Accepted` (created + disseminate attempted),
  `Placed` / `Failed` (set from node webhooks).
- **Oort**: namespace→config inheritance relationship per stored config (async).

## Status

Stable for the config chain as exercised by the experiments and manual
cockpit runs; the placement robustness gaps below are accepted for the lab.

## Bugs, disabled and unfinished features

### Confirmed bugs

| Where | Effect |
|---|---|
| `internal/services/placement.go` — `Place` | Node-selection errors are **discarded**: `nodes, _ = s.placeByQuery(...)` / `placeByGossip(...)`. If Magnetar is unreachable the call returns **zero tasks and no error** — the operator sees a successful placement that reached nobody. |
| `internal/services/placement.go` — `Place` | Per-node failures (`store.Place` error, command marshal error) are logged and skipped; the returned task list silently omits nodes. A partially placed config looks like a success. |
| `selectRandmNodes` | ~~RNG re-seeded with `time.Now().Unix()` on every call: two gossip placements in the same second select the same nodes.~~ Fixed: `rand.Perm` over the global source, empty/zero/negative guards, clamp to 100, no mutation of the caller's slice. Note `ceil()` means a 5 % gossip with 1–19 nodes always selects at least 1 node — fine for intent, worth knowing. |
| `internal/services/*.go` — Oort calls | `CreateInheritanceRel` is fire-and-forget (callback only logs): a failed policy write leaves the config stored but not grant-propagated. `Delete` does **not** remove the Oort relationship at all — deleted configs keep their namespace→config edge in the graph. |
| Stale placement tasks | Delivery to an offline node produces no webhook, so the task stays `Accepted` **forever** — there is no timeout/GC job and no retry (gravity's fire-and-forget path). |

### Disabled / unfinished features

- **Tests:** placement selection now has unit tests
  (`internal/services/placement_test.go`); etcd store and webhook handling do
  not.
- Task webhook server has no authentication: anything that can reach
  `kuiper:9002` can forge `Placed`/`Failed` statuses (lab boundary; the real
  node reply path goes through gravity).
- Diff is implemented client-side over two full `Get`s (both must be
  readable); no pagination on list endpoints.
- Placement has no scheduling/rollout concept: `gossip` is a one-shot random
  sample, not a progressive rollout (that role belongs to `cockpit image
  rollout` for OS updates).

### PoC boundaries (by design, not bugs)

- Plaintext gRPC/HTTP/NATS inside the compose network; no TLS.
- `Kuiper → Gravity` calls forward the caller's `authz-token` metadata (fix
  `fb118cb`); NATS 503 no-responder replies are skipped rather than decoded
  (fix `83bbe82`).
- Config bodies are opaque protobuf blobs once validated — Kuiper does not
  interpret parameters beyond validation and diffing.

References: `tools/docs/README.md` (corrected API conventions),
`../plan/architecture.md` (placement flow), `../magnetar/README.md`,
`../gravity/README.md`.
