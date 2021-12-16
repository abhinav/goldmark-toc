# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1] - 2021-12-15
### Fixed
- inspect: Correctly handle escaped punctuation in titles.
- render: Don't unintentionally interpret escape sequences in titles.

[0.2.1]: https://github.com/abhinav/goldmark-toc/releases/tag/v0.2.1

## [0.2.0] - 2021-04-04
### Added
- Add `toc.Transformer` to generate a table of contents to the front of any
  document parsed by a Goldmark parser.
- Add `toc.Extender` to extend a `goldmark.Markdown` object with the
  transformer.

[0.2.0]: https://github.com/abhinav/goldmark-toc/releases/tag/v0.2.0

## [0.1.0] - 2021-03-23
- Initial release.

[0.1.0]: https://github.com/abhinav/goldmark-toc/releases/tag/v0.1.0
