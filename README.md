## Bloom
A [bloom filter](https://en.wikipedia.org/wiki/Bloom_filter) implementation using 32bit [fnv](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function) and [murmur3](https://en.wikipedia.org/wiki/MurmurHash) hash functions.

Adding an item to the filter:

    m := 32768
    f := bloom.New(m)
    f.Add([]byte("item"))

Checkng if an item might be in the filter:

    if f.Check([]byte("item"))