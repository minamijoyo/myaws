## master (Unreleased)

## 0.4.3 (2022/04/12)

NEW FEATURES:

* [ssm parameter env] add quote parameter ([#59](https://github.com/minamijoyo/myaws/pull/59))

ENHANCEMENTS:

* Use golangci-lint instead of golint ([#57](https://github.com/minamijoyo/myaws/pull/57))
* Update golangci-lint-action to v3 ([#58](https://github.com/minamijoyo/myaws/pull/60))
* [github actions] fix checkout step ([#61](https://github.com/minamijoyo/myaws/pull/61))

## 0.4.2 (2021/11/18)

ENHANCEMENTS:

* Add release target to arm64 architecture ([#54](https://github.com/minamijoyo/myaws/pull/54))
* Update Go to v1.17.3 and Alpine to 3.14 ([#55](https://github.com/minamijoyo/myaws/pull/55))
* Remove docker build on pull request ([#56](https://github.com/minamijoyo/myaws/pull/56))

## 0.4.1 (2021/10/28)

* Restrict permissions for GitHub Actions ([#51](https://github.com/minamijoyo/myaws/pull/51))
* Set timeout for GitHub Actions ([#52](https://github.com/minamijoyo/myaws/pull/52))

## 0.4.0 (2021/07/19)

This releases contains small breaking changes to improve CI/CD workflow. AWS related functionalities didn't change.

BREAKING CHANGES:

* Build & push docker images on GitHub Actions ([#50](https://github.com/minamijoyo/myaws/pull/50))

The `latest` tag of docker image now points at the latest release. Previously the `latest` tag pointed at the master branch, if you want to use the master branch, use the `master` tag instead.

* Set a version number explicitly in source ([#43](https://github.com/minamijoyo/myaws/pull/43))

The `version` command now contains only a version number, not a revision (commit SHA1).

ENHANCEMENTS:

* Fix release archives in goreleaser.yml ([#49](https://github.com/minamijoyo/myaws/pull/49))
* Fix go mod tidy ([#48](https://github.com/minamijoyo/myaws/pull/48))
* Drop goreleaser dependencies ([#47](https://github.com/minamijoyo/myaws/pull/47))
* Move CI to GitHub Actions ([#46](https://github.com/minamijoyo/myaws/pull/46))
* Ignore updating README and CHANGELOG in release notes ([#45](https://github.com/minamijoyo/myaws/pull/45))
* Cache go modules in docker build ([#44](https://github.com/minamijoyo/myaws/pull/44))

