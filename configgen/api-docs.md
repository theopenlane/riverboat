# object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**refreshInterval**|`integer`|||
|[**jobQueue**](#jobqueue)|`object`|Config is the configuration for the river server<br/>||

**Additional Properties:** not allowed  
<a name="jobqueue"></a>
## jobQueue: object

Config is the configuration for the river server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**databaseHost**|`string`|DatabaseHost for connecting to the postgres database<br/>||
|[**queues**](#jobqueuequeues)|`array`|||
|[**workers**](#jobqueueworkers)|`object`|Workers that will be enabled on the server<br/>||

**Additional Properties:** not allowed  
<a name="jobqueuequeues"></a>
### jobQueue\.queues: array

**Items**

<a name="jobqueueworkers"></a>
### jobQueue\.workers: object

Workers that will be enabled on the server


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**emailWorker**](#jobqueueworkersemailworker)|`object`|||
|[**databaseWorker**](#jobqueueworkersdatabaseworker)|`object`|||

**Additional Properties:** not allowed  
<a name="jobqueueworkersemailworker"></a>
#### jobQueue\.workers\.emailWorker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**email**](#jobqueueworkersemailworkeremail)|`object`|||

**Additional Properties:** not allowed  
<a name="jobqueueworkersemailworkeremail"></a>
##### jobQueue\.workers\.emailWorker\.email: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**devMode**|`boolean`|enable dev mode<br/>||
|**testDir**|`string`|the directory to use for dev mode<br/>||
|**token**|`string`|the token to use for the email provider<br/>||
|**fromEmail**|`string`|||

**Additional Properties:** not allowed  
<a name="jobqueueworkersdatabaseworker"></a>
#### jobQueue\.workers\.databaseWorker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#jobqueueworkersdatabaseworkerconfig)|`object`|||

**Additional Properties:** not allowed  
<a name="jobqueueworkersdatabaseworkerconfig"></a>
##### jobQueue\.workers\.databaseWorker\.config: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**enabled**|`boolean`|Enable the dbx client<br/>||
|**baseUrl**|`string`|Base URL for the dbx service<br/>||
|**endpoint**|`string`|Endpoint for the graphql api<br/>||
|**debug**|`boolean`|Enable debug mode<br/>||

**Additional Properties:** not allowed  

