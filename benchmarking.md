Benchmark these:
* no concurrency
* multiple readers, one writer
* multiple readers, multiple writers

### Architecture (multiple writers, multiple readers)
**writers send messages to channel**
channel subscribers: function that writes to file, function that updates keydir

**readers send messages to channel?**
or do they just read from keydir/file drectly?

### TODO
* - [ ] move keydir updater into its own function (package?)
* - [ ] move file thing into its own package
