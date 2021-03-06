# Changelog

## Unreleased

## [v0.0.9]

- Fix: sometimes before/after words are too big, now they are limited to
       30 characters. Nomenclatural annotations cannot be something like
       "n. n." or "sp. sp." anymore.

## [v0.0.8]

- Add [#20]: Find nomenclatural annotations and save them in output.

## [v0.0.7]

- Add [#19]: Provide sha256 for each title
- Add [#18]: Save given number or words before and after names candidates.

## [v0.0.6]

- Fix [#17]: Output is broken, many records repeat many times.

## [v0.0.5]

- Fix [#16]: Wrong bad pages names detection.
- Fix [#15]: Releases do not show version.

## [v0.0.4]

- Add [#14]: Further performance optimizations.

## [v0.0.3]

- Add [#13]: Speedup by adding Bayes training data preload.
- Add [#12]: Register an error if no pages found in a title.

## [v0.0.2]

- Add [#11]: Flag for number of titles in a progress report line.
- Add [#10]: Add tests that public data from HathiTrust.
- Add [#7]: Save metainformation about titles.
- Fix [#9]: Index is out of range in gnfinder.

## [v0.0.1]

- Add [#6]: prepare Makefile for the first release.
- Add [#5]: documentation in README.md.
- Add [#4]: Save errors and results to filesystem.
- Add [#3]: Find scientific names in provided content.
- Add [#2]: Read content data using data from input file.
- Add [#1]: Read configuration from config file and flags.

## Footnotes

This document follows [changelog guidelines]

[v0.0.8]: https://github.com/gnames/htindex/compare/v0.0.7...v0.0.8
[v0.0.7]: https://github.com/gnames/htindex/compare/v0.0.6...v0.0.7
[v0.0.6]: https://github.com/gnames/htindex/compare/v0.0.5...v0.0.6
[v0.0.5]: https://github.com/gnames/htindex/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/gnames/htindex/compare/v0.0.3...v0.0.4
[v0.0.3]: https://github.com/gnames/htindex/compare/v0.0.2...v0.0.3
[v0.0.2]: https://github.com/gnames/htindex/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/gnames/htindex/compare/v0.0.0...v0.0.1

[#20]: https://github.com/gnames/htindex/issues/20
[#19]: https://github.com/gnames/htindex/issues/19
[#18]: https://github.com/gnames/htindex/issues/18
[#17]: https://github.com/gnames/htindex/issues/17
[#16]: https://github.com/gnames/htindex/issues/16
[#15]: https://github.com/gnames/htindex/issues/15
[#14]: https://github.com/gnames/htindex/issues/14
[#13]: https://github.com/gnames/htindex/issues/13
[#12]: https://github.com/gnames/htindex/issues/12
[#11]: https://github.com/gnames/htindex/issues/11
[#10]: https://github.com/gnames/htindex/issues/10
[#9]: https://github.com/gnames/htindex/issues/9
[#8]: https://github.com/gnames/htindex/issues/8
[#7]: https://github.com/gnames/htindex/issues/7
[#6]: https://github.com/gnames/htindex/issues/6
[#5]: https://github.com/gnames/htindex/issues/5
[#4]: https://github.com/gnames/htindex/issues/4
[#3]: https://github.com/gnames/htindex/issues/3
[#2]: https://github.com/gnames/htindex/issues/2
[#1]: https://github.com/gnames/htindex/issues/1

[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog
