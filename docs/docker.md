# Docker

### Manual

Pass all parameters to `docker run` command, for example: 

```
docker run -it \
    -e GIT_REMOTE_ENABLED=false \
    -e CLOUDFLARE_ENABLED=true \
    -e CLOUDFLARE_EMAIL=owner@domain.com \
    -e CLOUDFLARE_TOKEN=1zx9234c56789012d3ef45g6789h0123ij456789 \
    antonyurchenko/dns-exporter:latest
```

### Docker-Compose

Create a **.env** file with your unique parameters:

<details><summary>.env</summary>

```
DELAY=1
GIT_REMOTE_ENABLED=true
GIT_URL=https://github.com/user/dns-archive.git
GIT_BRANCH=master
GIT_USER=machine-user
GIT_EMAIL=machine-user@domain.com
GIT_TOKEN=0ab1234c56789012d3ef45g6789h0123ij456789
CLOUDFLARE_ENABLED=true
CLOUDFLARE_EMAIL=owner@domain.com
CLOUDFLARE_TOKEN=1zx9234c56789012d3ef45g6789h0123ij456789
ROUTE53_ENABLED=true
AWS_REGION=us-west-2
```

</details>
<br />

Create a **docker-compose.yaml**:

<details><summary>docker-compose.yaml</summary>

```yaml
version: "3"

services:
  dns-exporter:
    image: antonyurchenko/dns-exporter:latest
    env_file:
      - .env
```

</details>
<br />

Execute with `docker-compose -f docker-compose.yaml up`
