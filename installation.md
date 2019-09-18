
Installing The Provider
-----------------------

The simplest way to install the provider is to run the following command from the directory containing your Terraform config: 

```
terraform init -plugin-dir=<PATH>
``` 
where `<PATH>` is the directory containing the provider file `terraform-provider-turbot`.    


Full details on installing terraform providers may be found using following links:  [installing a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) and [plugin discovery](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery).

Credentials
-----------
Credentials may be set either using environment variables, or using provider config in the Terraform overrides file.

#####Environment Variables
```
export TURBOT_SECRET_ACCESS_KEY=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  
export TURBOT_ACCESS_KEY_ID=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export TURBOT_WORKSPACE=https://bananaman-turbot.putney.turbot.io
```
#####Override file
```
provider "turbot" {
  access_key_id     = "xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  secret_access_key = "xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  workspace         = "https://bananaman-turbot.putney.turbot.io"
}
```
