/*
Package sllm is the 3rd iteration of the reference implementation for
the Structured Logging Lightweight Markup format.

The goal is to create a human-readable format for the message part of
log entries that also makes the parameters in the log message reliably
machine-readable. This is a task generally performed by markup
languages. However, sllm is intended to be much less intrusive than,
for example, XML or JSON. The traditional log message:

	2019/01/11 19:32:44 added 7 ⨉ Hat to shopping cart by John Doe

would become something like (depending on the choice of parameter names)

	2019/01/11 19:32:44 added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`

This is still readable by humans but can also be parsed reliably by
machines. Machine reading would not fail even if the message template
changes the order of the parameters. Careful choice of parameter names
can make the messages even more meaningful.

This package is no logging library—it provides functions to create and parse
sllm messages.
*/
package sllm
