# Server selector for gomemcached based on consistent hashing

Provides a ServerSelector interface implementation for picking the next server.
It is safe for conurrent use by multiple go routines.
It uses consistent hashing strategy to distribute keys across nodes

This is just an experimental implementation.
Although it should work, package was not tested in production, so please use with care if intending to use it.

# Inspired by

* [gomemcached client](https://github.com/bradfitz/gomemcache)
* [consistent by StatHat](https://github.com/stathat/consistent)

# About

Written by Ivan Jovanovic at [loopthrough GmbH](http://loopthrough.ch)
