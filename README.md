gocoa
=====

Objective-C bridge for Go

status
------

The code is still rapidly evolving, and as a result does not represent a stable interface. I am still working out the semantics, fundamental things will change, caveat emptor, etc. It is however a fairly concise example of how to go about interfacing the Objective-C runtime.


examples
--------

To run examples use, for example:

	export GOPATH=`pwd` 
	go run HelloWorld.go 


license
-------

Three-clause BSD, attribution with liberal redistribution rights.


reasoning
---------

Rob Pike has gone to great lengths to evangelize Go, and in particular, its current feature set as being carefully designed to leave out those language features that introduce complexity that bears a cost greater than their utility. So the notion of trying to add a smalltalk style object runtime is likely to upset some people who've embraced the Go idea. On the other hand, how else am I going to write a Mac application?

Given that Go was not designed with either objects or dynamic loading in mind, there's a great deal of conceptual impedence mismatch. However, the Objective-C runtime has a fairly straightforward C interface. Even still, I think it would probably be easier to add a new object runtime to Go (Objective-Go?) than to retrofit this existing one. I've read some papers on small runtimes for GNUStep that are interesting (the perennial problem that new Objective-C implementations face is how to get garbage collection). 

As an aside, whether Go needs OOP features or not is the subject of all these mailing list gripes. My thinking is, why are these people always begging the implementors for OOP features? The runtime source is available, the thing is well understood, it's a neat pet project to build for yourself.

I looked at several Objective-C bridges before choosing a strategy, and my thinking was most influenced by Cocoa#. Although C# supports all sorts of constructs irrelevant to Go, and although Cocoa# reimplements features everywhere with a peer infrastructure, the convention is to build out from the C interface to the Objective-C runtime. When I looked at other projects, PyObjC in particular, a lot of unnecessary peer code is being written in Objective-C, and that looks like it might have been the maintenance headache that eventually led to bit rot.

With gocoa, I'm working toward solutions to similar problems that allow me to write less code, without generating additional dependencies, and to write everything in Go (failing that, Go and cgo).

