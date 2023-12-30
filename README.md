# Kuiper - dynamic configuration service

This service was realized and organized in two endpoints.

- /create - endpoint for creation of versions of config or group
- /list - endpoint to list (timeseries representation) versions of exactly entity (config or group) by id and other parameters.
- /diff/list - endpoint to list (timeseries representation) diffs (of versions) of exactly entity (config or group) by id and other parameters.

## POST /create route

### Description

This route is used for creating versions of configuration or group, depends on provided type in json payload.
|parameter| type  |                    description                      |
|---------|-------|-----------------------------------------------------|
| type    | enum  | **Required.** Applicable values: *config* , *group* |

### Example of *create config* request

1. In case of creating fully new configuration or group (first version of entity), we should provide all informations except configurationID(unique identifier for config or group).

    Request - *create first version of config*

    ```json
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "configurationID":"",
        "versionTag": "v1",
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        }
    }
    ```

    Response - *create first version of config*

    ```json
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "appName": "app",
        "versionTag": "v1",
        "configurationID": "da97b0d6-bf49-4f14-8f36-1c77c23173e1",
        "createdAt": 1703086398,
        "weight": 1,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        }
    }
    ```

2. In case of creating new version of entity whom minimum one version already exists in system **we have to provide configurationID and new value of versionTag**.

    Request - *create new version of existing configuration*

    ```json
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "configurationID":"da97b0d6-bf49-4f14-8f36-1c77c23173e1",
        "versionTag": "v1.1",
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat-new"
            }
        }
    }
    ```

    Response - *create new version of existing configuration*

    ```json
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "appName": "app",
        "versionTag": "v1.1",
        "configurationID": "da97b0d6-bf49-4f14-8f36-1c77c23173e1",
        "createdAt": 1703109407,
        "weight": 2,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat-new"
            }
        },
        "diff": [
            {
                "type": "replace",
                "key": "etcdHostUAT",
                "new": "exampleCloud.timeseriesEtcd.cluster-uat-new",
                "old": "exampleCloud.timeseriesEtcd.cluster-uat"
            }
        ]
    }   
    ```

    Service on new version creation of existing config with exact id, will get newly before created version, compare new labels with old labels of config, and add appropriate diffs in diff slice. Type of diff can be: **addition**, **deletion** and **replace**.

### Example of *create group* request

1. Same as logic on **create config**, if we creating version of group for first time we shouldn't provide *configurationID*. Also for configs in group, we shouldn't provide id.

    Request - *create first version of group*

   ```json
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "configurationID":"",
        "versionTag": "v1",
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

    Response - *create first version of group*

    ```json
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "appName": "app",
        "versionTag": "v1",
        "configurationID": "8cea77df-6d07-4464-9bad-4ea219bf4247",
        "createdAt": 1703110310,
        "weight": 1,
        "type": "group",
        "config": {
            "configs": [
                {
                    "id": "fbec496b-d96a-4e0c-a3ba-e3fcb22c17ea",
                    "labels": {
                        "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany",
                        "etcdHostUAT": "exampleCloud.timeseriesEtcd.cluster-uat"
                    },
                    "diff": null
                },
                {
                    "id": "ddb2facc-58e3-43d4-b882-609ffea0b8ef",
                    "labels": {
                        "natsHostPROD": "exampleCloud.nats-provider.prod"
                    },
                    "diff": null
                }
            ]
        }
    }
    ```

2. If we have to create new version of group **we have to provide configurationID**, also if we have to edit configuration in group, we have to send id of **that group config** and edit labels.

3. If we have to add configuration in group, we have to provide new object to configs array **with labels map and without id field**.

4. If we have to delete configuration from group we have to send full configs list **without element which have to be deleted**.

### Cases 2, 3, 4 example

Request - *add config, delete config, and edit labels of existing config in group*

```json
{
    "namespace": "spacename",
    "creatorUsername": "silja",
    "configurationID":"8cea77df-6d07-4464-9bad-4ea219bf4247",
    "versionTag": "v2",
    "type": "group",
    "config": {
        "configs": [
            {
                "id": "fbec496b-d96a-4e0c-a3ba-e3fcb22c17ea",
                "labels": {
                    "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany"
                }
            },
            {
                "labels": {
                    "newConfigLabel": "newConfigLabelValue"
                }
            }
        ]
    }
}   
```

In this request we had ignore config which added on first version creation with id = "ddb2facc-58e3-43d4-b882-609ffea0b8ef" (delete of config in group), added new configuration with one example label and edited configuration in group with id = "fbec496b-d96a-4e0c-a3ba-e3fcb22c17ea" - deleted one of its label.

Response - *add config, delete config, and edit labels of existing config in group*

```json
{
    "namespace": "spacename",
    "creatorUsername": "silja",
    "appName": "app",
    "versionTag": "v2",
    "configurationID": "8cea77df-6d07-4464-9bad-4ea219bf4247",
    "createdAt": 1703110855,
    "weight": 2,
    "type": "group",
    "config": {
        "configs": [
            {
                "id": "fbec496b-d96a-4e0c-a3ba-e3fcb22c17ea",
                "labels": {
                    "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany"
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
                "id": "b3b2232e-456a-45e8-87c0-1a1836387674",
                "labels": {
                    "newConfigLabel": "newConfigLabelValue"
                },
                "diff": null
            }
        ]
    },
    "diff": [
        {
            "type": "addition",
            "key": "b3b2232e-456a-45e8-87c0-1a1836387674",
            "value": {
                "newConfigLabel": "newConfigLabelValue"
            }
        },
        {
            "type": "deletion",
            "key": "ddb2facc-58e3-43d4-b882-609ffea0b8ef",
            "value": {
                "natsHostPROD": "exampleCloud.nats-provider.prod"
            }
        }
    ]
}
```

Diffs in group are organized with: *diff of group* which is in root object and *config diff* which is in each configuration in group. In group diff array will be added items which describes **addition** or **deletion** of configs in group. In group config diff is organized same as only config, there can be **addition**, **deletion** or **replace** of label.

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

All list operations have sorting of versions to achieve timeseries structure(sort based on weight parameter which is increment integer of each version with starting position 1).

### Request list config URL example

*<http://localhost:8000/list?type=config&namespace=spejsnejm&appName=app&id=594012e8-ff3b-4db8-96ef-ee1a8fefc54d>*

### Response of list config request

```json
[
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "appName": "app",
        "versionTag": "v2",
        "configurationID": "2fdd969c-b821-4523-9c9c-874f5b3c77c4",
        "createdAt": 1703933385,
        "weight": 1,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostFR": "exampleCloud.timeseriesEtcd.cluster-france",
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-germany1123"
            }
        }
    },
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "appName": "app",
        "versionTag": "v3.beta",
        "configurationID": "2fdd969c-b821-4523-9c9c-874f5b3c77c4",
        "createdAt": 1703933432,
        "weight": 2,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostFR": "exampleCloud.timeseriesEtcd.cluster-france-beta",
                "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-germany-beta"
            }
        },
        "diff": [
            {
                "type": "replace",
                "key": "etcdHostFR",
                "new": "exampleCloud.timeseriesEtcd.cluster-france-beta",
                "old": "exampleCloud.timeseriesEtcd.cluster-france"
            },
            {
                "type": "replace",
                "key": "etcdHostGER",
                "new": "exampleCloud.timeseriesEtcd.cluster-germany-beta",
                "old": "exampleCloud.timeseriesEtcd.cluster-germany1123"
            }
        ]
    },
    {
        "namespace": "spacename",
        "creatorUsername": "silja",
        "appName": "app",
        "versionTag": "v3",
        "configurationID": "2fdd969c-b821-4523-9c9c-874f5b3c77c4",
        "createdAt": 1703933461,
        "weight": 3,
        "type": "config",
        "config": {
            "labels": {
                "etcdHostFR": "exampleCloud.timeseriesEtcd.cluster-france-v3"
            }
        },
        "diff": [
            {
                "type": "replace",
                "key": "etcdHostFR",
                "new": "exampleCloud.timeseriesEtcd.cluster-france-v3",
                "old": "exampleCloud.timeseriesEtcd.cluster-france-beta"
            },
            {
                "type": "deletion",
                "key": "etcdHostGER",
                "value": "exampleCloud.timeseriesEtcd.cluster-germany-beta"
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
        "weight": 1,
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
        "weight": 2,
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

## GET /diff/list route

This route is used for different list operations for diffs of versions. Depends on provided type in query parameters, it can be list diffs (of versions) of configuration or group (same as **/list** route, but returns timeseries representation list with only diffs and versionTag). Other parameters are concerned with needs of list operation (search).

### Request list config versions diffs URL example

*<http://localhost:8000/list?type=config&namespace=spejsnejm&appName=app&id=594012e8-ff3b-4db8-96ef-ee1a8fefc54d>*

### Response list config versions diffs example

```json
[
    {
        "versionTag": "v2",
        "diffs": null
    },
    {
        "versionTag": "v3.beta",
        "diffs": [
            {
                "type": "replace",
                "key": "etcdHostFR",
                "new": "exampleCloud.timeseriesEtcd.cluster-france-beta",
                "old": "exampleCloud.timeseriesEtcd.cluster-france"
            },
            {
                "type": "replace",
                "key": "etcdHostGER",
                "new": "exampleCloud.timeseriesEtcd.cluster-germany-beta",
                "old": "exampleCloud.timeseriesEtcd.cluster-germany1123"
            }
        ]
    },
    {
        "versionTag": "v3",
        "diffs": [
            {
                "type": "replace",
                "key": "etcdHostFR",
                "new": "exampleCloud.timeseriesEtcd.cluster-france-v3",
                "old": "exampleCloud.timeseriesEtcd.cluster-france-beta"
            },
            {
                "type": "deletion",
                "key": "etcdHostGER",
                "value": "exampleCloud.timeseriesEtcd.cluster-germany-beta"
            }
        ]
    }
]
```

(v2 diffs fields are empty because it is first version of config)

Each element of this list of timeseries configuration diffs contains versionTag field and diffs field (differences between this version and previous version).

### Request list group versions diffs URL example

*<http://localhost:8000/diff/list?type=group&id=991a2b2f-eae3-42b4-aafd-89e2be616f07&namespace=spacename&appName=app>*

### Response list group versions diffs example

```json
[
    {
        "versionTag": "v1",
        "groupConfigsDiff": [],
        "groupDiffs": null
    },
    {
        "versionTag": "v3",
        "groupConfigsDiff": [
            {
                "configID": "9d196653-d93a-43f8-aa71-f0daeef57cd2",
                "diffs": [
                    {
                        "type": "addition",
                        "key": "configLabelNew",
                        "value": "configLabelNew"
                    },
                    {
                        "type": "replace",
                        "key": "newConfigLabel",
                        "new": "newConfigLabelValue123",
                        "old": "newConfigLabelValue"
                    }
                ]
            }
        ],
        "groupDiffs": [
            {
                "type": "addition",
                "key": "29b4efc6-e14a-44fc-bbe7-b866a92e2e2b",
                "value": {
                    "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany"
                }
            }
        ]
    },
    {
        "versionTag": "v4.1",
        "groupConfigsDiff": [],
        "groupDiffs": [
            {
                "type": "deletion",
                "key": "29b4efc6-e14a-44fc-bbe7-b866a92e2e2b",
                "value": {
                    "etcdHostGER": "exampleCloud.timeseriesEtcd.cluster-dev-germany"
                }
            }
        ]
    }
]
```

(v1 diffs fields are empty because it is first version of group)

Each element of this list of timeseries group diffs contains versionTag, groupConfigsDiff and groupDiffs field.

**groupConfigsDiff** - diffs for each configuration in group(difference between those configurations in group with previous version of this configurations in group)

**groupDiffs** - diffs for each configuration in group(example: addition or deletion of configuration in group)
