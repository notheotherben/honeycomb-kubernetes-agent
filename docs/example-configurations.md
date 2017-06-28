Here are some example configurations for the Honeycomb agent:


Parse logs from pods labelled with `app: nginx`:
```
---
writekey: "YOUR_HONEYCOMB_WRITEKEY_HERE"
parsers:
  - labelSelector: app=nginx
    parser: nginx

    processors:
    - request_shape:
      field: request
```

Send logs from different services to different datasets:
```
writekey: "YOUR_HONEYCOMB_WRITEKEY_HERE"
parsers:
  - labelSelector: "app=nginx"
    parser: nginx
    dataset: nginx-kubernetes

  - labelSelector: "app=frontend-web"
    parser: json
    dataset: frontend
```


Sample events from a service (only send one in 20 events), and drop the
`user_email` field from all events:
```
writekey: "YOUR_HONEYCOMB_WRITEKEY_HERE"
parsers:
  - labelSelector: "app=frontend-web"
    parser: json
    dataset: frontend

    processors:
      - sample:
        type: static
        rate: 20
      - drop_field:
        field: user_email
```

Only process logs from the `sidecar` container in a multi-container pod:
```
---
writekey: "YOUR_HONEYCOMB_WRITEKEY_HERE"
parsers:
  - labelSelector: "app=frontend-web"
    containerName: sidecar
    parser: json
    dataset: frontend
```
