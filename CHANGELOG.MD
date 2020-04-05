# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0]  

This release is the first major release and includes an upgrade to TeamCity 2019.2.2 as the supported version.
Whenever a major version is bumped, backwards compatibility with previous TeamCity versions aren't guaranteed.

### Changes:
- Updating to use TeamCity 2019.2.2 [#72]
- Added support for TeamCity User Group API [#68]
- Support token auth [#67]

## [0.5.1]

### Added

- Add support for `octopus.push.package` step. It allows creating and pushing packages from [TeamCity to Octopus](https://octopus.com/docs/api-and-integration/teamcity). [#53]

### Fixed
- Issue with serialization and newlines for VCSRoot branchSpec property

### Changed

- (**BREAKING CHANGE**) Move package code from `pkg/teamcity` to `teamcity`.
  - The former structure caused some issues with tooling such as `gocode` and `godef`.
  Since this project follows [go-github](https://github.com/google/go-github) design and conventions, I've adapted the project structure as well.

## [0.4.0] - 2018-10-11

### Added

- Add `GetName()` to Step interface.

## [0.2.0] - 2018-07-16

## [0.1.0] - 2018-07-16

First release.

[//]: # (Release links)
[1.0.0]: https://github.com/cvbarros/go-teamcity-sdk/releases/tag/v1.0.0
[0.5.1]: https://github.com/cvbarros/go-teamcity-sdk/releases/tag/v0.5.1
[0.4.0]: https://github.com/cvbarros/go-teamcity-sdk/releases/tag/v0.4.0
[0.2.0]: https://github.com/cvbarros/go-teamcity-sdk/releases/tag/v0.2.0
[0.1.0]: https://github.com/cvbarros/go-teamcity-sdk/releases/tag/v0.1.0

[//]: # (Issue/PR links)
[#53]: https://github.com/cvbarros/go-teamcity-sdk/pull/53
[#67]: https://github.com/cvbarros/go-teamcity/pull/67
[#68]: https://github.com/cvbarros/go-teamcity/pull/68
[#72]: https://github.com/cvbarros/go-teamcity/pull/72