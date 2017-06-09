# Server selector for gomemcached based on consistent hashing

Provides a ServerSelector interface implementation for picking the next server.
It is safe for conurrent use by multiple go routines.
It uses consistent hashing strategy to distribute keys across nodes

## Libraries

Uses GoDS TreeMap for keeping the sorted map of hashes -> servers.
https://github.com/emirpasic/gods#treemap
