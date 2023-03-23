# Artifact Diff

Compare directories and zip/jar artifacts.

Artifact Diff helps to create reports in plain text, so a generic diff tool can highlight
differences between to paths or archives.

Reports will be written in JSON, and YAML format.

Our use case was a migration of our build environment. We wanted to ensure that the new
artifacts would be equivalent to the previous ones.

## Install

The latest release is available at https://github.com/gesellix/artifact-diff/releases/latest.

You may also install the Golang package like this:

```shell
go install github.com/gesellix/artifact-diff/cmd/artifact-diff@latest
```

## Usage

Please ensure that the binary is executable and in your `$PATH`.

```shell
artifact-diff <report directory> <path1> [path2]
```

- `report directory`: The reports for `path1` and optionally `path2` will be written to the report directory
- `path1`: A directory or zip-compatible archive to be examined
- `path2` (_optional_): Another directory or zip-compatible archive to be examined

## Good to Know

Artifact Diff will extract archives into your `$TEMP` directory, so that archives inside archives
can be examined, too. Please ensure that your storage has enough space for the extracted data.
Temporary files will be removed automatically.

## Credits

The example code from https://github.com/bitfield/tpg-tools helped to get started with Golang's filesystem abstraction.

## License

Copyright (c) 2023 Tobias Gesellchen.

See the LICENSE file in the root directory.
