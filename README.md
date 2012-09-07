gocoa
=====

Objective-C bridge for Go

status
------

You can get a quick synopsis of the current status with:

	go test gocoa

The package isn't yet go-gettable because there are crashers and missing fundamental types. Just download the zipfile and set a local GOPATH to the current working directory. I don't want you to put this code in your pkg directory, it's ready to play with, but not ready for work.

The Objective-C runtime interface I've defined is settling down now, but there are some details I want to consider before baking anything else in.

A problem with the Go linker isn't a showstopper, but requires some dancing around with export names. I know that substantial changes to the linker are coming in a future Go version anyway, so a fix may be forthcoming besides.


examples
--------

* HelloWorld.go demonstrates messaging and events between Go code and Interface Builder controls
* SimpleView.go displays a subclass NSView and responds to drawRect: requests from NSApplication
* TableView.go displays and adds arbitrary data to a NSTableView defined in Interface Builder (broken)
* Nibless.go defines its UI in code rather than using Interface Builder
* WebKit? forthcoming...


To run any of the examples use, for example:

	export GOPATH=`pwd` 
	go build -o hello HelloWorld.go
	./hello 


nibble
------

Nibble edits .nib files and prints the result to stdout. .nib files are, conveniently, these days .xib files defined in XML. What this buys you is the ability to add IB outlets, actions and delegates without creating an Objective-C class, writing a header, defining your outlets programmatically, and then letting Interface Builder import them for you. That's tedious and stupid when you know what you're doing.

Like a cowboy, you can just draw your app in IB, drag some custom objects on the thing, connect them with something you define later in Go, and smash them together with nibble. For example:

	nibble -a outlet mainWindow ApplicationController NSWindowTemplate designable.nib > out.nib

Adds an outlet, 'mainWindow' from 'ApplicationController' to NSWindowTemplate. Or, more commonly:

	nibble -a appdelegate ApplicationController designable.nib > out.nib

Connects 'ApplicationController' as delegate to NSApplication, allowing it to receive and handle events.

If you don't like bugs, don't use nibble just yet. There are some limitations in the feature set, error handling isn't robust, and I'd like to write some formal tests before calling it finished.


license
-------

Three-clause BSD, attribution with liberal redistribution rights. The Gopher icns is a derived work and CC-BY Renee French & the Go Authors.


reasoning
---------

Rob Pike has gone to great lengths to evangelize Go, and in particular, its current feature set as being carefully designed to leave out those language features that introduce complexity that bears a cost greater than their utility. So the notion of trying to add a smalltalk style object runtime is likely to upset some people who've embraced the Go idea. On the other hand, how else am I going to write a Mac application?

Given that Go was not designed with either objects or dynamic loading in mind, there's some conceptual impedance mismatch. However, the Objective-C runtime has a fairly straightforward C interface, and bridging that interface is certainly enough to use Cocoa effectively from within Go.


scope
-----

With gocoa, I'm working toward a solution that allows me to write less code, without generating runtime dependencies on additional libraries and tools. I'm also interested in things like the iPhone and GNUStep and Étoilé, so while I've seen other work in this area, I'm solving for the general case.

I looked at several Objective-C bridges before choosing a strategy, and my thinking was most influenced by Cocoa# from the mono project. Although C# has all sorts of constructs irrelevant to Go, the convention is to simply build out from the C interface. That's a heck of a lot less to debug.

In studying other projects, I noticed that in each of them a lot of unnecessary peer code was written, and that looks like a maintenance headache that will eventually lead to bit rot. If you must have peers, use a tool to generate lightweight ones using introspection at runtime in your target language. Don't write them, you're probably just wasting time handcrafting bugs.

I think it would probably be easier to add a new object runtime to Go (Objective-Go?) than to retrofit this existing one. I've read some papers on very small runtimes that are interesting. OTOH, I don't yet have a strong argument for one.

