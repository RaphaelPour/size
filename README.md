# size

[![PkgGoDev](https://pkg.go.dev/badge/github.com/raphaelpour/size)](https://pkg.go.dev/github.com/raphaelpour/size)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Like `time.Duration` but for byte sizes like `1*KB` or `16*Gigabyte`. Supports base-2 (IEC) and base-10 (SI) units, string formatting, and text (un)marshaling.

```go
chunkSize := 16 * size.KiB

fmt.Printf("chunk size: %d\n", chunkSize.Bytes()) // chunk size: 16384
```

## Units

Build a size by multiplying an integer with a unit constant. Both the long and short names are available:

```go
sectorSize := 64 * size.Byte
blockSize := 4 * size.Megabyte
partitionSize := 120 * size.GB
maxUsage := 2 * size.TiB
trafficVolume := 3 * size.PiB
arcticArchiveMaxSize := 5 * size.EiB
```

| Base-2 (IEC, ×1024) | Base-10 (SI, ×1000) |
| --- | --- |
| `Byte` / `B` | `Byte` / `B` |
| `Kibibyte` / `KiB` | `Kilobyte` / `KB` |
| `Mebibyte` / `MiB` | `Megabyte` / `MB` |
| `Gibibyte` / `GiB` | `Gigabyte` / `GB` |
| `Tebibyte` / `TiB` | `Terabyte` / `TB` |
| `Pebibyte` / `PiB` | `Petabyte` / `PB` |
| `Exbibyte` / `EiB` | `Exabyte` / `EB` |

The underlying type is `uint64`, so the largest representable unit is the Exbibyte (2^60 bytes). Zebibyte and above are intentionally omitted. Arithmetic is plain `uint64` arithmetic: there are no negative sizes, and overflow wraps rather than saturating.

Read a size back out with `Bytes()` for the exact count, or with the per-unit `float64` accessors (`Kibibytes()`, `Megabytes()`, `Gibibytes()`, …) for a scaled value:

```go
(6 * size.GiB).Gibibytes() // 6
(1500 * size.KB).Megabytes() // 1.5
```

## Formatting

`Size` implements `fmt.Stringer` with an opinionated default: base-2 units with the empty fraction trimmed.

```go
fmt.Println(42 * size.GiB) // 42GiB
fmt.Println(42 * size.GB)  // 39.12GiB
```

For explicit control use the `Format*` methods:

```go
s := 42 * size.GiB

s.FormatIEC()          // "42.00GiB"  — auto-fit to the largest base-2 unit
s.FormatSI()           // "45.10GB"   — auto-fit to the largest base-10 unit
s.Format(size.UnitMiB) // "43008.00MiB" — force a specific unit
```

Options tune precision and trailing zeros:

```go
s.FormatSI(size.WithPrecision(4))          // "45.0972GB"
(10 * size.MiB).FormatIEC(size.WithCutEmptyFraction()) // "10MiB"
```

## Text marshaling

`Size` implements `encoding.TextMarshaler` / `TextUnmarshaler`, so it round-trips through JSON, YAML, and config files. The wire format is the exact byte count with a `B` suffix, which is lossless at every magnitude.

```go
text, _ := (5 * size.MiB).MarshalText() // "5242880B"

var s size.Size
_ = s.UnmarshalText([]byte("5MiB"))     // s.Bytes() == 5242880
```

To parse a string directly, use `Parse` (the `time.ParseDuration` analog); `UnmarshalText` is built on it:

```go
s, _ := size.Parse("5MiB") // s.Bytes() == 5242880
```

Parsing accepts any unit suffix (`"5MiB"`, `"5MB"`, `"5kB"`, `"5242880B"`), tolerates surrounding whitespace, and supports fractional values (`"0.04TiB"`). It is **case-insensitive** (`"5gb"` == `"5GB"`), and the base follows the suffix: `KB`/`kB` are base-10 (1000) while `KiB` is base-2 (1024), so `"1KB"` parses to 1000 bytes, not 1024.

## Alternatives

- [go-humanize](https://github.com/dustin/go-humanize#sizes)

## License

MIT License Copyright (c) 2026 Raphael Pour
