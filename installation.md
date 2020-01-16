
Installing The Provider
-----------------------

Full details on installing terraform providers may be found using following links:  [installing a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) and [plugin discovery](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery).

Credentials
-----------
Assuming you have a credentials file (default location `~/.config/turbot/credentials.yml`), the terraform provider will use the default profile credentials if no other credentials are passed.

Credentials may be set (in order of precedence):
 1) using provider config with key and workspace values
 2) using environment variables
 3) using provider config with a profile value (and optionally a path to the credentials file)

#####Provider config - credentials values
```
provider "turbot" {
  access_key     = "xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  secret_key = "xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  workspace         = "https://bananaman-turbot.putney.turbot.io"
}
```
#####Environment Variables
```
export TURBOT_SECRET_KEY=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  
export TURBOT_ACCESS_KEY=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export TURBOT_WORKSPACE=https://bananaman-turbot.putney.turbot.io
```

#####Provider config - profile

```
provider "turbot" {
  profile = bananaman
}
```
