# bedrockbuf
The Minecraft protocol implemented in Google's protobuf format 

## This was a massive waste of time
When starting this project, I expected protobufs to be a lot more efficient and optimized than the bedrock protocol, but I was wrong... I spent a whole day writing the generator for the entire protocol without testing a single packet first, and I extremely regret that now.
```
goos: windows
goarch: amd64
pkg: github.com/twistedasylummc/bedrockbuf
cpu: AMD Ryzen 7 3700X 8-Core Processor
BenchmarkProtocol
BenchmarkProtocol-16             8664360               140.3 ns/op
BenchmarkProtobuf
BenchmarkProtobuf-16               21295             56844 ns/op
PASS
```
The benchmarks above were taken using LevelChunk packets with identical content. Already the protobuf format is ~405x slower in this scenario. With the same packet, the final output from protobuf was 798 bytes larger than the bedrock protocol, 9788 vs 10586. This is mainly due to the lack of byte types in protobuf, as bytes are written as int32s/uint32s.