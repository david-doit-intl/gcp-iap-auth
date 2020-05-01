# gcp-iap-auth

`gcp-iap-auth` is a simple server implementation and package in
[Go](http://golang.org) for helping you secure your web apps running on GCP
behind a
[Google Cloud Platform's IAP (Identity-Aware Proxy)](https://cloud.google.com/iap/docs/) by validating IAP signed headers in the requests.
[Ldap Server](https://support.google.com/a/answer/9089736?hl=en)

## Why

Validating signed headers helps you protect your app from the following kinds of risks:

- IAP is accidentally disabled;
- Misconfigured firewalls;
- Access from within the project.
- You wish to connect and validate users via ldap

## Using it with Kubernetes

### As a reverse proxy

A simple way to use it with
[Kubernetes](https://github.com/kubernetes/kubernetes) and without any other
dependencies is to run it as a reverse proxy that validates and forwards
requests to a backend server.

```yaml
      - name: gcp-iap-auth
        image: gcp-iap-auth:latest
        env:
        - name: GCP_IAP_AUTH_AUDIENCES
          value: "YOUR_AUDIENCE"
        - name: GCP_IAP_AUTH_LISTEN_PORT
          value: "1080"
        - name: GCP_IAP_AUTH_BACKEND
          value: "http://YOUR_BACKEND_SERVER"
        ports:
        - name: proxy
          containerPort: 1080
        readinessProbe:
          httpGet:
            path: /healthz
            scheme: HTTP
            port: proxy
          periodSeconds: 1
          timeoutSeconds: 1
          successThreshold: 1
          failureThreshold: 10
        livenessProbe:
          httpGet:
            path: /healthz
            scheme: HTTP
            port: proxy
          timeoutSeconds: 5
          initialDelaySeconds: 10
```


### Notes

To use HTTPS just make sure:
- You set up `GCP_IAP_AUTH_TLS_CERT=/path/to/tls_cert_file` and `GCP_IAP_AUTH_TLS_KEY=/path/to/tls_key_file` environment variables.
- You set up volumes for [secrets](https://kubernetes.io/docs/concepts/configuration/secret/) in Kubernetes so it knows where to find them.
- Change the scheme in readiness and liveness probes to `HTTPS`.
- Adjust your nginx.conf as necessary to proxy pass the auth requests to gcp-iap-auth as HTTPS.

