# 2019-07-30
**BREAKING CHANGE**: This project now has been renamed from `go-teamcity-sdk` to `go-teamcity`. 
This will break any imports from `master` as of https://github.com/cvbarros/go-teamcity/commit/b637df4f43e5bc61501b97e8d88448317fb7c241.

# 2018-12-20
BREAKING CHANGE: Moved package code from `pkg/teamcity` to `teamcity`. The former structure caused some issues with tooling such as `gocode` and `godef`.
Since this project follows [go-github](https://github.com/google/go-github) design and conventions, I've adapted the project structure as well.