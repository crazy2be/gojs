#include <JavaScriptCore/JSObjectRef.h>
#include <assert.h>
#include <stdio.h>
#include "_cgo_export.h"
#include "callback.h"

FILE* fp = 0;

static JSValueRef JSObjectCallAsFunctionCallback_trampoline(
	JSContextRef ctx, JSObjectRef function, 
	JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], 
	JSValueRef* exception)
{
	void* value = 0;

	// Extract point to go callback
	fprintf( fp, "callback ref = %p, %p\n", function, thisObject ); fflush(fp);
	void* data = JSObjectGetPrivate( function );
	assert( data );
	JSObjectCallAsFunctionCallback_go( data, value );
	return (JSValueRef)value;
}

JSObjectRef JSObjectMakeFunctionWithCallback_wka( JSContextRef ctx, JSStringRef name, void* go_object )
{
	assert( ctx );
	assert( go_object );

	if ( !fp ) {
		fp = fopen( "./tmp.log", "w" );
		assert( fp );
	}

	JSObjectRef ref = JSObjectMakeFunctionWithCallback( ctx, name, JSObjectCallAsFunctionCallback_trampoline );
	fprintf( fp, "ref = %p\n", ref );  fflush(fp);
	if (ref) {
		JSObjectSetPrivate( ref, go_object );
		assert( JSObjectGetPrivate(ref) == go_object );
	}
	return ref;
}
	
