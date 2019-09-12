# htindex

The purpose of `htindex` is to create an index of scientific names in
HathiTrust Digital Library. This library contains large amount of scientific
literature (40% public, 60% private). This program will allow to add
biodiversity information to their metadata. It will make possible to search
their corpus by scientific names.

- [htindex](#htindex)
  - [Installation](#installation)
  - [Usage](#usage)
  - [License](#license)
  - [Authors](#authors)

## Installation

For Linux download the [latest release], untar, and copy it to /usr/local/bin,
or any other directory that is in the PATH.

In the home directory of the user that will run the app create
`.htindex.yaml`. Use an [example .htindex.yaml file] for reference. The
example file explains configuration parameters. You can skip creation of the
`.htindex.yaml` file, if you are planning to provide all the needed settings
via command line flags.

## Usage

The `htindex` reads a file that contains paths to zipped
`volumes/books/titles` from HathiTrust, finds these files, extracts text from
all the pages, finds scientific names in them and saves results to a given
output directory.

If `~/.htindex.yaml` file already contains all the settings it is sufficient
to run

```bash
htindex
# To see help message:
htindex -h
# To see version of the app:
htindex -v
```

If some settings for the app need to be modified during command line
execution, use the following flags:

`-h, --help`
:help for htindex

`-j, --jobs`
:takes an positive integer. Sets the number of workers (jobs). It looks like
optimal number is `number_of_threads * 3`.

`-i, --input`
:takes a string. Sets a path to the input data file

`-o, --output`
: takes a string. Sets a path to the output directory. This directory will
contain error log and results data.

`-r, --root`
: takes a string. Sets a root path to add to the input file data. This creates
complete absolute path to zip files with volumes.

`-v, --version`
: Shows htindex version and build timestamp

## License
Released under [MIT license]

## Authors

- [Dmitry Mozzherin]

[Dmitry Mozzherin]: https://gitlab.com/dimus
[example .htindex.yaml file]: https://raw.githubusercontent.com/gnames/htindex/master/files/.htindex.yaml
[MIT license]: https://raw.githubusercontent.com/gnames/htindex/master/LICENSE
[latest release]: https://github.com/gnames/htindex/releases/latest