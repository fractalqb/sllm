/*
Package sllm is the reference implementation for the sllm (pronounced “slim”)
log message format.

sllm is short for Structured Logging Lightweight Markup.  Its goal is
to provide a human readable format for the message part of log-entries
that allows parameters in the log message to be reliably recognized by
programs. A task generally addressed by markup languages. For sllm we
want something much less obstrusive than e.g. XML. The traditional log
message:

  2019/01/11 19:32:44 added 7 ⨉ Hat to shopping cart by John Doe

would become something like (depending on the choice of parameter names)

  2019/01/11 19:32:44 added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`

Still human readable but also easy to be read by machines. Also machine reading
would not break even when the message template changes the order of parameters.
Careful choice of parameter names can make messages even more expressive.

This package is no logging library—it provides functions to create and parse
sllm messages.
*/
package sllm
