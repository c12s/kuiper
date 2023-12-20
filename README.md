# Kuiper - dynamic configuration service

This service was realized and organized in two endpoints.

- /create/{type} - endpoint for creation of versions of config or group
- /list - endpoint to list versions of exactly entity (config or group) with id and other parameters.

## POST /create route

### Description

This route is used for creating versions of configuration or group, depends on provided type in json payload.
|parameter| type  |                    description                      |
|---------|-------|-----------------------------------------------------|
| type    | enum  | **Required.** Applicable values: *config* , *group* |

### Example of *create config* request

```json
{
    "namespace": "spacename",
    "creatorUsername": "silja",
    "configurationID":"",
    "tag": "v1",
    "type": "config",
    "config": {
        "labels": {
            "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
            "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
        }
    }
}
```

1. In case of creating fully new configuration or group (first version of entity), we should provide all informations except configurationID(unique identifier for config or group).

2. In case of creating new version of entity whom minimum one version already exists in system **we have to provide configurationID**.

### Response - *create config*

```json
{
    "namespace": "spacename",
    "creatorUsername": "silja",
    "appName": "app",
    "tag": "v1",
    "configurationID": "da97b0d6-bf49-4f14-8f36-1c77c23173e1",
    "createdAt": 1703086398,
    "type": "config",
    "config": {
        "labels": {
            "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
            "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
        }
    }
}
```

### Example of *create group* request

```json
{
    "namespace": "spacename",
    "creatorUsername": "silja",
    "configurationID":"",
    "tag": "v1",
    "type": "group",
    "config": {
        "configs": [
            {
                "labels": {
                    "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                    "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
                }
            },
            {
                
                "labels": {
                    "natsHostPROD": "exampleCloud.nats-provider.prod"
                }
            }
        ]
    }
}
```

1. Same as logic on **create config**, if we creating version of group for first time we shouldn't provide *configurationID*. Also for configs in group, we shouldn't provide id.

2. If we have to create version new version of group **we have to provide configurationID**, also if we have to edit configuration in group, we have to send id and edit labels.

3. If we have to add configuration in group, we have to provide new object **with labels map and without id field**

4. If we have to delete configuration from group we have to send full configs list **without element which have to be deleted**.

### Response - *create group*

```json
{
    "namespace": "spacename",
    "creatorUsername": "silja",
    "configurationID":"",
    "tag": "v1",
    "type": "group",
    "config": {
        "configs": [
            {
                "id": "18c591e6-a9da-4bb6-a350-52a4505d3818",
                "labels": {
                    "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                    "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
                }
            },
            {
                "id": "d6b72ac6-1223-4b7b-842f-1818188e1db4",
                "labels": {
                    "natsHostPROD": "exampleCloud.nats-provider.prod"
                }
            }
        ]
    }
}
```

## GET /list route

This route is used for different list operations for versions. Depends on provided type in query parameters, it can be list of versions of configuration or group (same as **create** route). Other parameters are concerned with needs of list operation (search).

|parameter| type  |                    description                      |
|---------|-------|-----------------------------------------------------|
| type    | enum  | **Required.** Applicable values: *config* , *group* |
| id   | uuid  | **Required.** Unique identifier of entity (config or group) for which a listing is requested |
| namespace    | string  | **Not required.** (if not provided service will use default value - *namespace*)  |
| appName    | string  | **Not required.** (if not provided service will use default value - *app*) |
| fromVersion | string | **Not required.** Parameter which represents version tag from which will listing start
| withFrom | string | **Not required.** Parameter which represents should fromVersion be in list or not
| toVersion | string | **Not required.** Parameter which represents version tag from which will listing stop
| withTo | string | **Not required.** Parameter which represents should toVersion be in list or not
| sortType | enum | **Not required.** Applicable values: *lexically* , *timestamp*. Parameter which represents a sortType of return value. default: *lexically*

###

When we provide fromVersion and toVersion, database will return all object whose key fall in lexically range. Example if we provide next query params:

- type = group
- namespace = namespace
- app = app
- id = 123-123-123
- fromVersion = v1

Route URL example: *<http://localhost:8000/list?type=group&namespace=namespace&app=app&id=123-123-123&fromVersion=v1>*

it will return all versions of group with id which have lexically higher version than v1.

example:

- /namespace/app/group/123-123-123/v2

- /namespace/app/group/123-123-123/v3

- and all versions whose key have prefix */namespace/app/group/123-123-123/* and whom lexically higher than v1 version tag

but we shouldn't get element v1, **because we didn't provided withFrom = true parameter**.

### Request list group URL example

*<http://localhost:8000/list?type=config&namespace=spejsnejm&appName=app&id=594012e8-ff3b-4db8-96ef-ee1a8fefc54d&sortType=timestamp>*

### Response of list config request

```json
[
    {
        "namespace": "spejsnejm",
        "creatorUsername": "silja",
        "appName": "app",
        "tag": "v10",
        "configurationID": "594012e8-ff3b-4db8-96ef-ee1a8fefc54d",
        "createdAt": 1702748764,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        }
    },
    {
        "namespace": "spejsnejm",
        "creatorUsername": "silja",
        "appName": "app",
        "tag": "v1",
        "configurationID": "594012e8-ff3b-4db8-96ef-ee1a8fefc54d",
        "createdAt": 1702748787,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        }
    },
    {
        "namespace": "spejsnejm",
        "creatorUsername": "silja",
        "appName": "app",
        "tag": "micko",
        "configurationID": "594012e8-ff3b-4db8-96ef-ee1a8fefc54d",
        "createdAt": 1702748792,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        }
    },
    {
        "namespace": "spejsnejm",
        "creatorUsername": "silja",
        "appName": "app",
        "tag": "laza",
        "configurationID": "594012e8-ff3b-4db8-96ef-ee1a8fefc54d",
        "createdAt": 1702748806,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany"
            }
        },
        "diff": [
            {
                "type": "deletion",
                "key": "etcdHostUAT",
                "value": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        ]
    },
    {
        "namespace": "spejsnejm",
        "creatorUsername": "silja",
        "appName": "app",
        "tag": "pera",
        "configurationID": "594012e8-ff3b-4db8-96ef-ee1a8fefc54d",
        "createdAt": 1702748819,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        },
        "diff": [
            {
                "type": "deletion",
                "key": "etcdHostGER",
                "value": "exampleCloud.timeseriesEtcd.cluster-dev-germany"
            }
        ]
    }
]
```

### Request list group URL example

*<http://localhost:8000/list?type=group&id=7f23cfd3-10a9-42b8-9925-c01849e95909&namespace=spejsnejm&appName=app&sortType=timestamp>*

### Response of list group request

```json
[
    {
        "namespace": "spejsnejm",
        "creatorUsername": "silja",
        "appName": "app",
        "tag": "v1",
        "configurationID": "7f23cfd3-10a9-42b8-9925-c01849e95909",
        "createdAt": 1702578908,
        "type": "group",
        "config": {
            "configs": [
                {
                    "id": "18c591e6-a9da-4bb6-a350-52a4505d3818",
                    "labels": {
                        "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany-1",
                        "etcdHostUATDEV": "exampleCloud.timeseriesEtcd.cluster-uatdev"
                    },
                    "diff": null
                },
                {
                    "id": "d6b72ac6-1223-4b7b-842f-1818188e1db4",
                    "labels": {
                        "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany-12312312"
                    },
                    "diff": null
                },
                {
                    "id": "89483fb3-a2ba-4d4b-8684-da1505205f84",
                    "labels": {
                        "etcdHostUATNew": "exampleCloud.timeseriesEtcd.cluster-uat-new"
                    },
                    "diff": null
                }
            ]
        }
    },
    {
        "namespace": "spejsnejm",
        "creatorUsername": "silja",
        "appName": "app",
        "tag": "v2",
        "configurationID": "7f23cfd3-10a9-42b8-9925-c01849e95909",
        "createdAt": 1702578923,
        "type": "group",
        "config": {
            "configs": [
                {
                    "id": "b8e0ab29-fcaf-422b-a663-b059eecf5ab9",
                    "labels": {
                        "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                        "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
                    },
                    "diff": null
                },
                {
                    "id": "18c591e6-a9da-4bb6-a350-52a4505d3818",
                    "labels": {
                        "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany-1",
                        "etcdHostUATDEV": "exampleCloud.timeseriesEtcd.cluster-uatdev"
                    },
                    "diff": null
                },
                {
                    "id": "d6b72ac6-1223-4b7b-842f-1818188e1db4",
                    "labels": {
                        "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany-12312312"
                    },
                    "diff": null
                },
                {
                    "id": "89483fb3-a2ba-4d4b-8684-da1505205f84",
                    "labels": {
                        "etcdHostUATNew": "exampleCloud.timeseriesEtcd.cluster-uat-new"
                    },
                    "diff": null
                }
            ]
        },
        "diff": [
            {
                "type": "addition",
                "key": "b8e0ab29-fcaf-422b-a663-b059eecf5ab9",
                "value": {
                    "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                    "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
                }
            }
        ]
    }
]
```
