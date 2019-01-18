# sllm – Structured Logging Leightweight Markup
[![Test Coverage](https://img.shields.io/badge/coverage-73%25-orange.svg)](file:coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/fractalqb/sllm)](https://goreportcard.com/report/github.com/fractalqb/sllm)
[![GoDoc](https://godoc.org/github.com/fractalqb/sllm?status.svg)](https://godoc.org/github.com/fractalqb/sllm)

A human readable approach to make parameters from an actual log
message recognizable for machines.

(Pronounce it like “slim”)

## Rationale
Logging is subject to two conflicting requirements. Log entries should
be understandable and easy to read. On the other hand, they should be
able to be reliably processed automatically, i.e. to extract relevant
data from the entry.

Current technologies force a decision for one of the two
requirements. This means that significant cutbacks are made to the
other requirement. _sllm_ picks up the idea of "markup", which is
typically intended to bring human readability and machine processing
together. At the same time, _sllm_ wants to remain simple and
unobtrusive – unlike XML ;).

Let's take the example of a standard log message that states some
business relevant event. Besides that message the log entry conatis
some standard fields `time`, `thread`, `level` and `module`.  Further
more we have the transaction id that is added as `tx.id` to an MDC
because this is some piece of technical information that shall not be
available in the business code but can be immensly helpful as part of
a log entry.

```
// code that knows something about transaction:
MDC.put("tx.id", tx.getId());
…
// down the stack: do an actual business related message:

log.Info("added {} x {} to shopping cart by {}", 7, "Hat", "John Doe");
…
// finally when climbing up the stack:
MDC.clear();
…
```

How would the output of such a scenario look like? – Either one gets
a visually pleasant message that is rather good to read or you get
output that is easy to be consumed by computers. Currently one has
to make a choice. But, can't we have both in one?

**Note** that _sllm_ focuses on the log message.  We do not propose how
to handle other pieces of information! So lets have a look what _sllm_
would do with some decent log formats.

### Classic log output

```
2018-07-02 20:52:39 [main] INFO sllm.Example - added 7 x Hat to ↩
↪ shopping cart by John Doe - tx.id=4711
```

This message is nice to read for humans. Especially the message part
is easy to understand because humans are used to gain some
understanding from natural language. However, the relevant parameters—
i.e. the number of items, which type of item and the respective
user—is not easy to extract from the text by machines. Even if you do,
simple changes of the text template can easily break the mechanism to
extrac those parameters.

**With sllm message:**
```
2018-07-02 20:52:39 [main] INFO sllm.Example - added `count:7` x  ↩
↪ `item:Hat` to shopping cart by `user:John Doe` - tx.id=4711
```

The _sllm_'ed format is still quite readable but lets one reliably
identify the business relevant values.

### logfmt
```
time=2018-07-02T20:52:39 thread=main level=INFO module=sllm.Example  ↩
↪ ix.id=4711 number=7 item=Hat user=John_Doe tag=fill_shopping_cart
```

The [logftm page](https://www.brandur.org/logfmt#human) itself states
that the human readability of logftm is far from perferct and encourages
the approach to include a human readable message with every log line:

```
time=2018-07-02T20:52:39 thread=main level=INFO module=sllm.Example ↩
↪ ix.id=4711 msg="added 7 x Hat to shopping cart by John Doe" ↩
↪ number=7 item=Hat user=John_Doe tag=fill_shopping_cart
```

Once you find the message by skimming the log entry its meaning is not
subject of presonal interpretation of technical key/value pairs any
more. That's fine! But there is still a significant amount of “visual
noise”. However the result is quite acceptable. But one still may ask
if the redundancy in the log entry is necessary. With a “slim” message
it is not.

**With sllm message:**
```
time=2018-07-02T20:52:39 thread=main level=INFO module=sllm.Example ↩
↪ ix.id=4711 msg=added `count:7` x `item:Hat` to shopping cart by ↩
↪ `user:John Doe`
```

### JSON Lines
```
{"time":"2018-07-02T20:52:39","thread":"main","level":"INFO", ↩
↪ "module":"sllm.Example","ix.id"="4711","number":"7","item":"Hat", ↩
↪ "user":"John Doe","tag":"fill_shopping_cart"}
```

Obviously, JSON is the least readable format. However, JSON has an
outstanding advantage. JSON can display structured data. Structured
data is deliberately avoided with _sllm_. Taking this path would
inevitably lead to something with the complexity of XML.

However, similar to the logfmt example, you can use _sllm_ to insert a
machine-readable message into the entry.

**With sllm message:**
```
{"time":"2018-07-02T20:52:39","thread":"main","level":"INFO", ↩
↪ "module":"sllm.Example","ix.id"="4711","msg":"added `count:7` x ↩
↪ `item:Hat` to shopping cart by `user:John Doe`"}
```
