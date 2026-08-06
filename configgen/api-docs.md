# object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**refreshinterval**|`integer`|||
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
            "emailconfig": {
                "urls": {}
            },
            "emailworker": {
                "config": {}
            },
            "exportcontentworker": {
                "config": {}
            },
            "deleteexportcontentworker": {
                "config": {}
            }
        },
        "trustcenterworkers": {
            "openlaneconfig": {}
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
        "emailconfig": {
            "urls": {}
        },
        "emailworker": {
            "config": {}
        },
        "exportcontentworker": {
            "config": {}
        },
        "deleteexportcontentworker": {
            "config": {}
        }
    },
    "trustcenterworkers": {
        "openlaneconfig": {}
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
|[**emailconfig**](#riverworkersemailconfig)|`object`|EmailTemplateConfig contains configuration that can be shared across workers instead of each worker redefining theirs.<br/>||
|[**emailworker**](#riverworkersemailworker)|`object`|EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider<br/>||
|[**exportcontentworker**](#riverworkersexportcontentworker)|`object`|ExportContentWorker exports the content into csv and makes it downloadable<br/>||
|[**deleteexportcontentworker**](#riverworkersdeleteexportcontentworker)|`object`|DeleteExportContentWorker deletes exports that are older than the configured cutoff duration<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "openlaneconfig": {},
    "emailconfig": {
        "urls": {}
    },
    "emailworker": {
        "config": {}
    },
    "exportcontentworker": {
        "config": {}
    },
    "deleteexportcontentworker": {
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
<a name="riverworkersemailconfig"></a>
#### river\.workers\.emailconfig: object

EmailTemplateConfig contains configuration that can be shared across workers instead of each worker redefining theirs.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**companyname**|`string`|||
|**companyaddress**|`string`|||
|**corporation**|`string`|||
|**year**|`integer`|||
|**fromemail**|`string`|||
|**supportemail**|`string`|||
|**questionnaireemail**|`string`|||
|**logourl**|`string`|||
|[**urls**](#riverworkersemailconfigurls)|`object`|||
|**templatespath**|`string`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "urls": {}
}
```

<a name="riverworkersemailconfigurls"></a>
##### river\.workers\.emailconfig\.urls: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**root**|`string`|||
|**product**|`string`|||
|**docs**|`string`|||
|**verify**|`string`|||
|**invite**|`string`|||
|**reset**|`string`|||
|**verifysubscriber**|`string`|||
|**verifybilling**|`string`|||
|**billing**|`string`|||
|**questionnaire**|`string`|||

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
|**enabled**|`boolean`|Enabled indicates if this job is enabled in the server<br/>||
|**maxzipsize**|`integer`|the maximum allowed size in bytes for a zip archive export<br/>||
|**cloudflareaccountid**|`string`|the cloudflare account id used for browser rendering pdf generation<br/>||
|**cloudflareapikey**|`string`|the cloudflare api key used for browser rendering pdf generation<br/>||
|**maxsnoozes**|`integer`|MaxSnoozes is the maximum number of times to snooze the job before giving up<br/>||
|**snoozeduration**|`integer`|SnoozeDuration is the duration to snooze between PDF render retries<br/>||

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
<a name="rivertrustcenterworkers"></a>
### river\.trustcenterworkers: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**openlaneconfig**](#rivertrustcenterworkersopenlaneconfig)|`object`|OpenlaneConfig contains the configuration for connecting to the Openlane API.<br/>||

**Additional Properties:** not allowed  
**Example**

```json
{
    "openlaneconfig": {}
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

