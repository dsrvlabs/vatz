# VATZ API Specs

RPC service is available for managing VATZ. For now, this RPC service is only available on local.

## Get Plugin's status

Querying plugin's status.

```
~$ curl localhost:19091/v1/plugin_status
{"status":"OK","pluginStatus":[{"status":"FAIL","pluginName":"up check"}]}
```
