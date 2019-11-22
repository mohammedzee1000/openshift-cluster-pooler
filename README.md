# openshift-cluster-pool
Code base for openshift cluster pool

## Try It (With minishift) - Developer

### Prepare

1. Make sure you have latest golang setup with support for go mod.

1. Make sure you have minishift installed (https://docs.okd.io/latest/minishift/getting-started/index.html)

1. Clone repo
```
$ git clone mohammedzee1000/openshift-cluster-pool && cd openshift-cluster-pool
```

**Note:** You may need to setup requirements in go.mod appropriately before proceeding

### Build the binaries

Build all the nessasary binaries with 

```bash
$ make build
# or
# scripts/build.sh
```

### Setup DB

Make a test DB directory
```
mkdir `pwd`/test-badger
```
Load sample pool config

```
BADGER_DIR="`pwd`/test-badger" ./db-cli save-pool `pwd`/pool-examples/minishift-simple/minishift-simple.json
```

**Warning:** `db-cli` is a quick and dirty cli built to overcome missing cli for badger db.
It has limited function and should not be considered as replacement for the admin cli
which is in the pipeline

### Copy minishift provision scripts

Copy scripts to /usr/bin/
```
$ cp -avrf pool-examples/minishift-simple/usr/bin/* /usr/bin/
```

### Start pool manager

```bash
BADGER_DIR="`pwd`/test-badger" ./pool-manager
```

### Other important options of db-cli

As already stated `db-cli` is only for testing purposes. But you can use still use it.
Just do `./db-cli help` to find other commands

### Client API Server

A basic client api server is provided that can be launched with

```bash
$ BADGER_DIR="`pwd`/test-badger" HOST_ON=":20000" ./api-server
```

#### Possible operations:

##### List pool
 
 ```bash
$ curl http://localhost:20000/v1beta/pools/list
```

##### Basic describe a pool

```bash
$ curl http://localhost:20000/v1beta/pool/minishift-simple/short-describe
```

##### Activate cluster

```bash
$ curl http://localhost:20000/v1beta/pool/minishift-simple/get-cluster
```

##### Get cluster information

```bash
$ curl -k http://localhost:20000/v1beta/cluster/{clusterid}/describe
```

##### Return used cluster early for cleanup

```bash
$ curl -k http://localhost:20000/v1beta/cluster/{clusterid}/return
```

## In the pipeline
 - Admin API Server + client to replace db-cli
 - Client API server + client to allow users to list pools and get a cluster from pool / list its information

#### SAMPLE RUN: 

```bash
$ BADGER_DIR="`pwd`/test-badger" ./db-cli save-pool `pwd`/pool-examples/minishift-simple/minishift-simple.json
INFO[2019-11-22T12:19:09+05:30] Loading pool into DB                          component-name=save-pool name=test-db-cli

$ cp -avrf pool-examples/minishift-simple/usr/bin/* /usr/bin/

$ BADGER_DIR="`pwd`/test-badger" ./pool-manager
INFO[2019-11-22T12:45:25+05:30] starting...                                   component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] initiating cluster management cycle           component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] managing pools to be removed                  component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] no pools to remove, skipping this turn        component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] managing available pools                      component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] initiating GC of pool minishift-simple        component-name="Pool GC" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] initiating cleanup of clusters that have met some conditions, pool minishift-simple  component-name="Pool GC" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] no garbage in pool minishift-simple, skipping  component-name="Pool GC - By Condition" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] initiating cleanup of clusters that need to be removed due to expected size reduction change, pool minishift-simple  component-name="Pool GC" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] no garbage in pool minishift-simple, skipping  component-name="Pool GC - By Config Change" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] initiating reconciliation for pool minishift-simple  component-name="Pool Reconcile" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] available clusters do not match expected for pool minishift-simple  component-name="Pool Reconcile" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] allocating 2 clusters for pool minishift-simple serially  component-name="Pool Reconcile" name=pool-manager
INFO[2019-11-22T12:45:25+05:30] provisioning clusters of pool minishift-simple  component-name="Pool provision" name=pool-manager
INFO[2019-11-22T12:55:13+05:30] successfully provisioned clusters 0eea7e12-8b0f-4879-a2fc-6b57156f2b11, pool minishift-simple  component-name="Pool provision" name=pool-manager
INFO[2019-11-22T12:55:13+05:30] provisioning clusters of pool minishift-simple  component-name="Pool provision" name=pool-manager
INFO[2019-11-22T13:08:12+05:30] successfully provisioned clusters 0507ab26-864b-42e5-ad01-ea8a0c699b34, pool minishift-simple  component-name="Pool provision" name=pool-manager
INFO[2019-11-22T13:11:12+05:30] initiating cluster management cycle           component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:11:12+05:30] managing pools to be removed                  component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:11:12+05:30] no pools to remove, skipping this turn        component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:11:12+05:30] managing available pools                      component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:11:13+05:30] initiating GC of pool minishift-simple        component-name="Pool GC" name=pool-manager
INFO[2019-11-22T13:11:13+05:30] initiating cleanup of clusters that have met some conditions, pool minishift-simple  component-name="Pool GC" name=pool-manager
INFO[2019-11-22T13:11:13+05:30] no garbage in pool minishift-simple, skipping  component-name="Pool GC - By Condition" name=pool-manager
INFO[2019-11-22T13:11:13+05:30] initiating cleanup of clusters that need to be removed due to expected size reduction change, pool minishift-simple  component-name="Pool GC" name=pool-manager
INFO[2019-11-22T13:11:13+05:30] no garbage in pool minishift-simple, skipping  component-name="Pool GC - By Config Change" name=pool-manager
INFO[2019-11-22T13:11:13+05:30] initiating reconciliation for pool minishift-simple  component-name="Pool Reconcile" name=pool-manager
INFO[2019-11-22T13:11:13+05:30] skipping reconcilation for pool minishift-simple as actual matches expected  component-name="Pool Reconcile" name=pool-manager
...


$ minishift profile list
- 0507ab26-864b-42e5-ad01-ea8a0c699b34	Running		(Active)
- 0eea7e12-8b0f-4879-a2fc-6b57156f2b11	Running
- minishift				Does Not Exist

$ BADGER_DIR="`pwd`/test-badger" HOST_ON=":20000" ./api-server
/v1beta/pools/list
/v1beta/pool/{poolname}/get-cluster
/v1beta/pool/{poolname}/short-describe
/v1beta/cluster/{clusterid}/describe
/v1beta/cluster/{clusterid}/return
/v1beta/admin/pools/save
/v1beta/admin/pools/{poolname}/delete

$ curl http://localhost:20000/v1beta/pools/list | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    64  100    64    0     0   4266      0 --:--:-- --:--:-- --:--:--  4266
{
  "data": [
    "minishift-simple"
  ],
  "api_version": "v1beta",
  "error": ""
}

$ curl http://localhost:20000/v1beta/pool/minishift-simple/short-describe | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    87  100    87    0     0     82      0  0:00:01  0:00:01 --:--:--    82
{
  "description": "minishift-simple",
  "current_count": 2,
  "api_version": "v1beta",
  "error": ""
}

$ curl http://localhost:20000/v1beta/pool/minishift-simple/get-cluster | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  4364    0  4364    0     0  47956      0 --:--:-- --:--:-- --:--:-- 47956
{
  "Cluster": {
    "ClusterID": "0507ab26-864b-42e5-ad01-ea8a0c699b34",
    "PoolName": "minishift-simple",
    "State": "Used",
    "URL": "https://192.168.42.143:8443",
    "AdminUser": "developer",
    "AdminPassword": "developer",
    "CAFile": [
      "-----BEGIN CERTIFICATE-----",
      "MIIC1DCCAbygAwIBAgIRAMGo6ym9hf88AMi6BRFj1KkwDQYJKoZIhvcNAQELBQAw",
      "EzERMA8GA1UEChMIbXplZTEwMDAwHhcNMTkxMTIyMDcyMDAwWhcNMjIxMTA2MDcy",
      "MDAwWjATMREwDwYDVQQKEwhtemVlMTAwMDCCASIwDQYJKoZIhvcNAQEBBQADggEP",
      "ADCCAQoCggEBAMXpz6U/0mPqbmXEj5ZNnAR78FiEuVBDeeXJO0u4mV4qishItkUd",
      "6Nh+CJGM3wIwQOzns/4UBU/ghkbQBQ4iUyfgcQMmcnKL369NbQ18V0rDhYqJQyrJ",
      "kqusoidhsXH29px1sFTELVl6CfsLCBC6gWE80cv0We8RGZ2nSjXHk1LCXcy17FBm",
      "EPXERFV6rDtDWGnCYn3OW/TWmJENHzQ1PzQM0Xmdjf1SBKNiSrfl7tUQF98qsyED",
      "lGwCZ+EGngMhteuR/VLCzP4PbnnI/a2Vlwj7Qo3EWuaP52HhVXJ9xuNbdbFERQLg",
      "uJtEYaShnNq8NQlQ3D/B+YLnx2nSQ2l2CPsCAwEAAaMjMCEwDgYDVR0PAQH/BAQD",
      "AgKsMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAMCuSF1zd2Ps",
      "vZzCbQjrDNNYAcMaSmNfh5mlp3Fv/lXY2HyhgD4EE1/8l8yI7Cm0UMcOGQlzE5g7",
      "wBRGqZ2h1cfiD9hby/qXE1BDu7LdtSY6rKRvEwsyD5c2sMddRyXjGwc7YIPPPgFs",
      "CS7wOM3hPLcYOQ4vhc/ICGeqUeBSvYlb7Y5tEODZmMtTHWkmSJFGsAB/C2RUPcpA",
      "T0uJPBVq2JDnQ+k1u/ao9YaMz7niYehzbGSwH9BVWhgXsiqGsffN/E5p5Tuyi1rm",
      "2IYVSnIdL+4KVNfe9YCdfiiH+8jIJuYmmbzXESMneYS5Ga1HNG+ZovP+ioALtOVp",
      "0DczO8fYQYA=",
      "-----END CERTIFICATE-----"
    ],
    "CertFile": [
      "-----BEGIN CERTIFICATE-----",
      "MIIC8TCCAdmgAwIBAgIQYGcXI8xegRiOfih8wlfdIjANBgkqhkiG9w0BAQsFADAT",
      "MREwDwYDVQQKEwhtemVlMTAwMDAeFw0xOTExMjIwNzIwMDBaFw0yMjExMDYwNzIw",
      "MDBaMB8xHTAbBgNVBAoMFG16ZWUxMDAwLjxib290c3RyYXA+MIIBIjANBgkqhkiG",
      "9w0BAQEFAAOCAQ8AMIIBCgKCAQEArkhwb3swJi4zfO/usM0PcIF8iaaVAsZ8hB+S",
      "a/cNBS44v8GGbMFE0yqcRqQsDiXWhTXzAPw9Tl/aBJ4tpHDDkWl4hzwH8ciBiG5y",
      "dNhBHmT6w8JLYcyEZbzkot1lMo3nz7zg64sQSvRO2LNsuxWIC0NoOtFOwDHrekOP",
      "94pLjBwVbpoHHbWxlxgDhGyvYlr0RQOmyegGpt+NONcUnH90KOn34z1sCEYt6jX1",
      "8B3+MLcqEGk7DaVEg6wXPf7a5NeW3XyPx329JzUp5W/8rCaOp/quixJBEqh0l322",
      "8h3YytMzXHIoA18yGLa0QHjDGHLZIEwGif4TzPQufz70hsR2twIDAQABozUwMzAO",
      "BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw",
      "ADANBgkqhkiG9w0BAQsFAAOCAQEAqzv4nZifjojFomrg7vDVN7R87IO3j/4D3Pi+",
      "iETPpwx3tve/q5IQlSov1edV6ZAX5c1duit5hz4b2c2JK5zCldMNOHK2b96U/Gdh",
      "90QKRpekT+nK6BZSco694AM3iLf4YfontaIaUi5hQA5wI6yMQY/omvcj2HlLJHBs",
      "aidpeUpdRRcx7UrWVbUoR7kNIL/ZSBVxXBrEi6yRIIuMKsgGWHOPSqfmltTPpAfp",
      "Q3JDCG7y73t+6jCVRzL+aqresccd8ZQnpukVM/TbAh0blMv4k2YeBFWI7AynBWbH",
      "SvUUciZTSJyauMbm6golrdUcpoAyb7KAGlnNk7RX7yb9rb8a/Q==",
      "-----END CERTIFICATE-----"
    ],
    "KeyFile": [
      "-----BEGIN RSA PRIVATE KEY-----",
      "MIIEpQIBAAKCAQEArkhwb3swJi4zfO/usM0PcIF8iaaVAsZ8hB+Sa/cNBS44v8GG",
      "bMFE0yqcRqQsDiXWhTXzAPw9Tl/aBJ4tpHDDkWl4hzwH8ciBiG5ydNhBHmT6w8JL",
      "YcyEZbzkot1lMo3nz7zg64sQSvRO2LNsuxWIC0NoOtFOwDHrekOP94pLjBwVbpoH",
      "HbWxlxgDhGyvYlr0RQOmyegGpt+NONcUnH90KOn34z1sCEYt6jX18B3+MLcqEGk7",
      "DaVEg6wXPf7a5NeW3XyPx329JzUp5W/8rCaOp/quixJBEqh0l3228h3YytMzXHIo",
      "A18yGLa0QHjDGHLZIEwGif4TzPQufz70hsR2twIDAQABAoIBAAm7qevP6WR4eA+m",
      "JqJhEVerI5VcZD3/b7zBNqAo7+U2K50p5aP1Ny7D1m5rhLpViqFt3eBUNehGmhpf",
      "6xSf54wbY8vJonfyRqmj4Wh9G0XjRc3g7+zKSyqTXgFqc9ha7HNBjR4aahKFilG3",
      "036vOSXH4e5G+irpnsj5NPUSGB3+6+bAS+9V1C7fzifKvE7d69T1KA9veXpLeD+p",
      "ly6gGpTj3S+PYJj0IAenTjUeN2agKvJKnHp9M3ETWZlOUy8ILNjSmCGbEzfY7y0q",
      "3PeVNJq56NVUePOQcVEGm/ub9r7BH6xjIflGy5a/7rlcuRDyEXa6vxbPmrncySfv",
      "NIa/FoECgYEAys8qMAZoLDd7FN9M68JK17WLoaoHXxdWdwxQYKuC6viW2YhnIKNM",
      "/9HXBeii9dDydEWOhDpU5QYY3/LB1SUJRpMxZP46iXc7xFJU+phTswVe+PDLhFc7",
      "m67Dh/umEcsYdT6l8M2jdRrJILyH2+8wggBICA7jQG+g96I7lCPcMYcCgYEA2/39",
      "OMKvemmvmwS/O4oVhLewL1SCTego8ZVQyTCL+rXjntCV7wsc7FyvGs5gIyyZO+pc",
      "VsiQlnK/g6DbWDCaSGL0FqC/7UICawXz+B0vLDeOT4keSgTpAjbKnrhcMa9w/YVm",
      "fWgKS9omRM5mzQZY1YfMv/sOaGKs0CIYDFC/nVECgYEAmH3OTc/zchPBUv9Xugkb",
      "9zeFJuhOpJxKoja7FQTA6mZCHoxmZm8DDXM9Ry8VoNkcBHrsXtXxUUcVWwYP4nD3",
      "mX5BbJuPbh8d7E6voMD6ZigKsgl0LSzeH//2+38m7kgUOswBP5+PYRTj196KFL+z",
      "bHxDrPNswd1tXeU5APk5rm8CgYEAjt+2zF1MYAExhkkf9Ygpj4dIyoRlGDnWFYf2",
      "7qMz1gC5MtSe+5/JCgzrwEoWD+IQJuR/UfFyTfN6Q/99VRpDqQ1zHxsJawp6zY0R",
      "NKunjl0KdMdFv6bOuZxiHZD4d2BMzqoLtRiTz01/myI9i5w6p3tJ08k2Qz8KoyXx",
      "XlY3C/ECgYEArIOkBBB1YYTfwSwJW/SLon0FU44W/U4jmS7ufGD1ASV8DB8jEH0M",
      "7YAwpR764OwFDrRKXvhFEftnUmi5LYOeVWqxb+N3m+vJaA8fsa+sIOqc68FsJB/e",
      "GvAqG7i7dMJTld8smYAUCbgG9LR1gXkR54OMUDFQXEwg4xwsYDvMhfo=",
      "-----END RSA PRIVATE KEY-----"
    ],
    "ExtraInfo": "",
    "CreatedOn": "2019-11-22T13:08:11.193430502+05:30",
    "ActivatedOn": "2019-11-22T13:12:20.034214837+05:30"
  },
  "ExpiresOn": "2024-06-07T14:32:16.322223029+05:30", //non output note, this display logic needs to be fixed, but it should work as expected
  "api_version": "v1beta",
  "error": ""
}

$ curl -k http://localhost:20000/v1beta/cluster/0507ab26-864b-42e5-ad01-ea8a0c699b34/describe | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  4364    0  4364    0     0   133k      0 --:--:-- --:--:-- --:--:--  137k
{
  "Cluster": {
    "ClusterID": "0507ab26-864b-42e5-ad01-ea8a0c699b34",
    "PoolName": "minishift-simple",
    "State": "Used",
    "URL": "https://192.168.42.143:8443",
    "AdminUser": "developer",
    "AdminPassword": "developer",
    "CAFile": [
      "-----BEGIN CERTIFICATE-----",
      "MIIC1DCCAbygAwIBAgIRAMGo6ym9hf88AMi6BRFj1KkwDQYJKoZIhvcNAQELBQAw",
      "EzERMA8GA1UEChMIbXplZTEwMDAwHhcNMTkxMTIyMDcyMDAwWhcNMjIxMTA2MDcy",
      "MDAwWjATMREwDwYDVQQKEwhtemVlMTAwMDCCASIwDQYJKoZIhvcNAQEBBQADggEP",
      "ADCCAQoCggEBAMXpz6U/0mPqbmXEj5ZNnAR78FiEuVBDeeXJO0u4mV4qishItkUd",
      "6Nh+CJGM3wIwQOzns/4UBU/ghkbQBQ4iUyfgcQMmcnKL369NbQ18V0rDhYqJQyrJ",
      "kqusoidhsXH29px1sFTELVl6CfsLCBC6gWE80cv0We8RGZ2nSjXHk1LCXcy17FBm",
      "EPXERFV6rDtDWGnCYn3OW/TWmJENHzQ1PzQM0Xmdjf1SBKNiSrfl7tUQF98qsyED",
      "lGwCZ+EGngMhteuR/VLCzP4PbnnI/a2Vlwj7Qo3EWuaP52HhVXJ9xuNbdbFERQLg",
      "uJtEYaShnNq8NQlQ3D/B+YLnx2nSQ2l2CPsCAwEAAaMjMCEwDgYDVR0PAQH/BAQD",
      "AgKsMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAMCuSF1zd2Ps",
      "vZzCbQjrDNNYAcMaSmNfh5mlp3Fv/lXY2HyhgD4EE1/8l8yI7Cm0UMcOGQlzE5g7",
      "wBRGqZ2h1cfiD9hby/qXE1BDu7LdtSY6rKRvEwsyD5c2sMddRyXjGwc7YIPPPgFs",
      "CS7wOM3hPLcYOQ4vhc/ICGeqUeBSvYlb7Y5tEODZmMtTHWkmSJFGsAB/C2RUPcpA",
      "T0uJPBVq2JDnQ+k1u/ao9YaMz7niYehzbGSwH9BVWhgXsiqGsffN/E5p5Tuyi1rm",
      "2IYVSnIdL+4KVNfe9YCdfiiH+8jIJuYmmbzXESMneYS5Ga1HNG+ZovP+ioALtOVp",
      "0DczO8fYQYA=",
      "-----END CERTIFICATE-----"
    ],
    "CertFile": [
      "-----BEGIN CERTIFICATE-----",
      "MIIC8TCCAdmgAwIBAgIQYGcXI8xegRiOfih8wlfdIjANBgkqhkiG9w0BAQsFADAT",
      "MREwDwYDVQQKEwhtemVlMTAwMDAeFw0xOTExMjIwNzIwMDBaFw0yMjExMDYwNzIw",
      "MDBaMB8xHTAbBgNVBAoMFG16ZWUxMDAwLjxib290c3RyYXA+MIIBIjANBgkqhkiG",
      "9w0BAQEFAAOCAQ8AMIIBCgKCAQEArkhwb3swJi4zfO/usM0PcIF8iaaVAsZ8hB+S",
      "a/cNBS44v8GGbMFE0yqcRqQsDiXWhTXzAPw9Tl/aBJ4tpHDDkWl4hzwH8ciBiG5y",
      "dNhBHmT6w8JLYcyEZbzkot1lMo3nz7zg64sQSvRO2LNsuxWIC0NoOtFOwDHrekOP",
      "94pLjBwVbpoHHbWxlxgDhGyvYlr0RQOmyegGpt+NONcUnH90KOn34z1sCEYt6jX1",
      "8B3+MLcqEGk7DaVEg6wXPf7a5NeW3XyPx329JzUp5W/8rCaOp/quixJBEqh0l322",
      "8h3YytMzXHIoA18yGLa0QHjDGHLZIEwGif4TzPQufz70hsR2twIDAQABozUwMzAO",
      "BgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw",
      "ADANBgkqhkiG9w0BAQsFAAOCAQEAqzv4nZifjojFomrg7vDVN7R87IO3j/4D3Pi+",
      "iETPpwx3tve/q5IQlSov1edV6ZAX5c1duit5hz4b2c2JK5zCldMNOHK2b96U/Gdh",
      "90QKRpekT+nK6BZSco694AM3iLf4YfontaIaUi5hQA5wI6yMQY/omvcj2HlLJHBs",
      "aidpeUpdRRcx7UrWVbUoR7kNIL/ZSBVxXBrEi6yRIIuMKsgGWHOPSqfmltTPpAfp",
      "Q3JDCG7y73t+6jCVRzL+aqresccd8ZQnpukVM/TbAh0blMv4k2YeBFWI7AynBWbH",
      "SvUUciZTSJyauMbm6golrdUcpoAyb7KAGlnNk7RX7yb9rb8a/Q==",
      "-----END CERTIFICATE-----"
    ],
    "KeyFile": [
      "-----BEGIN RSA PRIVATE KEY-----",
      "MIIEpQIBAAKCAQEArkhwb3swJi4zfO/usM0PcIF8iaaVAsZ8hB+Sa/cNBS44v8GG",
      "bMFE0yqcRqQsDiXWhTXzAPw9Tl/aBJ4tpHDDkWl4hzwH8ciBiG5ydNhBHmT6w8JL",
      "YcyEZbzkot1lMo3nz7zg64sQSvRO2LNsuxWIC0NoOtFOwDHrekOP94pLjBwVbpoH",
      "HbWxlxgDhGyvYlr0RQOmyegGpt+NONcUnH90KOn34z1sCEYt6jX18B3+MLcqEGk7",
      "DaVEg6wXPf7a5NeW3XyPx329JzUp5W/8rCaOp/quixJBEqh0l3228h3YytMzXHIo",
      "A18yGLa0QHjDGHLZIEwGif4TzPQufz70hsR2twIDAQABAoIBAAm7qevP6WR4eA+m",
      "JqJhEVerI5VcZD3/b7zBNqAo7+U2K50p5aP1Ny7D1m5rhLpViqFt3eBUNehGmhpf",
      "6xSf54wbY8vJonfyRqmj4Wh9G0XjRc3g7+zKSyqTXgFqc9ha7HNBjR4aahKFilG3",
      "036vOSXH4e5G+irpnsj5NPUSGB3+6+bAS+9V1C7fzifKvE7d69T1KA9veXpLeD+p",
      "ly6gGpTj3S+PYJj0IAenTjUeN2agKvJKnHp9M3ETWZlOUy8ILNjSmCGbEzfY7y0q",
      "3PeVNJq56NVUePOQcVEGm/ub9r7BH6xjIflGy5a/7rlcuRDyEXa6vxbPmrncySfv",
      "NIa/FoECgYEAys8qMAZoLDd7FN9M68JK17WLoaoHXxdWdwxQYKuC6viW2YhnIKNM",
      "/9HXBeii9dDydEWOhDpU5QYY3/LB1SUJRpMxZP46iXc7xFJU+phTswVe+PDLhFc7",
      "m67Dh/umEcsYdT6l8M2jdRrJILyH2+8wggBICA7jQG+g96I7lCPcMYcCgYEA2/39",
      "OMKvemmvmwS/O4oVhLewL1SCTego8ZVQyTCL+rXjntCV7wsc7FyvGs5gIyyZO+pc",
      "VsiQlnK/g6DbWDCaSGL0FqC/7UICawXz+B0vLDeOT4keSgTpAjbKnrhcMa9w/YVm",
      "fWgKS9omRM5mzQZY1YfMv/sOaGKs0CIYDFC/nVECgYEAmH3OTc/zchPBUv9Xugkb",
      "9zeFJuhOpJxKoja7FQTA6mZCHoxmZm8DDXM9Ry8VoNkcBHrsXtXxUUcVWwYP4nD3",
      "mX5BbJuPbh8d7E6voMD6ZigKsgl0LSzeH//2+38m7kgUOswBP5+PYRTj196KFL+z",
      "bHxDrPNswd1tXeU5APk5rm8CgYEAjt+2zF1MYAExhkkf9Ygpj4dIyoRlGDnWFYf2",
      "7qMz1gC5MtSe+5/JCgzrwEoWD+IQJuR/UfFyTfN6Q/99VRpDqQ1zHxsJawp6zY0R",
      "NKunjl0KdMdFv6bOuZxiHZD4d2BMzqoLtRiTz01/myI9i5w6p3tJ08k2Qz8KoyXx",
      "XlY3C/ECgYEArIOkBBB1YYTfwSwJW/SLon0FU44W/U4jmS7ufGD1ASV8DB8jEH0M",
      "7YAwpR764OwFDrRKXvhFEftnUmi5LYOeVWqxb+N3m+vJaA8fsa+sIOqc68FsJB/e",
      "GvAqG7i7dMJTld8smYAUCbgG9LR1gXkR54OMUDFQXEwg4xwsYDvMhfo=",
      "-----END RSA PRIVATE KEY-----"
    ],
    "ExtraInfo": "",
    "CreatedOn": "2019-11-22T13:08:11.193430502+05:30",
    "ActivatedOn": "2019-11-22T13:12:20.034214837+05:30"
  },
  "ExpiresOn": "2024-06-07T14:32:16.322223029+05:30",
  "api_version": "v1beta",
  "error": ""
}

$ curl -k http://localhost:20000/v1beta/cluster/0507ab26-864b-42e5-ad01-ea8a0c699b34/return | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    48  100    48    0     0    923      0 --:--:-- --:--:-- --:--:--   923
{
  "data": "OK",
  "api_version": "v1beta",
  "error": ""
}

// rerun just to skip wait cycle
$ BADGER_DIR="`pwd`/test-badger" ./pool-manager
INFO[2019-11-22T13:18:52+05:30] starting...                                   component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:18:52+05:30] initiating cluster management cycle           component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:18:52+05:30] managing pools to be removed                  component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:18:52+05:30] no pools to remove, skipping this turn        component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:18:52+05:30] managing available pools                      component-name="Pool Manager" name=pool-manager
INFO[2019-11-22T13:18:52+05:30] initiating GC of pool minishift-simple        component-name="Pool GC" name=pool-manager
INFO[2019-11-22T13:18:52+05:30] initiating cleanup of clusters that have met some conditions, pool minishift-simple  component-name="Pool GC" name=pool-manager
INFO[2019-11-22T13:18:52+05:30] deprovisioning clusters of pool minishift-simple  component-name="Pool deprovision" name=pool-manager
INFO[2019-11-22T13:18:53+05:30] successfully deprovisioned clusters 0507ab26-864b-42e5-ad01-ea8a0c699b34, pool minishift-simple  component-name="Pool deprovision" name=pool-manager
INFO[2019-11-22T13:18:53+05:30] initiating cleanup of clusters that need to be removed due to expected size reduction change, pool minishift-simple  component-name="Pool GC" name=pool-manager
INFO[2019-11-22T13:18:53+05:30] no garbage in pool minishift-simple, skipping  component-name="Pool GC - By Config Change" name=pool-manager
INFO[2019-11-22T13:18:53+05:30] initiating reconciliation for pool minishift-simple  component-name="Pool Reconcile" name=pool-manager
INFO[2019-11-22T13:18:53+05:30] available clusters do not match expected for pool minishift-simple  component-name="Pool Reconcile" name=pool-manager
INFO[2019-11-22T13:18:53+05:30] allocating 1 clusters for pool minishift-simple serially  component-name="Pool Reconcile" name=pool-manager
INFO[2019-11-22T13:18:53+05:30] provisioning clusters of pool minishift-simple  component-name="Pool provision" name=pool-manager

$ minishift profile list
- 0eea7e12-8b0f-4879-a2fc-6b57156f2b11	Running
- 648685b7-58e9-43cf-8a0d-130caa59f0bc	Running		(Active)
- minishift				Does Not Exist
```
