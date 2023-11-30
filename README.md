# Key-Value Store with Simple HTTP API

## Introduction

This project implements a persistent key-value store with a simple HTTP API. It exposes the following endpoints:

* `GET http://localhost:8081/get?key=keyName`: Retrieves the value associated with the specified key.
* `POST http://localhost:8081/set`: Sets the value associated with the specified key. The key-value pair is provided in the request body as JSON.
* `DELETE http://localhost:8081/del?key=keyName`: Deletes the specified key and returns its associated value.

The key-value store follows the LSM tree model for reading and writing data. Write operations are first written to the memtable, a sorted map of key-value pairs. The memtable is periodically flushed to disk as an SST file (Sorted String Table). To prevent the number of SST files from growing too large, compaction is performed to merge smaller files into larger ones. In fact, the latter feature is done in parallel with a go routine.

The SST files are in binary format and include the following fields:

* Magic Number: The unique identifier for the application.
* Entry Count: Number of the key-value pairs in the SST File.
* Bloom Filter: The bloom filter's bitset.
* Version: A version number to manage updates to the value.
* Key: The unique identifier for the value.
* Value: The data associated with the key.
* Checksum: A hash value to detect corrupted files.

## Added Dependencies

I have tried to keep third dependencies to a minimum, for I wanted to learn as much as possible from the project. Thus I have not importes any external libraries, and have coded my own Bloom Filter and Red-Black Tree.
The only part that I have found on the internet readily available and have copied is the insertion in the RB Tree (with said rotations).

## Extras

* Bloom filters: Bloom filters are used to quickly test for key existence in SST files.

## Problem Encountered - Wal Cleaning

At first (Refer to previous commits for details), I tried to implement the log file with a watermark. That decision has proven to be the most detrimental to both my project and my sanity. In fact, when renaming the temporary file to the log (supposedly it is atomic on unix based systems but not on windows), I had always gotten an Access Denied Error. I have spent 3 full days trying to debug the problem but to no avail. As such, now I only truncate the log file after flushing. Indeed an expensive approach, and not a standard, but I will try to implement Wal Cleaning correctly later.

## Future Improvements

* **Compression:** SST files are compressed to save disk space.
* **Ensuring Atomicity:** When Flushing, the creation of the SST File is not guaranteed to be atomic, and can lead to bugs when the application crashes (Never happened to me when testing). However, that is only dependent on the Write method of files (Operating System). As such, I am looking for methods to ensure atomicity of writing whole files.
* **Concurrent Distributed Database:** Implement a concurrent distributed database to handle multiple clients and achieve high availability.
* **Performance Enhancement:** Explore techniques to enhance the performance of the key-value store, such as utilizing Goroutines for parallel processing and optimizing data structures.

## Getting Started

To run the key-value store, follow these steps:

1. **Clone the repository.**

2. **In `Lstm.go`, change the constants as you see fit:**
   - `flushThreshold`: The threshold of bytes before flushing.
   - `CompactionThreshold`: The number of SST files before starting to compact.

3. **Start the server:**


You can then access the key-value store using aforementioned API endpoints.