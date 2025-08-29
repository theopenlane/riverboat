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
            "emailWorker": {
                "config": {}
            },
            "databaseWorker": {
                "config": {}
            },
            "createCustomDomainWorker": {
                "config": {}
            },
            "validateCustomDomainWorker": {
                "config": {}
            },
            "deleteCustomDomainWorker": {
                "config": {}
            },
            "exportContentWorker": {
                "config": {}
            },
            "deleteExportContentWorker": {
                "config": {}
            }
        }
    }
}
```

<a name="river"></a>
## river: object

Config is the configuration for the river server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**databaseHost**|`string`|DatabaseHost for connecting to the postgres database<br/>||
|[**queues**](#riverqueues)|`array`|||
|[**workers**](#riverworkers)|`object`|Workers that will be enabled on the server<br/>||
|**defaultMaxRetries**|`integer`|DefaultMaxRetries is the maximum number of retries for failed jobs, this can be set differently per job<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "queues": [
        {}
    ],
    "workers": {
        "emailWorker": {
            "config": {}
        },
        "databaseWorker": {
            "config": {}
        },
        "createCustomDomainWorker": {
            "config": {}
        },
        "validateCustomDomainWorker": {
            "config": {}
        },
        "deleteCustomDomainWorker": {
            "config": {}
        },
        "exportContentWorker": {
            "config": {}
        },
        "deleteExportContentWorker": {
            "config": {}
        }
    }
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
|[**emailWorker**](#riverworkersemailworker)|`object`|EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider<br/>||
|[**databaseWorker**](#riverworkersdatabaseworker)|`object`|DatabaseWorker is a worker to create a dedicated database for an organization<br/>||
|[**createCustomDomainWorker**](#riverworkerscreatecustomdomainworker)|`object`|||
|[**validateCustomDomainWorker**](#riverworkersvalidatecustomdomainworker)|`object`|||
|[**deleteCustomDomainWorker**](#riverworkersdeletecustomdomainworker)|`object`|||
|[**exportContentWorker**](#riverworkersexportcontentworker)|`object`|||
|[**deleteExportContentWorker**](#riverworkersdeleteexportcontentworker)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "emailWorker": {
        "config": {}
    },
    "databaseWorker": {
        "config": {}
    },
    "createCustomDomainWorker": {
        "config": {}
    },
    "validateCustomDomainWorker": {
        "config": {}
    },
    "deleteCustomDomainWorker": {
        "config": {}
    },
    "exportContentWorker": {
        "config": {}
    },
    "deleteExportContentWorker": {
        "config": {}
    }
}
```

<a name="riverworkersemailworker"></a>
#### river\.workers\.emailWorker: object

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
##### river\.workers\.emailWorker\.config: object

EmailConfig contains the configuration for the email worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|enable or disable the email worker<br/>||
|**devMode**|`boolean`|enable dev mode<br/>||
|**testDir**|`string`|the directory to use for dev mode<br/>||
|**token**|`string`|the token to use for the email provider<br/>||
|**fromEmail**|`string`|FromEmail is the email address to use as the sender<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersdatabaseworker"></a>
#### river\.workers\.databaseWorker: object

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
##### river\.workers\.databaseWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|Enable the dbx client<br/>||
|**baseUrl**|`string`|Base URL for the dbx service<br/>||
|**endpoint**|`string`|Endpoint for the graphql api<br/>||
|**debug**|`boolean`|Enable debug mode<br/>||

**Additional Properties:** not allowed  
<a name="riverworkerscreatecustomdomainworker"></a>
#### river\.workers\.createCustomDomainWorker: object

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
##### river\.workers\.createCustomDomainWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`||no|
|**cloudflareApiKey**|`string`||no|
|**openlaneAPIHost**|`string`||no|
|**openlaneAPIToken**|`string`||no|
|**databaseHost**|`string`||no|
|**validateInterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="riverworkersvalidatecustomdomainworker"></a>
#### river\.workers\.validateCustomDomainWorker: object

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
##### river\.workers\.validateCustomDomainWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`||no|
|**cloudflareApiKey**|`string`||no|
|**openlaneAPIHost**|`string`||no|
|**openlaneAPIToken**|`string`||no|
|**databaseHost**|`string`||no|
|**validateInterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="riverworkersdeletecustomdomainworker"></a>
#### river\.workers\.deleteCustomDomainWorker: object

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
##### river\.workers\.deleteCustomDomainWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`||no|
|**cloudflareApiKey**|`string`||no|
|**openlaneAPIHost**|`string`||no|
|**openlaneAPIToken**|`string`||no|
|**databaseHost**|`string`||no|
|**validateInterval**|`integer`||yes|

**Additional Properties:** not allowed  
<a name="riverworkersexportcontentworker"></a>
#### river\.workers\.exportContentWorker: object

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
##### river\.workers\.exportContentWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|||
|**openlaneAPIHost**|`string`|||
|**openlaneAPIToken**|`string`|||

**Additional Properties:** not allowed  
<a name="riverworkersdeleteexportcontentworker"></a>
#### river\.workers\.deleteExportContentWorker: object

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
##### river\.workers\.deleteExportContentWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`||no|
|**interval**|`integer`||yes|
|**openlaneAPIHost**|`string`||no|
|**openlaneAPIToken**|`string`||no|
|**cutoffDuration**|`integer`||yes|

**Additional Properties:** not allowed  

