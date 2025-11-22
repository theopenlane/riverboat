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
            "createcustomdomainworker": {
                "config": {}
            },
            "validatecustomdomainworker": {
                "config": {}
            },
            "deletecustomdomainworker": {
                "config": {}
            },
            "exportcontentworker": {
                "config": {}
            },
            "deleteexportcontentworker": {
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
            "slackworker": {
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
        "createcustomdomainworker": {
            "config": {}
        },
        "validatecustomdomainworker": {
            "config": {}
        },
        "deletecustomdomainworker": {
            "config": {}
        },
        "exportcontentworker": {
            "config": {}
        },
        "deleteexportcontentworker": {
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
        "slackworker": {
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
|[**openlaneconfig**](#riverworkersopenlaneconfig)|`object`|||
|[**emailworker**](#riverworkersemailworker)|`object`|EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider<br/>||
|[**databaseworker**](#riverworkersdatabaseworker)|`object`|DatabaseWorker is a worker to create a dedicated database for an organization<br/>||
|[**createcustomdomainworker**](#riverworkerscreatecustomdomainworker)|`object`|||
|[**validatecustomdomainworker**](#riverworkersvalidatecustomdomainworker)|`object`|||
|[**deletecustomdomainworker**](#riverworkersdeletecustomdomainworker)|`object`|||
|[**exportcontentworker**](#riverworkersexportcontentworker)|`object`|||
|[**deleteexportcontentworker**](#riverworkersdeleteexportcontentworker)|`object`|||
|[**watermarkdocworker**](#riverworkerswatermarkdocworker)|`object`|||
|[**createpirschdomainworker**](#riverworkerscreatepirschdomainworker)|`object`|||
|[**deletepirschdomainworker**](#riverworkersdeletepirschdomainworker)|`object`|||
|[**updatepirschdomainworker**](#riverworkersupdatepirschdomainworker)|`object`|||
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
    "createcustomdomainworker": {
        "config": {}
    },
    "validatecustomdomainworker": {
        "config": {}
    },
    "deletecustomdomainworker": {
        "config": {}
    },
    "exportcontentworker": {
        "config": {}
    },
    "deleteexportcontentworker": {
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
    "slackworker": {
        "config": {}
    }
}
```

<a name="riverworkersopenlaneconfig"></a>
#### river\.workers\.openlaneconfig: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||

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
<a name="riverworkerscreatecustomdomainworker"></a>
#### river\.workers\.createcustomdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkerscreatecustomdomainworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkerscreatecustomdomainworkerconfig"></a>
##### river\.workers\.createcustomdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**cloudflareapikey**|`string`||no|
|**validateinterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="riverworkersvalidatecustomdomainworker"></a>
#### river\.workers\.validatecustomdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersvalidatecustomdomainworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersvalidatecustomdomainworkerconfig"></a>
##### river\.workers\.validatecustomdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**cloudflareapikey**|`string`||no|
|**validateinterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="riverworkersdeletecustomdomainworker"></a>
#### river\.workers\.deletecustomdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersdeletecustomdomainworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersdeletecustomdomainworkerconfig"></a>
##### river\.workers\.deletecustomdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**cloudflareapikey**|`string`||no|
|**validateinterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="riverworkersexportcontentworker"></a>
#### river\.workers\.exportcontentworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersexportcontentworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersexportcontentworkerconfig"></a>
##### river\.workers\.exportcontentworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||

**Additional Properties:** not allowed  
<a name="riverworkersdeleteexportcontentworker"></a>
#### river\.workers\.deleteexportcontentworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersdeleteexportcontentworkerconfig)|`object`||yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersdeleteexportcontentworkerconfig"></a>
##### river\.workers\.deleteexportcontentworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`||no|
|**openlaneapitoken**|`string`||no|
|**enabled**|`boolean`||no|
|**interval**|`integer`||yes|
|**cutoffduration**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="riverworkerswatermarkdocworker"></a>
#### river\.workers\.watermarkdocworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkerswatermarkdocworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkerswatermarkdocworkerconfig"></a>
##### river\.workers\.watermarkdocworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||

**Additional Properties:** not allowed  
<a name="riverworkerscreatepirschdomainworker"></a>
#### river\.workers\.createpirschdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkerscreatepirschdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkerscreatepirschdomainworkerconfig"></a>
##### river\.workers\.createpirschdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**pirschclientid**|`string`|||
|**pirschclientsecret**|`string`|||

**Additional Properties:** not allowed  
<a name="riverworkersdeletepirschdomainworker"></a>
#### river\.workers\.deletepirschdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersdeletepirschdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersdeletepirschdomainworkerconfig"></a>
##### river\.workers\.deletepirschdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**pirschclientid**|`string`|||
|**pirschclientsecret**|`string`|||

**Additional Properties:** not allowed  
<a name="riverworkersupdatepirschdomainworker"></a>
#### river\.workers\.updatepirschdomainworker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersupdatepirschdomainworkerconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersupdatepirschdomainworkerconfig"></a>
##### river\.workers\.updatepirschdomainworker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|||
|**openlaneapitoken**|`string`|||
|**enabled**|`boolean`|||
|**pirschclientid**|`string`|||
|**pirschclientsecret**|`string`|||

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

