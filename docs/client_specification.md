# Redis Proxy
## Introduction
In this coding assignment, you will be building a transparent Redis proxy service. This proxy is implemented as an HTTP web service which allows the ability to add additional features on top of Redis (e.g. caching and sharding). In the following text, the term “the proxy” refers to this proxy you will create. You have to implement all the requirements below which define the minimally viable deliverable.

This problem is designed to test your abilities to do the kind of systems programming we do at Segment, with a focus on concurrency, networking, integration and some algorithmic optimizations as well as to test your ability to write software that would be maintained and extended by others.

## How to approach this
Where the specification is unclear or falls short, you should make reasonable assumptions and design choices, similar to how you would had this been a project you undertook as part of your regular job. When doing this, is it important to thoroughly document your assumptions and design in the README or other relevant documentation artifacts you choose to produce (e.g. code comments or user manual).

To help you make these decisions, keep in mind that we ask candidates to complete these coding exercises in order to:

1. See how they would implement a software solution based on a problem statement which reflects some aspects of the problems we solve on a
daily basis.
2. Gauge their technical strengths, which we can use for follow-up conversations.

Where a candidate already has a relevant public code portfolio, we often skip this step. We are, therefore, less interested in seeing how well you can stick to every single detail of a detailed specification (although that definitely helps, especially if the specification is clear and unambiguous) as much as we are interested in seeing what you are capable of.

## Evaluation
When we receive your submission, the first thing we’ll do is to unpack the code archive (or git clone it, if appropriate), enter the directory and run make test. The expectation is that, by following the steps above, the code would build itself and run all relevant tests. We expect it to, at least, contain an end-to-end test for each requirement you claim to implement. After successfully running the tests, we’ll review the code and design.

For example, this should “just work”:
```bash
tar -xzvf assignment.tar.gz
cd assignment
make test
```

## Background
Redis describes itself as an “in-memory data structure store” and is deployed as a server process which responds to various text commands. It has an impressive array of client libraries, making it easy to integrate to it from almost any programming language. It stores a variety of data types under string-valued keys, allowing values to be subsequently retrieved and manipulated, using the same key that was used to initially store them.

## Requirements
The table below defines requirements that the proxy has to meet and against which the implementation would be measured. It allows the proxy to be used as a simple read-through cache. When deployed in this fashion, it is assumed that all writes are directed to the backing Redis instance, bypassing the proxy.

TODO