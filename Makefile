include $(GOROOT)/src/Make.inc

TARG=javascriptcore
CGOFILES=\
	base.go \
	context.go \
	native.go \
	object.go \
	panic.go \
	reflect.go \
	string.go \
	value.go 
CGO_OFILES=callback.o
CGO_CFLAGS=`pkg-config --cflags webkit-1.0`
CGO_LDFLAGS=`pkg-config --libs webkit-1.0`

include $(GOROOT)/src/Make.pkg

main: install main.go
	$(GC) main.go
	$(LD) -o $@ main.$O
