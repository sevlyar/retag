# Retag [![TravisCI](https://api.travis-ci.org/sevlyar/retag.svg)](https://travis-ci.org/sevlyar/retag) [![GoDoc](https://godoc.org/github.com/sevlyar/retag?status.svg)](https://godoc.org/github.com/sevlyar/retag)

Package retag provides an ability to change tags of structures' fields in runtime
without copying of the data. It may be helpful in next cases:

* Automatic tags generation;
* Different views of the one data;
* Fixing of leaky abstractions with minimal boilerplate code
when application has layers of abstractions and model is
separated from storages and presentation layers.

Features:

* No memory allocations (for cached types);
* Fast converting (lookup in table and pointer creation for cached types).

The package requires go1.7+. 

**The package is still experimental and subject to change. The package can be broken by a next release of go.**

## Installation

    go get github.com/sevlyar/retag

You can use [gopkg.in](http://labix.org/gopkg.in):

    go get gopkg.in/sevlyar/retag.v0

If you want to use the library in production project, please use vendoring,
because i can not ensure backward compatibility before release v1.0.

## Documentation

Please see [godoc.org/github.com/sevlyar/retag](https://godoc.org/github.com/sevlyar/retag)
