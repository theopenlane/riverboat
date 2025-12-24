# object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**refreshInterval**|`integer`|||
|[**river**](#river)|`object`|Config is the configuration for the river server<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "river": {
        "queues": [
            {}
        ],
        "workers": {
            "openlaneconfig": {},
            "emailworker": {
                "config": {}
            },
            "databaseworker": {
                "config": {}
            },
            "exportcontentworker": {
                "config": {}
            },
            "deleteexportcontentworker": {
                "config": {}
            },
            "slackworker": {
                "config": {}
            }
        },
        "trustcenterworkers": {
            "openlaneconfig": {},
            "createcustomdomainworker": {
                "config": {}
            },
            "validatecustomdomainworker": {
                "config": {}
            },
            "deletecustomdomainworker": {
                "config": {}
            },
            "watermarkdocworker": {
                "config": {}
            },
            "createpirschdomainworker": {
                "config": {}
            },
            "deletepirschdomainworker": {
                "config": {}
            },
            "updatepirschdomainworker": {
                "config": {}
            },
            "createpreviewdomainworker": {
                "config": {}
            },
            "deletepreviewdomainworker": {
                "Config": {}
            },
            "validatepreviewdomainworker": {
                "config": {}
            }
        },
        "metrics": {}
    }
}
```

<a name="river"></a>
## river: object

Config is the configuration for the river server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**databasehost**|`string`|DatabaseHost for connecting to the postgres database<br/>||
|[**queues**](#riverqueues)|`array`|||
|[**workers**](#riverworkers)|`object`|Workers that will be enabled on the server<br/>||
|[**trustcenterworkers**](#rivertrustcenterworkers)|`object`|||
|**defaultmaxretries**|`integer`|DefaultMaxRetries is the maximum number of retries for failed jobs, this can be set differently per job<br/>||
|[**metrics**](#rivermetrics)|`object`|MetricsConfig is the configuration for metrics<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "queues": [
        {}
    ],
    "workers": {
        "openlaneconfig": {},
        "emailworker": {
            "config": {}
        },
        "databaseworker": {
            "config": {}
        },
        "exportcontentworker": {
            "config": {}
        },
        "deleteexportcontentworker": {
            "config": {}
        },
        "slackworker": {
            "config": {}
        }
    },
    "trustcenterworkers": {
        "openlaneconfig": {},
        "createcustomdomainworker": {
            "config": {}
        },
        "validatecustomdomainworker": {
            "config": {}
        },
        "deletecustomdomainworker": {
            "config": {}
        },
        "watermarkdocworker": {
            "config": {}
        },
        "createpirschdomainworker": {
            "config": {}
        },
        "deletepirschdomainworker": {
            "config": {}
        },
        "updatepirschdomainworker": {
            "config": {}
        },
        "createpreviewdomainworker": {
            "config": {}
        },
        "deletepreviewdomainworker": {
            "Config": {}
        },
        "validatepreviewdomainworker": {
            "config": {}
        }
    },
    "metrics": {}
}
```

<a name="riverqueues"></a>
### river\.queues: array

**Items**

**Example**

```json
[
    {}
]
```

<a name="riverworkers"></a>
### river\.workers: object

Workers that will be enabled on the server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**openlaneconfig**](#riverworkersopenlaneconfig)|`object`|OpenlaneConfig contains the configuration for connecting to the Openlane API.<br/>||
|[**emailworker**](#riverworkersemailworker)|`object`|EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider<br/>||
|[**databaseworker**](#riverworkersdatabaseworker)|`object`|DatabaseWorker is a worker to create a dedicated database for an organization<br/>||
|[**exportcontentworker**](#riverworkersexportcontentworker)|`object`|ExportContentWorker exports the content into csv and makes it downloadable<br/>||
|[**deleteexportcontentworker**](#riverworkersdeleteexportcontentworker)|`object`|DeleteExportContentWorker deletes exports that are older than the configured cutoff duration<br/>||
|[**slackworker**](#riverworkersslackworker)|`object`|SlackWorker sends messages to Slack.<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "openlaneconfig": {},
    "emailworker": {
        "config": {}
    },
    "databaseworker": {
        "config": {}
    },
    "exportcontentworker": {
        "config": {}
    },
    "deleteexportcontentworker": {
        "config": {}
    },
    "slackworker": {
        "config": {}
    }
}
```

<a name="riverworkersopenlaneconfig"></a>
#### river\.workers\.openlaneconfig: object

OpenlaneConfig contains the configuration for connecting to the Openlane API.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>||
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersemailworker"></a>
#### river\.workers\.emailworker: object

EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersemailworkerconfig)|`object`|EmailConfig contains the configuration for the email worker<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersemailworkerconfig"></a>
##### river\.workers\.emailworker\.config: object

EmailConfig contains the configuration for the email worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|enable or disable the email worker<br/>||
|**devmode**|`boolean`|enable dev mode<br/>||
|**testdir**|`string`|the directory to use for dev mode<br/>||
|**token**|`string`|the token to use for the email provider<br/>||
|**fromemail**|`string`|FromEmail is the email address to use as the sender<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersdatabaseworker"></a>
#### river\.workers\.databaseworker: object

DatabaseWorker is a worker to create a dedicated database for an organization


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersdatabaseworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersdatabaseworkerconfig"></a>
##### river\.workers\.databaseworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|Enable the dbx client<br/>||
|**baseUrl**|`string`|Base URL for the dbx service<br/>||
|**endpoint**|`string`|Endpoint for the graphql api<br/>||
|**debug**|`boolean`|Enable debug mode<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersexportcontentworker"></a>
#### river\.workers\.exportcontentworker: object

ExportContentWorker exports the content into csv and makes it downloadable


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersexportcontentworkerconfig)|`object`|ExportWorkerConfig configuration for the export content worker<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersexportcontentworkerconfig"></a>
##### river\.workers\.exportcontentworker\.config: object

ExportWorkerConfig configuration for the export content worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>||
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>||
|**enabled**|`boolean`|||

**Additional Properties:** not allowed  
<a name="riverworkersdeleteexportcontentworker"></a>
#### river\.workers\.deleteexportcontentworker: object

DeleteExportContentWorker deletes exports that are older than the configured cutoff duration


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersdeleteexportcontentworkerconfig)|`object`|DeleteExportWorkerConfig holds the configuration for the delete export worker<br/>|yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersdeleteexportcontentworkerconfig"></a>
##### river\.workers\.deleteexportcontentworker\.config: object

DeleteExportWorkerConfig holds the configuration for the delete export worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>|no|
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>|no|
|**enabled**|`boolean`||no|
|**interval**|`integer`||yes|
|**cutoffduration**|`integer`|CutoffDuration defines the tolerance for exports. If you set 30 minutes, all exports older than 30 minutes<br/>at the time of job execution will be deleted<br/>|yes|

**Additional Properties:** not allowed  
<a name="riverworkersslackworker"></a>
#### river\.workers\.slackworker: object

SlackWorker sends messages to Slack.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersslackworkerconfig)|`object`|SlackConfig configures the Slack worker.<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersslackworkerconfig"></a>
##### river\.workers\.slackworker\.config: object

SlackConfig configures the Slack worker.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|enable or disable the slack worker<br/>||
|**devmode**|`boolean`|enable dev mode<br/>||
|**token**|`string`|the token to use for the slack app<br/>||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkers"></a>
### river\.trustcenterworkers: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**openlaneconfig**](#rivertrustcenterworkersopenlaneconfig)|`object`|OpenlaneConfig contains the configuration for connecting to the Openlane API.<br/>||
|[**createcustomdomainworker**](#rivertrustcenterworkerscreatecustomdomainworker)|`object`|||
|[**validatecustomdomainworker**](#rivertrustcenterworkersvalidatecustomdomainworker)|`object`|||
|[**deletecustomdomainworker**](#rivertrustcenterworkersdeletecustomdomainworker)|`object`|||
|[**watermarkdocworker**](#rivertrustcenterworkerswatermarkdocworker)|`object`|||
|[**createpirschdomainworker**](#rivertrustcenterworkerscreatepirschdomainworker)|`object`|||
|[**deletepirschdomainworker**](#rivertrustcenterworkersdeletepirschdomainworker)|`object`|||
|[**updatepirschdomainworker**](#rivertrustcenterworkersupdatepirschdomainworker)|`object`|||
|[**createpreviewdomainworker**](#rivertrustcenterworkerscreatepreviewdomainworker)|`object`|||
|[**deletepreviewdomainworker**](#rivertrustcenterworkersdeletepreviewdomainworker)|`object`|||
|[**validatepreviewdomainworker**](#rivertrustcenterworkersvalidatepreviewdomainworker)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "openlaneconfig": {},
    "createcustomdomainworker": {
        "config": {}
    },
    "validatecustomdomainworker": {
        "config": {}
    },
    "deletecustomdomainworker": {
        "config": {}
    },
    "watermarkdocworker": {
        "config": {}
    },
    "createpirschdomainworker": {
        "config": {}
    },
    "deletepirschdomainworker": {
        "config": {}
    },
    "updatepirschdomainworker": {
        "config": {}
    },
    "createpreviewdomainworker": {
        "config": {}
    },
    "deletepreviewdomainworker": {
        "Config": {}
    },
    "validatepreviewdomainworker": {
        "config": {}
    }
}
```

<a name="rivertrustcenterworkersopenlaneconfig"></a>
#### river\.trustcenterworkers\.openlaneconfig: object

OpenlaneConfig contains the configuration for connecting to the Openlane API.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>||
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkerscreatecustomdomainworker"></a>
#### river\.trustcenterworkers\.createcustomdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkerscreatecustomdomainworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkerscreatecustomdomainworkerconfig"></a>
##### river\.trustcenterworkers\.createcustomdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**cloudflareapikey**|`string`||no|
|**validateinterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkersvalidatecustomdomainworker"></a>
#### river\.trustcenterworkers\.validatecustomdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkersvalidatecustomdomainworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkersvalidatecustomdomainworkerconfig"></a>
##### river\.trustcenterworkers\.validatecustomdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**cloudflareapikey**|`string`||no|
|**validateinterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkersdeletecustomdomainworker"></a>
#### river\.trustcenterworkers\.deletecustomdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkersdeletecustomdomainworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkersdeletecustomdomainworkerconfig"></a>
##### river\.trustcenterworkers\.deletecustomdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**cloudflareapikey**|`string`||no|
|**validateinterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkerswatermarkdocworker"></a>
#### river\.trustcenterworkers\.watermarkdocworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkerswatermarkdocworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkerswatermarkdocworkerconfig"></a>
##### river\.trustcenterworkers\.watermarkdocworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkerscreatepirschdomainworker"></a>
#### river\.trustcenterworkers\.createpirschdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkerscreatepirschdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkerscreatepirschdomainworkerconfig"></a>
##### river\.trustcenterworkers\.createpirschdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**pirschclientid**|`string`|||
|**pirschclientsecret**|`string`|||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkersdeletepirschdomainworker"></a>
#### river\.trustcenterworkers\.deletepirschdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkersdeletepirschdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkersdeletepirschdomainworkerconfig"></a>
##### river\.trustcenterworkers\.deletepirschdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**pirschclientid**|`string`|||
|**pirschclientsecret**|`string`|||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkersupdatepirschdomainworker"></a>
#### river\.trustcenterworkers\.updatepirschdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkersupdatepirschdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkersupdatepirschdomainworkerconfig"></a>
##### river\.trustcenterworkers\.updatepirschdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**pirschclientid**|`string`|||
|**pirschclientsecret**|`string`|||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkerscreatepreviewdomainworker"></a>
#### river\.trustcenterworkers\.createpreviewdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkerscreatepreviewdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkerscreatepreviewdomainworkerconfig"></a>
##### river\.trustcenterworkers\.createpreviewdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**cloudflareapikey**|`string`|||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkersdeletepreviewdomainworker"></a>
#### river\.trustcenterworkers\.deletepreviewdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**Config**](#rivertrustcenterworkersdeletepreviewdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "Config": {}
}
```

<a name="rivertrustcenterworkersdeletepreviewdomainworkerconfig"></a>
##### river\.trustcenterworkers\.deletepreviewdomainworker\.Config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**cloudflareapikey**|`string`|||

**Additional Properties:** not allowed  
<a name="rivertrustcenterworkersvalidatepreviewdomainworker"></a>
#### river\.trustcenterworkers\.validatepreviewdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#rivertrustcenterworkersvalidatepreviewdomainworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="rivertrustcenterworkersvalidatepreviewdomainworkerconfig"></a>
##### river\.trustcenterworkers\.validatepreviewdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**cloudflareapikey**|`string`||no|
|**maxsnoozes**|`integer`||yes|
|**snoozeduration**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="rivermetrics"></a>
### river\.metrics: object

MetricsConfig is the configuration for metrics


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enablemetrics**|`boolean`|Enable toggles otel metrics middleware<br/>||
|**metricsdurationunit**|`string`|DurationUnit sets the duration unit for metrics<br/>||
|**enablesemanticmetrics**|`boolean`|EnableSemanticMetrics toggles semantic metrics<br/>||

**Additional Properties:** not allowed  

