# Hasher
## Background
I've been reading through Alex Xu and Sahm Lam books on system design and decided it would be fun and cool to implement 
a project which could really show off the differences and benefits in using consistent hashing for distributing keys 
across servers in a cluster for something like sharding rather than using something like moding the hash value of the key.

I thought the project would better help solidify the concept as consistent hashing seems to come up several times all 
over the book and can be used in different applications like evenly distributing keys across several nodes in a cluster 
while minimizing the cost of redistributing those keys when the cluster size changes (nodes are added/deleted).

The project really relies on the happy path currently since what I was really interested was looking at
how many keys needed to be redistributed when the cluster was updated (by adding / removing nodes).

## Analysis
Its fairly straight forward to see the benefits of both approaches. Using a hashmod strategy is fairly straight forward
to implement while implementing a hashring strategy is more complex. In terms of redistributing keys the hashring has
the added benefit of only needing to redistribute the keys belonging to a single node (node being deleted or node whos
segment has been partitioned by another node) vs having to redistribute all of the keys for each node.

## How to run it
* First make sure you have the go toolchain setup
* `go build`
* `./hasher`
* Then interact with the Repl given the documentation printed. All logged information for a particular command is contained
within the dashed delimited section:
```
---
put name kevin
... 
...
(logged information for put command)
---
```

## What's next
* Address missing / not existent error handling
* Refactor distribution strategy handling logic (I believe this can be cleaner)
* Add in some tests showcasing redistribution performance as well
* Implement virtual nodes to showcase how we can get a more even distribution to solve for possible hot partition
problems that can occur with the hashring implementation