# 2018-12-20
BREAKING CHANGE: Moved package code from `pkg/teamcity` to `teamcity`. The former structure caused some issues with tooling such as `gocode` and `godef`.
Since this project follows [go-github](https://github.com/google/go-github) design and conventions, I've adapted the project structure as well.