# params - command line parameter library


[![GoDoc](https://godoc.org/github.com/ExploratoryEngineering/params?status.svg)](https://godoc.org/github.com/ExploratoryEngineering/params)
[![Build Status](https://travis-ci.org/ExploratoryEngineering/params.svg?branch=master)](https://travis-ci.org/ExploratoryEngineering/params)
[![Go Report Card](https://goreportcard.com/badge/github.com/ExploratoryEngineering/params)](https://goreportcard.com/report/github.com/ExploratoryEngineering/params)
[![codecov](https://codecov.io/gh/ExploratoryEngineering/params/branch/master/graph/badge.svg)](https://codecov.io/gh/ExploratoryEngineering/params)


This is (yet another) command line parameter library. The main goal for this
library is to make it easy to set up and use the command line parameters. The
standard Go `flags` library is quite tedious to use and you find yourself
doing the same validation over and over again.

Command line parameters are declared as annotations on structs. Reflection is
used to build the set of command line parameters.

Declaring the command line parameters with the defaults are straightforward:

```golang
type HTTPConfig struct {
    Endpoint      string `param:"desc=Server endpoint;default=:8080"`
    TLSCertFile   string `param:"desc=TLS cert file;file"`
    TLSKeyFile    string `param:"desc=TLS key file;file"`
    ACMECert      bool   `param:"desc=Let's Encrypt ACME certs;default=false"`
    ACMEHosts     string `param:"desc=ACME host names"`
    ACMESecretDir string `param:"desc=ACME secret dir"`
}
```

The parameters is named according to the casing and only the public parameters are used. Members that aren't annotated or public won't be used. The above example produces the following parameters:

* `--endpoint`
* `--tls-cert-file`
* `--tls-key-file`
* `--acme-cert`
* `--acme-hosts`
* `--acme-secret-dir`

The following data types are supported. The type is inferred from the type of the struct member.

* strings (also as a set of options, see example below)
* integers
* booleans
* floats
* duration
* files (see the file directive above). The file must exist for a valid parameter.

## Nesting structures

Parameter structs can be nested. The parameters inside the struct will be prefixed according to the name of the containing struct. Note that the parameter struct itself doesn't have an annotation.

If you create a new struct that wraps the `HTTPConfig` struct above like this:
```golang
type parameters struct {
    HTTP        httpConfig
    LogType     string        `param:"desc=Log type;options=plain,syslog,fancy,ansi,full;default=plain"`
    MyOtherBool bool          `param:"default=true"`
}
```

You would get the following list of command line parameters:

* `--log-type`
* `--my-other-bool`
* `--http-endpoint`
* `--http-tls-cert-file`
* `--http-tls-key-file`
* `--http-acme-cert`
* `--http-acme-hosts`
* `--http-acme-secret-dir`

Parameters can be nested within parameters but don't go overboard. Remember: Someone has to type the parameters at one point and it might be you. It might be tempting to name the parameter struct `HTTPConfig` but that would result in some odd looking command line parameters so you should spend a few seconds contemplating the naming of member structs and parameters.

## Reading the parameters

A single call will read and check the parameters. The error message can be used directly on the console:

```golang
var config parameters
if err := params.NewFlag(&config, os.Args[1:]); err == nil {
    fmt.Println(err.Error())
    return
}
```

## Environment variables

Parameters can be specified via environment variables as well. The environment variables are ALL_CAPS and substitutes the dash for underscore. The parameter `htt-tls-cert-file` would be `HTTP_TLS_CERT_FILE`. The environment setting overrides any command line parameters that are used.

```shell
# This will set the LogType parameter to "plain"
[local ~]$ LOG_TYPE=plan ./my-command

# This will also set the LogType parameter to "plain" since the environment variable overrides
# the parameter
[local ~]$ LOG_TYPE=plan ./my-command --log-type=fancy
```

You can either read *just* the environment variables with `params.NewEnv` or read both environment and command line parameters at the same time with `params.NewEnvFlag` which is probably the one you are going to use the most:

```golang
var config parameters
if err := params.NewEnvFlag(&config, os.Args[1:]); err == nil {
    fmt.Println(err.Error())
    return
}
```

## Parameter files

The third option is to use a configuration file. Each property is camelCased and nested ccording to the same rules so if you want to read the `parameters` struct above you can use this configuration file:

```json
{
    "logType": "plain",
    "myOtherBool": false,
    "http": {
        "endpoint": "localhost:1234",
        "tlsCertFile": "",
        "tlsKeyFile": "",
        "acmeCert": true,
        "acmeHosts": "some.example.com",
        "acmeSecretDir": "/var/secret"
    }
}
```

Reading the file takes a reader:

```golang

var config parameters
f, err := os.Open("config.json")
if err != nil {
    panic(err.Error())
}
defer f.Close()
if err := params.NewFile(&config, f); err == nil {
    fmt.Println(err.Error())
    return
}
```

Note that the reader doesn't have to be a file. A reader is a reader so you could just as easily read the configuration from a network stream.
