gocoa
=====

Objective-C bridge for Go

Rob Pike has gone to great lengths to evangelize Go, and in particular, its current feature set as being carefully designed to leave out those language features that aren't purely orthoganal or that intruduce complexity that bears a cost greater than its utility. So the notion of trying to add a smalltalk style object runtime is likely to upset some people who've embraced the Go idea. On the other hand, how else am I going to write a Mac application?

The code is at this point still rapidly evolving, and as a result does not represent a stable interface as much as a collection of examples and growing set of tests. Given that Go was not designed with either objects or dynamic loading in mind, there's a great deal of conceptual impedence mismatch. Given that the Objective-C runtime has a fairly straightforward C interface, building a library is straightforward, however I think it would probably be easier in the long run to add a new object runtime to Go (Objective-Go?) than to retrofit this existing one.

I am still working out the semantics, fundamental things will change, caveat emptor, etc.