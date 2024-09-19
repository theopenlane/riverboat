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
|[**workers**](#jobqueueworkers)|`object`|||

**Additional Properties:** not allowed  
<a name="jobqueuequeues"></a>
### jobQueue\.queues: array

**Items**

<a name="jobqueueworkers"></a>
### jobQueue\.workers: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**emailWorker**](#jobqueueworkersemailworker)|`object`|||

**Additional Properties:** not allowed  
<a name="jobqueueworkersemailworker"></a>
#### jobQueue\.workers\.emailWorker: object

**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**devMode**|`boolean`|enable dev mode<br/>||
|**testDir**|`string`|the directory to use for dev mode<br/>||
|**token**|`string`|the token to use for the email provider<br/>||
|**fromEmail**|`string`|||

**Additional Properties:** not allowed  

