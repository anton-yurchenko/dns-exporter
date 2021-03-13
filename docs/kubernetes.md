# Kubernetes

Create a ConfigMap / Secret with your unique parameters:

<details><summary>ConfigMap</summary>

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: dns-exporter
  namespace: default
  labels:
    app: dns-exporter
data:
  delay: "1"
  git.remote: "true"
  git.url: "https://github.com/user/dns-archive.git"
  git.branch: "master"
  cloudflare.enabled: "true"
  route53.enabled: "true"
  route53.region: "us-west-2"
```

</details>
<br />

<details><summary>Secret</summary>

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: dns-exporter
  namespace: default
  labels:
    app: dns-exporter
type: Opaque
data:
  git.user: bWFjaGluZS11c2Vy                                                      # machine-user
  git.email: bWFjaGluZS11c2VyQGRvbWFpbi5jb20=                                     # machine-user@domain.com
  git.token: MGFiMTIzNGM1Njc4OTAxMmQzZWY0NWc2Nzg5aDAxMjNpajQ1Njc4OQ==             # 0ab1234c56789012d3ef45g6789h0123ij456789
  cloudflare.email: b3duZXJAZG9tYWluLmNvbQ==                                      # owner@domain.com
  cloudflare.token: MXp4OTIzNGM1Njc4OTAxMmQzZWY0NWc2Nzg5aDAxMjNpajQ1Njc4OQ==      # 1zx9234c56789012d3ef45g6789h0123ij456789
```

</details>
<br />

- Remember to store values in base64 form with no new lines (`echo -n <token> | base64`)

Spin up a pod manually or create a cronjob:

<details><summary>POD</summary>

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dns-exporter
  namespace: default
  labels:
    app: dns-exporter
spec:
  containers:
    - name: dns-exporter
      image: antonyurchenko/dns-exporter:latest
      env:
        - name: DELAY
          valueFrom:
            configMapKeyRef:
              name: dns-exporter
              key: delay
        - name: GIT_REMOTE_ENABLED
          valueFrom:
            configMapKeyRef:
              name: dns-exporter
              key: git.remote
        - name: GIT_URL
          valueFrom:
            configMapKeyRef:
              name: dns-exporter
              key: git.url
        - name: GIT_BRANCH
          valueFrom:
            configMapKeyRef:
              name: dns-exporter
              key: git.branch
        - name: GIT_USER
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: git.user
        - name: GIT_EMAIL
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: git.email
        - name: GIT_TOKEN
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: git.token
        - name: CLOUDFLARE_ENABLED
          valueFrom:
            configMapKeyRef:
              name: dns-exporter
              key: cloudflare.enabled
        - name: CLOUDFLARE_EMAIL
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: cloudflare.email
        - name: CLOUDFLARE_TOKEN
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: cloudflare.token
        - name: ROUTE53_ENABLED
          valueFrom:
            configMapKeyRef:
              name: dns-exporter
              key: route53.enabled
        - name: AWS_REGION
          valueFrom:
            configMapKeyRef:
              name: dns-exporter
              key: route53.region
```

</details>
<br />

<details><summary>CronJob</summary>

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: dns-exporter
  namespace: default
  labels:
    app: dns-exporter
spec:
  schedule: "30 6 1 * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: dns-exporter
            image: antonyurchenko/dns-exporter:latest
            env:
              - name: DELAY
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: delay
              - name: GIT_REMOTE_ENABLED
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: git.remote
              - name: GIT_URL
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: git.url
              - name: GIT_BRANCH
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: git.branch
              - name: GIT_USER
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: git.user
              - name: GIT_EMAIL
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: git.email
              - name: GIT_TOKEN
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: git.token
              - name: CLOUDFLARE_ENABLED
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: cloudflare.enabled
              - name: CLOUDFLARE_EMAIL
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: cloudflare.email
              - name: CLOUDFLARE_TOKEN
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: cloudflare.token
              - name: ROUTE53_ENABLED
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: route53.enabled
              - name: AWS_REGION
                valueFrom:
                  secretKeyRef:
                    name: dns-exporter
                    key: route53.region
          restartPolicy: OnFailure
```

</details>

### Local Archive Only

You will have to create a persistent volume and attach it as `/opt/data`. 
Here is an example of changes you might be required to apply to your manifests:

<details><summary>POD</summary>

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dns-exporter
  namespace: default
  labels:
    app: dns-exporter
spec:
  containers:
    - name: dns-exporter
      image: antonyurchenko/dns-exporter:latest
      volumeMounts:
        - name: archive
          mountPath: "/opt/data"
      env:
        - name: DELAY
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: delay
        - name: GIT_REMOTE_ENABLED
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: git.remote
        - name: CLOUDFLARE_ENABLED
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: cloudflare.enabled
        - name: CLOUDFLARE_EMAIL
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: cloudflare.email
        - name: CLOUDFLARE_TOKEN
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: cloudflare.token
        - name: ROUTE53_ENABLED
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: route53.enabled
        - name: AWS_REGION
          valueFrom:
            secretKeyRef:
              name: dns-exporter
              key: route53.region
  volumes:
    - name: archive
      persistentVolumeClaim:
        claimName: local-dns-archive
```

</details>
<br />
