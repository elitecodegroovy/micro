# BuffIO

Package bufio implements buffered I/O. It provides buffering and some help for textual I/O. The default buffer size is 4096.
Scanner provides a convenient interface for reading data such as a file of newline-delimited lines of text. 
Sequential scans on a reader should use `bufio.Reader`, or you should use `bufio.NewScanner(reader)`

Use Scan case: scan the reader and the handle the content.

```
    buf := strings.NewReader(reader)
    s := bufio.NewScanner(buf)
    s.Split(bufio.ScanBytes)
```