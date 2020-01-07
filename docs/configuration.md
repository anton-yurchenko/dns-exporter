# Configuration
**DNS-EXPORTER** configuration is managed via the following environmental variables:
- `DELAY`: Providers API calls delays, applied per provider, default 1(sec).
- `GIT_REMOTE_ENABLED`: Set to `"true"` if you want to push exported files to remote git repository
- `GIT_URL`: Git URL in form of HTTPS. For example: `"https://github.com/user/dns-archive.git"`
- `GIT_BRANCH`: If remote git is enabled, you may choose which branch to clone/pull/push
- `GIT_USER`: Committer username, also used for authentication against remote repository 
- `GIT_EMAIL`: Committer email.
- `GIT_TOKEN`: Committer token, used for authentication against remote repository
- `CLOUDFLARE_ENABLED`: set to `"true"` to enable that provider
- `CLOUDFLARE_EMAIL`: Cloudflare user email address, required for authentication
- `CLOUDFLARE_TOKEN`: Global API Key, required for authentication
- `ROUTE53_ENABLED`: Set to `"true"` to enable that provider
- `AWS_REGION`: Substitute your desired AWS Region

In addition to that, enabling **AWS Route53**, it is expected that AWS authentication is pre-configured by:
- attaching AWS IAM Role
- settings additional environmental variables `AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY`
- mounting `/home/app/.aws` directory with `credentials / config` files
