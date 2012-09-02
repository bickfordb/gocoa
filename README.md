gocoa
=====

Objective-C bridge for Go

status
------

The code is still rapidly evolving, and as a result does not represent a stable interface. I am still working out the semantics, fundamental things will change, caveat emptor, etc. It is however a fairly concise example of how to go about interfacing the Objective-C runtime.


examples
--------

* HelloWorld.go demonstrates messaging and events between Go code and Interface Builder controls.
* SimpleView.go (broken) displays a subclass NSView defined in Interface Builder and responds to drawRect: requests from NSApplication.
* TableView.go (incomplete) displays and adds data to a NSTableView defined in Interface Builder.
* Nibless.go (incomplete) defines its UI in code rather than using Interface Builder.

...and what I really should be working on is getting an example using WebKit working. That's where Go's rubber meets the road.


To run any of the examples use, for example:

	export GOPATH=`pwd` 
	go run HelloWorld.go 


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

Whether Go needs OOP features or not is the subject of plenty of mailing list gripes. My thinking is: it's unreasonable to ask somebody else to build and maintain something that you want for yourself, and similarly unreasonable to ask someone to understand and then maintain something for you. 
