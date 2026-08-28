# object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**refreshinterval**|`integer`|||
|[**river**](#defsriverconfig)|`object`|Config is the configuration for the river server<br/>||

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

   
<a name="defsriverconfig"></a>
## $defs/river\.Config: object

Config is the configuration for the river server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**databasehost**|`string`|DatabaseHost for connecting to the postgres database<br/>||
|[**queues**](#defsriverqueue)|`array`|||
|[**workers**](#defsriverworkers)|`object`|Workers that will be enabled on the server<br/>||
|[**trustcenterworkers**](#defstrustcenterworkers)|`object`|||
|**defaultmaxretries**|`integer`|DefaultMaxRetries is the maximum number of retries for failed jobs, this can be set differently per job<br/>||
|[**metrics**](#defsriverqueuemetricsconfig)|`object`|MetricsConfig is the configuration for metrics<br/>||

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

   
<a name="defsriverqueue"></a>
### $defs/\[\]river\.Queue: array

**Items**

**Example**

```json
[
    {}
]
```

   
<a name="defsriverworkers"></a>
### $defs/river\.Workers: object

Workers that will be enabled on the server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**openlaneconfig**](#defsjobsopenlaneconfig)|`object`|OpenlaneConfig contains the configuration for connecting to the Openlane API.<br/>||
|[**emailconfig**](#defsjobsemailtemplateconfig)|`object`|EmailTemplateConfig contains configuration that can be shared across workers instead of each worker redefining theirs.<br/>||
|[**emailworker**](#defsjobsemailworker)|`object`|EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider<br/>||
|[**exportcontentworker**](#defsjobsexportcontentworker)|`object`|ExportContentWorker exports the content into csv and makes it downloadable<br/>||
|[**deleteexportcontentworker**](#defsjobsdeleteexportcontentworker)|`object`|DeleteExportContentWorker deletes exports that are older than the configured cutoff duration<br/>||

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

   
<a name="defsjobsdeleteexportcontentworker"></a>
#### $defs/jobs\.DeleteExportContentWorker: object

DeleteExportContentWorker deletes exports that are older than the configured cutoff duration


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#defsjobsdeleteexportworkerconfig)|`object`|DeleteExportWorkerConfig holds the configuration for the delete export worker<br/>|yes|

**Additional Properties:** not allowed   
**Example**

```json
{
    "config": {}
}
```

   
<a name="defsjobsdeleteexportworkerconfig"></a>
##### $defs/jobs\.DeleteExportWorkerConfig: object

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
   
<a name="defsjobsemailtemplateconfig"></a>
#### $defs/jobs\.EmailTemplateConfig: object

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
|[**urls**](#defsemailtemplatesurlconfig)|`object`|||
|**templatespath**|`string`|||

**Additional Properties:** not allowed   
**Example**

```json
{
    "urls": {}
}
```

   
<a name="defsemailtemplatesurlconfig"></a>
##### $defs/emailtemplates\.URLConfig: object

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
   
<a name="defsjobsemailworker"></a>
#### $defs/jobs\.EmailWorker: object

EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#defsjobsemailconfig)|`object`|EmailConfig contains the configuration for the email worker<br/>||

**Additional Properties:** not allowed   
**Example**

```json
{
    "config": {}
}
```

   
<a name="defsjobsemailconfig"></a>
##### $defs/jobs\.EmailConfig: object

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
   
<a name="defsjobsexportcontentworker"></a>
#### $defs/jobs\.ExportContentWorker: object

ExportContentWorker exports the content into csv and makes it downloadable


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#defsjobsexportworkerconfig)|`object`|ExportWorkerConfig configuration for the export content worker<br/>||

**Additional Properties:** not allowed   
**Example**

```json
{
    "config": {}
}
```

   
<a name="defsjobsexportworkerconfig"></a>
##### $defs/jobs\.ExportWorkerConfig: object

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
   
<a name="defsjobsopenlaneconfig"></a>
#### $defs/jobs\.OpenlaneConfig: object

OpenlaneConfig contains the configuration for connecting to the Openlane API.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>||
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>||

**Additional Properties:** not allowed   
   
<a name="defsriverqueuemetricsconfig"></a>
### $defs/riverqueue\.MetricsConfig: object

MetricsConfig is the configuration for metrics


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enablemetrics**|`boolean`|Enable toggles otel metrics middleware<br/>||
|**metricsdurationunit**|`string`|DurationUnit sets the duration unit for metrics<br/>||
|**enablesemanticmetrics**|`boolean`|EnableSemanticMetrics toggles semantic metrics<br/>||

**Additional Properties:** not allowed   
   
<a name="defstrustcenterworkers"></a>
### $defs/trustcenter\.Workers: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**openlaneconfig**](#defsjobsopenlaneconfig)|`object`|OpenlaneConfig contains the configuration for connecting to the Openlane API.<br/>||

**Additional Properties:** not allowed   
**Example**

```json
{
    "openlaneconfig": {}
}
```

   
<a name="defsjobsopenlaneconfig"></a>
#### $defs/jobs\.OpenlaneConfig: object

OpenlaneConfig contains the configuration for connecting to the Openlane API.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>||
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>||

**Additional Properties:** not allowed   

