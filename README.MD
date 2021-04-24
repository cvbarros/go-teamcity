# go-teamcity #

[![GoDoc](https://godoc.org/github.com/cvbarros/go-teamcity/github?status.svg)](https://godoc.org/github.com/cvbarros/go-teamcity) [![Build Status](https://travis-ci.com/cvbarros/go-teamcity.svg?branch=master)](https://travis-ci.com/cvbarros/go-teamcity)


## LOOKING FOR MAINTAINTERS ##

This project is in need of new maintainers, please see here:  
https://github.com/cvbarros/go-teamcity/issues/94.  
https://github.com/cvbarros/terraform-provider-teamcity/issues/113

---
_go-teamcity_ is a Go client library interacting with Jetbrain's [TeamCity CI server REST API](https://confluence.jetbrains.com/display/TCD10/REST+API).  
Currently, guaranteed supported API version (on CI tests) is **TeamCity 2019.2.2 or up**.

## Motivation ##
go-teamcity is intended for automating Jetbrain's TeamCity tasks and configurations using Go.
It is also the client library used in [TeamCity Provider](https://github.com/cvbarros/terraform-provider-teamcity), a [Hashicorp's Terraform](https://www.terraform.io) custom provider to configure TeamCity as code.

## Usage ##

```go
import teamcity "github.com/cvbarros/go-teamcity"
```

Construct a new Teamcity client with a default `http.Client` for the server at the `TEAMCITY_ADDR` environment variable, then use the various services exposed by the client to use TeamCity API's features.

```go
client := teamcity.NewClient(teamcity.BasicAuth("admin", "admin"), http.DefaultClient)
```
or
```go
client := teamcity.NewClient(teamcity.TokenAuth("token"), http.DefaultClient)
``` 
```go
// retrieve a project identified 'myproject', using the id: locator
project,  err := client.Projects.GetById("myproject")
```

> You may also use an `teamcity.NewClientWithAddress(...)` to explicitly provide the server address.

### Examples ###

For now, [integration tests](https://github.com/cvbarros/go-teamcity/search?q=filename%3A*_test.go&unscoped_q=filename%3A*_test.go) are the best examples on how to use the library to interact with the several services.

## Other Information ##

We follow [semantic versioning](https://semver.org) conventions. Thus releases are tagged in the `<major>.<minor>.<patch>` format, meaning:

- Public API breaking changes or major TeamCity version support (2018, 2019, etc) **MUST** increment the `major` version.
- Features and considerable refactorings without breaking compatibility should increment the `minor`.
- Bugfixes and minor patches should increment `patch`

All changes released are documented on the [CHANGELOG](CHANGELOG.MD).  

Issues, contributions and feature requests are welcome, please see [contributing](CONTRIBUTING.MD) for guidance on how to get involved.

The development and design of this library is based on the [go-github](https://github.com/google/go-github) project.

## License ##

This library is distributed under the MPL-2.0 license fully described in the [LICENSE](./LICENSE) file.
