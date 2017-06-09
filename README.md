# Server selector for gomemcached based on consistent hashing

Provides a ServerSelector interface implementation for picking the next server.
It is safe for conurrent use by multiple go routines.
It uses consistent hashing strategy to distribute keys across nodes
