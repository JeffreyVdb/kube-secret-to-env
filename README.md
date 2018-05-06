# kube-secret-to-env
A tool to export a kubernetes secret to an environment variable file

# Installing

```bash
go get github.com/JeffreyVdb/kube-secret-to-env
```
## Examples

In dotenv style output:

```bash
kubectl get secrets api-secrets -o json | kube-secret-to-env -type env
3RD_PARTY_SERVICE_USER=admin
3RD_PARTY_SERVICE_PASSWORD=secure_password
```

In shell output:

```bash
kubectl get secrets api-secrets -o json | kube-secret-to-env -type shell
export 3RD_PARTY_SERVICE_USER="admin"
export 3RD_PARTY_SERVICE_PASSWORD="secure_password"
```
