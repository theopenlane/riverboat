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
|[**emailWorker**](#riverworkersemailworker)|`object`|EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider<br/>||
|[**databaseWorker**](#riverworkersdatabaseworker)|`object`|DatabaseWorker is a worker to create a dedicated database for an organization<br/>||
|[**createCustomDomainWorker**](#riverworkerscreatecustomdomainworker)|`object`|CreateCustomDomainWorker creates a custom hostname in cloudflare, and creates and updates the records in our system<br/>||
|[**validateCustomDomainWorker**](#riverworkersvalidatecustomdomainworker)|`object`|ValidateCustomDomainWorker checks cloudflare custom domain(s), and updates the status in our system<br/>||
|[**deleteCustomDomainWorker**](#riverworkersdeletecustomdomainworker)|`object`|DeleteCustomDomainWorker delete the custom hostname from cloudflare and updates the records in our system<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersemailworker"></a>
#### river\.workers\.emailWorker: object

EmailWorker is a worker to send emails using the resend email provider the config defaults to dev mode, which will write the email to a file using the mock provider a token is required to send emails using the actual resend provider


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**config**](#riverworkersemailworkerconfig)|`object`|EmailConfig contains the configuration for the email worker<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersemailworkerconfig"></a>
##### river\.workers\.emailWorker\.config: object

EmailConfig contains the configuration for the email worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
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

CreateCustomDomainWorker creates a custom hostname in cloudflare, and creates and updates the records in our system


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**Config**](#riverworkerscreatecustomdomainworkerconfig)|`object`|CreateCustomDomainConfig contains the configuration for the worker<br/>||

**Additional Properties:** not allowed  
<a name="riverworkerscreatecustomdomainworkerconfig"></a>
##### river\.workers\.createCustomDomainWorker\.Config: object

CreateCustomDomainConfig contains the configuration for the worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**cloudflareApiKey**|`string`|||

**Additional Properties:** not allowed  
<a name="riverworkersvalidatecustomdomainworker"></a>
#### river\.workers\.validateCustomDomainWorker: object

ValidateCustomDomainWorker checks cloudflare custom domain(s), and updates the status in our system


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**Config**](#riverworkersvalidatecustomdomainworkerconfig)|`object`|ValidateCustomDomainConfig contains the configuration for the worker<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersvalidatecustomdomainworkerconfig"></a>
##### river\.workers\.validateCustomDomainWorker\.Config: object

ValidateCustomDomainConfig contains the configuration for the worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**cloudflareApiKey**|`string`|||

**Additional Properties:** not allowed  
<a name="riverworkersdeletecustomdomainworker"></a>
#### river\.workers\.deleteCustomDomainWorker: object

DeleteCustomDomainWorker delete the custom hostname from cloudflare and updates the records in our system


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|[**Config**](#riverworkersdeletecustomdomainworkerconfig)|`object`|DeleteCustomDomainConfig contains the configuration for the example worker<br/>||

**Additional Properties:** not allowed  
<a name="riverworkersdeletecustomdomainworkerconfig"></a>
##### river\.workers\.deleteCustomDomainWorker\.Config: object

DeleteCustomDomainConfig contains the configuration for the example worker


**Properties**

|Name|Type|Description|Required|
|----|----|-----------|--------|
|**cloudflareApiKey**|`string`|||

**Additional Properties:** not allowed  

