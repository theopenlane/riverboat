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
            },
            "slackworker": {
                "config": {}
            },
            "organizationdeletionreminderworker": {
                "config": {
                    "email": {
                        "config": {
                            "urls": {}
                        }
                    }
                }
            },
            "organizationdeletionworker": {
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
        },
        "slackworker": {
            "config": {}
        },
        "organizationdeletionreminderworker": {
            "config": {
                "email": {
                    "config": {
                        "urls": {}
                    }
                }
            }
        },
        "organizationdeletionworker": {
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
|[**slackworker**](#riverworkersslackworker)|`object`|SlackWorker sends messages to Slack.<br/>||
|[**organizationdeletionreminderworker**](#riverworkersorganizationdeletionreminderworker)|`object`|OrganizationPaymentReminderWorker fetches organizations for payment reminder processing.<br/>||
|[**organizationdeletionworker**](#riverworkersorganizationdeletionworker)|`object`|OrganizationDeleteWorker deletes organizations in Openlane.<br/>||

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
    },
    "slackworker": {
        "config": {}
    },
    "organizationdeletionreminderworker": {
        "config": {
            "email": {
                "config": {
                    "urls": {}
                }
            }
        }
    },
    "organizationdeletionworker": {
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
|**enabled**|`boolean`|||
|**maxzipsize**|`integer`|the maximum allowed size in bytes for a zip archive export<br/>||

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
<a name="riverworkersorganizationdeletionreminderworker"></a>
#### river\.workers\.organizationdeletionreminderworker: object

OrganizationPaymentReminderWorker fetches organizations for payment reminder processing.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersorganizationdeletionreminderworkerconfig)|`object`|OrganizationPaymentReminderConfig contains the configuration for the organization payment reminder worker.<br/>|yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {
        "email": {
            "config": {
                "urls": {}
            }
        }
    }
}
```

<a name="riverworkersorganizationdeletionreminderworkerconfig"></a>
##### river\.workers\.organizationdeletionreminderworker\.config: object

OrganizationPaymentReminderConfig contains the configuration for the organization payment reminder worker.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>|no|
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>|no|
|**paymentmethodinterval**|`integer`|PaymentMethodInterval is the amount of days an org must have a payment method attached or else it will be earmarked for deletion<br/>This is after org creation. So if an org is created 7 days ago and this is set to 6 days, the org will be marked<br/>as pending deletion. But if set to say 8 days, nothing happens<br/>|yes|
|**deletiondays**|`integer`|DeletionDays is the number of days an org has before the deletion actually occurs. Once an org is earmarked for<br/>deletion, we do not delete immediately, instead we send them an email and update "pending_deletion_at". SO if<br/>DeletionDays is set to 30, the org will be deleted at in 30 days ( pending_deletion_at set to today + 30 days)<br/>|yes|
|**enabled**|`boolean`|Enabled is used to determine if to register this worker or not<br/>|no|
|**dryrun**|`boolean`|if true<br/>|no|
|[**email**](#riverworkersorganizationdeletionreminderworkerconfigemail)|`object`||no|

**Additional Properties:** not allowed  
**Example**

```json
{
    "email": {
        "config": {
            "urls": {}
        }
    }
}
```

<a name="riverworkersorganizationdeletionreminderworkerconfigemail"></a>
###### river\.workers\.organizationdeletionreminderworker\.config\.email: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|||
|[**config**](#riverworkersorganizationdeletionreminderworkerconfigemailconfig)|`object`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {
        "urls": {}
    }
}
```

<a name="riverworkersorganizationdeletionreminderworkerconfigemailconfig"></a>
####### river\.workers\.organizationdeletionreminderworker\.config\.email\.config: object

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
|[**urls**](#riverworkersorganizationdeletionreminderworkerconfigemailconfigurls)|`object`|||
|**templatespath**|`string`|||

**Additional Properties:** not allowed  
**Example**

```json
{
    "urls": {}
}
```

<a name="riverworkersorganizationdeletionreminderworkerconfigemailconfigurls"></a>
######## river\.workers\.organizationdeletionreminderworker\.config\.email\.config\.urls: object

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
<a name="riverworkersorganizationdeletionworker"></a>
#### river\.workers\.organizationdeletionworker: object

OrganizationDeleteWorker deletes organizations in Openlane.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersorganizationdeletionworkerconfig)|`object`|OrganizationDeleteConfig contains the configuration for the organization deletion worker.<br/>|yes|

**Additional Properties:** not allowed  
**Example**

```json
{
    "config": {}
}
```

<a name="riverworkersorganizationdeletionworkerconfig"></a>
##### river\.workers\.organizationdeletionworker\.config: object

OrganizationDeleteConfig contains the configuration for the organization deletion worker.


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**openlaneapihost**|`string`|OpenlaneAPIHost is the host URL for the Openlane API<br/>|no|
|**openlaneapitoken**|`string`|OpenlaneAPIToken is the API token for authenticating with the Openlane API<br/>|no|
|**runinterval**|`integer`||yes|
|**maxdeletesperrun**|`integer`||yes|
|**enabled**|`boolean`||no|

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

