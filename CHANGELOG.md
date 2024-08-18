## master (Unreleased)

## 0.4.9 (2024/08/18)

ENHANCEMENTS:

* Update actions/checkout to v4 ([#70](https://github.com/minamijoyo/myaws/pull/70))
* Update setup-go to v5 ([#71](https://github.com/minamijoyo/myaws/pull/71))
* Update golangci lint to v1.59.1 ([#72](https://github.com/minamijoyo/myaws/pull/72))
* Update Go to v1.22 ([#73](https://github.com/minamijoyo/myaws/pull/73))
* Update docker/build-push-action to v4 ([#74](https://github.com/minamijoyo/myaws/pull/74))
* Update goreleaser to v2 ([#75](https://github.com/minamijoyo/myaws/pull/75))
* Switch to the official action for creating GitHub App token ([#76](https://github.com/minamijoyo/myaws/pull/76))

## 0.4.8 (2022/08/24)

ENHANCEMENTS:

* Update Go to 1.19 ([#67](https://github.com/minamijoyo/myaws/pull/67))

## 0.4.7 (2022/07/28)

ENHANCEMENTS:

* Add support for linux/arm64 Docker image ([#66](https://github.com/minamijoyo/myaws/pull/66))

## 0.4.6 (2022/07/21)

ENHANCEMENTS:

* Restrict repository of token for release ([#65](https://github.com/minamijoyo/myaws/pull/65))

## 0.4.5 (2022/07/14)

ENHANCEMENTS:

* Use GitHub App token for updating brew formula on release ([#64](https://github.com/minamijoyo/myaws/pull/64))

## 0.4.4 (2022/05/10)

ENHANCEMENTS:

* bump golang.org/x/crypto to v0.0.0-20220507011949-2cf3adece122 ([#63](https://github.com/minamijoyo/myaws/pull/63))

## 0.4.3 (2022/04/12)

NEW FEATURES:

* [ssm parameter env] add quote parameter ([#59](https://github.com/minamijoyo/myaws/pull/59))

ENHANCEMENTS:

* Use golangci-lint instead of golint ([#57](https://github.com/minamijoyo/myaws/pull/57))
* Update golangci-lint-action to v3 ([#60](https://github.com/minamijoyo/myaws/pull/60))
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

