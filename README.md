# gocompat [![Build Status](https://travis-ci.org/s2gatev/gocompat.svg?branch=master)](https://travis-ci.org/s2gatev/gocompat) [![Coverage Status](https://coveralls.io/repos/s2gatev/gocompat/badge.svg?branch=master&service=github)](https://coveralls.io/github/s2gatev/gocompat?branch=master)

Backwards compatibility checker for Go APIs.

## Introduction

**gocompat** allows you to verify backwards compatibility of your project interface.
It stores an index of all exported symbols in a `.gocompat` file that allows comparisons with
newer versions of the interfaces at later point. 

## Installation

`go get -u github.com/s2gatev/gocompat`

## Usage

Execute `gocompat` inside your project directory. You can modify the command by inserting:
* `-f` for storing the current interface in the index even if it is not compatible with the previous one.

## TODO

A list of things that should be taken care of:
* Handle methods for structs.
* Handle interfaces and their methods.
* Handle nested packages property. `./a/test` is different thank `./b/test`.
* Handle interface conversion - stricker to more relaxed interface should not break compatibility.

## Contribution

Contributions via issues and pull requests are more than welcome!

## License

gocompat is licensed under the [MIT License](LICENSE).
