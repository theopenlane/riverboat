# object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**refreshInterval**|`integer`|||
|[**river**](#river)|`object`|Config is the configuration for the river server<br/>||

**Additional Properties:** not allowed  
<a name="river"></a>
## river: object

Config is the configuration for the river server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**databaseHost**|`string`|DatabaseHost for connecting to the postgres database<br/>||
|[**queues**](#riverqueues)|`array`|||
|[**workers**](#riverworkers)|`object`|Workers that will be enabled on the server<br/>||

**Additional Properties:** not allowed  
<a name="riverqueues"></a>
### river\.queues: array

**Items**

<a name="riverworkers"></a>
### river\.workers: object

Workers that will be enabled on the server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**emailWorker**](#riverworkersemailworker)|`object`|||
|[**databaseWorker**](#riverworkersdatabaseworker)|`object`|||
|[**onboardingWorker**](#riverworkersonboardingworker)|`object`|||

**Additional Properties:** not allowed  
<a name="riverworkersemailworker"></a>
#### river\.workers\.emailWorker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersemailworkerconfig)|`object`|||

**Additional Properties:** not allowed  
<a name="riverworkersemailworkerconfig"></a>
##### river\.workers\.emailWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**devMode**|`boolean`|enable dev mode<br/>||
|**testDir**|`string`|the directory to use for dev mode<br/>||
|**token**|`string`|the token to use for the email provider<br/>||
|**fromEmail**|`string`|||

**Additional Properties:** not allowed  
<a name="riverworkersdatabaseworker"></a>
#### river\.workers\.databaseWorker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersdatabaseworkerconfig)|`object`|||

**Additional Properties:** not allowed  
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
<a name="riverworkersonboardingworker"></a>
#### river\.workers\.onboardingWorker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersonboardingworkerconfig)|`object`|||

**Additional Properties:** not allowed  
<a name="riverworkersonboardingworkerconfig"></a>
##### river\.workers\.onboardingWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**StarterTasks**](#riverworkersonboardingworkerconfigstartertasks)|`array`|||
|**APIBaseURL**|`string`|the base URL for the Openlane API<br/>Format: `"uri"`<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersonboardingworkerconfigstartertasks"></a>
###### river\.workers\.onboardingWorker\.config\.StarterTasks: array

**Items**


