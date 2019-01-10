# sllm – Structured Logging Leightweight Markup

A human readable approach to make parameters from an actual log message
recognizable for machines.

(Pronounce it like “slim”)

## Rationale
Let's take the example of a standard log message that states some
business relevant event. Besides that message the log entry conatis
some standard fields `time`, `thread`, `level` and `module`.  Further
more we have the transaction id that is added as `tx.id` to the MDC
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

*Note* that `sllm` focuses on the log message.  We do not propose how
to handle other pieces of information! So lets have a look what sllm
would do with some decent log formats.

### Classic log output

```
2018-07-02 20:52:39 [main] INFO sllm.Example - added 7 x Hat to shopping cart by John Doe - tx.id=4711
```

This message is nice to read for humans. Especially the message part
is easy to understand because humans are used to gain some
understanding from natural language. However, the relevant parameters—
i.e. the number of items, which type of item and the respective
user—is not easy to extract from the text by machines. Even if you do,
simple changes of the text template can easily break the mechanism to
extrac those parameters.

**With sllm message**
```
2018-07-02 20:52:39 [main] INFO sllm.Example - added `count:7` x `item:Hat` to shopping cart by `user:John Doe` - tx.id=4711
```

The sllm'ed format is still quite readable but lets one reliably identify the
business relevant values.

### logfmt
```
time=2018-07-02T20:52:39 thread=main level=INFO module=sllm.Example ix.id=4711 number=7 item=Hat user=John_Doe tag=fill_shopping_cart
```

**With sllm message**
```
time=2018-07-02T20:52:39 thread=main level=INFO module=sllm.Example ix.id=4711 msg=added `count:7` x `item:Hat` to shopping cart by `user:John Doe`
```

### JSON Lines
```
{"time":"2018-07-02T20:52:39","thread":"main","level":"INFO","module":"sllm.Example","ix.id"="4711","number":"7","item":"Hat","user":"John Doe","tag":"fill_shopping_cart"}
```

**With sllm message**
```
{"time":"2018-07-02T20:52:39","thread":"main","level":"INFO","module":"sllm.Example","ix.id"="4711","msg":"added `count:7` x `item:Hat` to shopping cart by `user:John Doe`"}
```
