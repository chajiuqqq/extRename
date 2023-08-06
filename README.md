# Rename File Extension Name Tool
## Install

    go build -o extChange ./cmd

## Usage

If you want to rename files in `~/movie` from `.7zz` to `.7z`, run this cmd:

    extChange -R 7zz -T 7z -D ~/movie

`-R`: source ext name, supported Regex pattern.

`-T`: new ext name

`-D`: target directory, absPath or relativePath is both supported

## Attention

This operation will overwrite the source file, also is equal to rename file. So be careful.
