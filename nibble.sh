#!/bin/sh

# generate for Application.go
# go run nibble.go -a outlet scrollTable1 ApplicationController NSScrollView designable.nib > test.nib
# go run nibble.go -a outlet mainWindow ApplicationController NSWindowTemplate test.nib > test2.nib
# go run nibble.go -a appdelegate ApplicationController test2.nib > test3.nib
# cp test3.nib Application.nib/designable.nib

# generate for HelloWorld.go
# go run nibble.go -a outlet textBox1 ApplicationController NSTextField designable.nib > test.nib
# go run nibble.go -a outlet mainWindow ApplicationController NSWindowTemplate test.nib > test2.nib
# go run nibble.go -a appdelegate ApplicationController test2.nib > test3.nib
# go run nibble.go -a action "buttonClick:" NSButton ApplicationController test3.nib > test4.nib
# cp test4.nib HelloWorld.nib/designable.nib

# generate for SimpleView.go
# go run nibble.go -a outlet itsView Controller SimpleView designable.nib > test.nib
# cp test2.nib SimpleView.nib/designable.nib