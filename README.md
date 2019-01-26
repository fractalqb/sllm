# sllm – Structured Logging Leightweight Markup
[![Build Status](https://travis-ci.org/fractalqb/sllm.svg)](https://travis-ci.org/fractalqb/sllm)
[![codecov](https://codecov.io/gh/fractalqb/sllm/branch/master/graph/badge.svg)](https://codecov.io/gh/fractalqb/sllm)
[![Test Coverage](https://img.shields.io/badge/coverage-72%25-orange.svg)](file:coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/fractalqb/sllm)](https://goreportcard.com/report/github.com/fractalqb/sllm)
[![GoDoc](https://godoc.org/github.com/fractalqb/sllm?status.svg)](https://godoc.org/github.com/fractalqb/sllm)

A human readable approach to make parameters from an actual log
message recognizable for machines.

(Pronounce it like “slim”) – For a Go logging lib that uses _sllm_ see
[qbsllm](https://github.com/fractalqb/qbsllm).

## Rationale
Logging is subject to two conflicting requirements. Log entries should
be _understandable and easy to read_. On the other hand, they should
be able to be _reliably processed automatically_, i.e. to extract
relevant data from the entry. Current technologies force a decision
for one of the two requirements.

To decide for either human readability or machine processing means
that significant cutbacks are made to the other requirement. _sllm_
picks up the idea of "markup", which is typically intended to bring
human readability and machine processing together. At the same time,
_sllm_ remains simple and unobtrusive—unlike XML ;).

Let's take the example of a standard log message that states some
business relevant event an might be generated from the following
pseudo-code:

```
…
// code that knows something about transaction:
DC.put("tx.id", tx.getId());
…
// down the stack: an actual business related message:
log.Info("added {} x {} to shopping cart by {}", 7, "Hat", "John Doe");
…
// finally(!) when climbing up the stack:
DC.clear();
…
```

Besides the message the log entry contains some standard fields
`time`, `thread`, `level` and `module` and some _Diagnostic Context_
information to give a more complete picture.

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
2018-07-02 20:52:39 [main] INFO sllm.Example - added `count:7` x ↩
↪ `item:Hat` to shopping cart by `user:John Doe` - tx.id=4711
```

The _sllm_'ed format is still quite readable but lets one reliably
identify the business relevant values.

### logfmt
```
time=2018-07-02T20:52:39 thread=main level=INFO module=sllm.Example ↩
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
outstanding advantage: JSON can display structured data. Structured
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

## About the Markup Rules

The markup is simple and sticks to the following requirements:

1. _No support for multi-line messages_

   Spreading a single log entry over multiple lines is considered a
   bad practice. However there may be use cases, e.g. logging a stack
   trace, that justify the multi-line approach. But in any case the
   message of a log entry shall not exceed a single line!
   
2. _Message arguments are unstructured character sequences_

   _sllm_ works on the text level. There is no type system implied by
   _sllm_. As such the arguments of a _sllm_ message are simply
   sub-strings of the message string. 
   
3. _Arguments are identified by a parameter name_

   Within a message each argument is identified by its parameter
   name. A parameter name also is a sub-strings of the message string.
   
4. _Reliable and robust recognition of parameters and arguments_

   The argument and the parameter can be uniquely recognised within a
   message. Changes of a message that do not affect neither the
   parameters nor the arguments do not break the recognition.

5. _Be transparent, simple and unobtrusive_

   A message shall be human readable so that the meaning of the
   message is easy to get. The system must be transparent in the sense
   that even the human reader can easily recognize the parameters with
   their arguments.
   
   _Note that the readability of a message also depends to a certain
   extent on its author._

With this requirements, why was the backtick '`' chosen for markup? –
The backtick is a rarely used character from the ASCII characters set,
i.e. it is also compatible with UTF-8. The fact that it is rarely used
implies that we don't have to escape it often. This affects backticks
in the message template and the arguments. In parameter names
backticks are simply not allowed.

And last but not least: Simpler markup rules make simpler software
implementations ([as long as it is not too
simple](https://en.wikiquote.org/wiki/Albert_Einstein#1930s)). Besides
many advantages this gives room for efficient implementations. Part of
this repository is a [Go
implementation](https://godoc.org/github.com/fractalqb/sllm) that does
not strive so much for efficiency but for hackability.
