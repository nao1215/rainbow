<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
![Coverage](https://raw.githubusercontent.com/nao1215/octocovs-central-repo/main/badges/nao1215/rainbow/coverage.svg)
[![LinuxUnitTest](https://github.com/nao1215/rainbow/actions/workflows/linux_test.yml/badge.svg)](https://github.com/nao1215/rainbow/actions/workflows/linux_test.yml)
[![WindowsUnitTest](https://github.com/nao1215/rainbow/actions/workflows/windows_test.yml/badge.svg)](https://github.com/nao1215/rainbow/actions/workflows/windows_test.yml)
[![MacUnitTest](https://github.com/nao1215/rainbow/actions/workflows/mac_test.yml/badge.svg)](https://github.com/nao1215/rainbow/actions/workflows/mac_test.yml)
[![reviewdog](https://github.com/nao1215/rainbow/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/rainbow/actions/workflows/reviewdog.yml)
[![Gosec](https://github.com/nao1215/rainbow/actions/workflows/security.yml/badge.svg)](https://github.com/nao1215/rainbow/actions/workflows/security.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/rainbow)](https://goreportcard.com/report/github.com/nao1215/rainbow)

## rainbow 
> [!IMPORTANT]  
> This project is under development. Do not use it in production environments.

The rainbow project is a toolset for managing AWS resources. This project consists of multiple CLI and CloudFormation templates. This project adopts README-Driven Development. Therefore, while there may be a README, there might not be any code yet. If you have any feedback regarding the README, please write down in the Issues.

## Supported OS & Go version
- Linux
- Mac
- Windows
- Go 1.19 or later

## CLI List
[WIP]
|Name|README|implementation|Description|
|:--|:--|:--|:--|
|[s3hub](./doc/s3hub/README.md)|✅||User-friendly s3 management tool|
|[spare](./doc/spare/README.md)|✅||Single Page Application Release Easily|

### s3hub example
#### Create a bucket(s)
![create_bucket](./doc/img/s3hub-mb.gif)

#### List buckets
![ls_bucket](./doc/img/s3hub-ls.gif)

#### Remove a bucket with all objects
![rm_bucket](./doc/img/s3hub-rm-all.gif)

#### Interactive mode
![interactive_mode](./doc/img/s3hub-interactive.gif)



## Template List
[WIP]

## LICENSE
This project is licensed under the terms of the MIT license. See the [LICENSE](./LICENSE) file.

## Contributing
Contributions are welcome! Please see the following documents for details:
- [CONTRIBUTING.md](./CONTRIBUTING.md)
- [Code of Conduct](./CODE_OF_CONDUCT.md)
- [Principle](./doc/common/principle.md) 
- [Development Eviorment Setup](./doc/common/developers.md)

This project incurs costs on AWS, and financial support from you would make it easier to maintain the project. If you wish to provide financial support, please do so through [GitHub Sponsors](https://github.com/sponsors/nao1215)

## GitHub Star History
GitHub Star is motivation for me. If you like this project, please star it.
[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/rainbow&type=Date)](https://star-history.com/#nao1215/rainbow&Date)

## Special Thanks
![localstack](./doc/img/localstack-readme-banner.svg)
[LocalStack](https://www.localstack.cloud/) is a service that mocks AWS, covering a wide range of AWS services. It is not easy to set up an AWS infrastructure for personal development, but LocalStack has lowered the barrier for server application development.

It has been incredibly helpful for my technical learning, and among the open-source software (OSS) I encountered in 2023, LocalStack is undoubtedly the best tool. I would like to take this opportunity to express my gratitude.

## Contributors ✨
Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://debimate.jp/"><img src="https://avatars.githubusercontent.com/u/22737008?v=4?s=80" width="80px;" alt="CHIKAMATSU Naohiro"/><br /><sub><b>CHIKAMATSU Naohiro</b></sub></a><br /><a href="https://github.com/nao1215/rainbow/commits?author=nao1215" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!