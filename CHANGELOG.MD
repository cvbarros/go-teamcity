# 2019-07-30
**BREAKING CHANGE**: This project now has been renamed from `go-teamcity-sdk` to `go-teamcity`. 
This will break any imports from `master` as of b267d09410b4aa86910d469f0ef11709ab543ae2.

# 2018-12-20
BREAKING CHANGE: Moved package code from `pkg/teamcity` to `teamcity`. The former structure caused some issues with tooling such as `gocode` and `godef`.
Since this project follows [go-github](https://github.com/google/go-github) design and conventions, I've adapted the project structure as well.