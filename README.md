# size

Like `time.Duration` but for bytesizes like `1*Kb` or `16*Gigabyte`. Supports base-2 and base-10 sizes. Conversion to bytes only.

```go
chunkSize = 16 * size.Kib

fmt.Printf("chunk size: %d\n", chunkSize.Bytes()) // chunk size: 16384
```

All other sizes:
```go
sectorSize = 64 * size.Byte
blockSize = 4 * size.Megabyte
partitionSize = 120 * size.Gb
maxUsage = 2 * size.Tib
trafficVolume = 3 * size.Pib
arcticArchiveMaxSize = 5 * size.Eib
```

## Alternatives

If you need to pretty-print sizes, use [go-humanize](https://github.com/dustin/go-humanize#sizes) instead.
