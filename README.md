# structure_fi_coding_challenge

See [PROBLEM.md](PROBLEM.md) for the description of the problem statement.

## Usage
This requires Go 1.17 toolchain to be installed on the machine.

```
$ make && ./structure_fi_coding_challenge
```

In case you have docker and not have Go toolchain installed, you can also use Docker to build and run the project using the following commands

```
# Build the docker image
$ docker build . -t structure_fi_by_ashwanth_kumar
$ docker run -p 8080:8080 -it structure_fi_by_ashwanth_kumar
```

This should start the application on port `8080` as default and we should be consuming the streams on the background.

Once the service starts, you can visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) for a swagger interface to test out the API implementation.

## Solution Notes
> First of all, thank you for the problem statement I had so much fun learning about t-digest and Perfect Hash Functions along the way.

The solution is written in Go and I've written tests pragmatically.

### Median Calculation
Median calculation is backed by [t-digest](https://github.com/tdunning/t-digest/) data structure implementation. I'm using a [go](https://github.com/spenczar/tdigest) implementation in my solution. You can find on [store/store.go](store/store.go).

### Hash Map

> If you follow through the commit history (see `3d12171`, `107122d`, `f344776`), you will notice that I got a working solution (a thin slice) using golang's in-built map data structure with a [Store abstraction](store/store.go) and then eventually replaced it with the custom hash map implementation. This is usually my development style, where I would like to focus on small slices of work and incrementally change / improve the design / code / performance as needed.

Custom Hash Map implementation can be found at [store/custom_map.go](store/custom_map.go). I spent good time going through various [lectures](https://www.youtube.com/watch?v=0M_kIqhwbFo&list=PLUl4u3cNGP61Oq3tWYp6V_F-5jb5L2iHb&index=9), to find an general purpose hash function which will not have any collision and later stumbled upon a blog post that described about Minimal Perfect Hash Function against a known set of keys. Given our keys are the total list of symbols, we can indeed build a MPH using that list. I ended up using the implementation of ["Hash, displace, and compress"](http://cmph.sourceforge.net/papers/esa09.pdf) algorithm by the blog author at https://github.com/dgryski/go-mph. Now given a perfect hash function, I wrote a custom wrapper on top of this to store the values for each key in an array. The maximum memory requirement of this map is `O(n)` where `n` is the total number of symbols. Given the list is finite and known upfront I went with this solution.

Thank you for reading so far, irrespective of the outcome if you've any comments please do send them across.
