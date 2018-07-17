# go-teamcity-sdk #

[![GoDoc](https://godoc.org/github.com/cvbarros/go-teamcity-sdk/github?status.svg)](https://godoc.org/github.com/cvbarros/go-teamcity-sdk) [![Build Status](https://travis-ci.org/cvbarros/go-teamcity-sdk.svg?branch=master)](https://travis-ci.org/cvbarros/go-teamcity-sdk)

go-teamcity-sdk is a Go client library interacting with Jetbrain's [TeamCity CI server REST API](https://confluence.jetbrains.com/display/TCD10/REST+API).  
Currently, supported API version is TeamCity 10.x (2017.2.1) or up.

## Motivation ##
go-teamcity-sdk is intended for automating Jetbrain's TeamCity tasks and configurations using Go.  
It is also the client library used in [TeamCity Provider](https://github.com/cvbarros/terraform-provider-teamcity), a [Hashicorp's Terraform](https://www.terraform.io) custom provider to configure TeamCity fully in code.

## Usage ##

```go
import teamcity "github.com/cvbarros/go-teamcity-sdk"
```

Construct a new Teamcity client with a default `http.Client` for the server at the `TEAMCITY_ADDR` environment variable, then use the various services exposed by the client to use TeamCity API's features.


```go
client := teamcity.New("admin", "admin", http.DefaultClient)

// retrieve a project identified 'myproject', using the id: locator
project,  err := client.Projects.GetById("myproject")
```

> You may also use an `teamcity.NewWithAddress(...)` to explicitly provide the server address.

### Examples
For now, [integration tests](https://github.com/cvbarros/go-teamcity-sdk/search?q=filename%3A*_test.go&unscoped_q=filename%3A*_test.go) are the best examples on how to use the library to interact with the several services.

## Project Status ##

This project is under heavy development for now. We follow [semantic versioning](https://semver.org) conventions, meaning that the public API compability under `1.0` is suitable to change without prior notice.

Issues, contributions and feature requests are welcome, please see [contributing](CONTRIBUTING.MD) for guidance on how to get involved.

The development and design of this library is based on the [go-github](https://github.com/google/go-github) project.

## License ##

This library is distributed under the MPL-2.0 license fully described in the [LICENSE](./LICENSE) file.