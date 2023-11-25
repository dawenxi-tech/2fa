//go:build darwin && !ios
// +build darwin,!ios

package tray

/*
#cgo CFLAGS: -Werror -Wno-deprecated-declarations  -DDARWIN -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework WebKit -framework AppKit -framework QuartzCore

#include <AppKit/AppKit.h>

__attribute__ ((visibility ("hidden"))) void show_tray(void);

static CFTypeRef newNSString(unichar *chars, NSUInteger length) {
	@autoreleasepool {
		NSString *s = [NSString string];
		if (length > 0) {
			s = [NSString stringWithCharacters:chars length:length];
		}
		return CFBridgingRetain(s);
	}
}

static CFTypeRef newNSArray() {
	@autoreleasepool {
		NSMutableArray* arr = [NSMutableArray array];
		return CFBridgingRetain(arr);
	}
}

static void array_add_object(CFTypeRef arr, CFTypeRef name, CFTypeRef secret) {
	NSMutableArray*_arr = (__bridge NSMutableArray*)arr;
	NSString*_name = (__bridge NSString*)name;
	NSString*_secret = (__bridge NSString*)secret;
	[_arr addObject:@{@"name": _name, @"secret": _secret}];
}


*/
import "C"

func show_tray() {
	C.show_tray()
}
