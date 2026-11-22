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
