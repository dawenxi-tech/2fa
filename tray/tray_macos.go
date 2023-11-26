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

static void nsstringGetCharacters(CFTypeRef cstr, unichar *chars, NSUInteger loc, NSUInteger length) {
	NSString *str = (__bridge NSString *)cstr;
	[str getCharacters:chars range:NSMakeRange(loc, length)];
}

static NSUInteger nsstringLength(CFTypeRef cstr) {
	NSString *str = (__bridge NSString *)cstr;
	return [str length];
}

*/
import "C"
import (
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/xlzd/gotp"
	"unicode/utf16"
	"unsafe"
)

func show_tray() {
	C.show_tray()
}

//export export_codes
func export_codes() C.CFTypeRef {
	arr := C.newNSArray()
	codes := storage.LoadCodes()
	for _, code := range codes {
		C.array_add_object(arr, stringToNSString(code.Name), stringToNSString(code.Secret.Val()))
	}
	return arr
}

//export export_2fa_code
func export_2fa_code(str C.CFTypeRef) C.CFTypeRef {
	secret := nsstringToString(str)
	totp := gotp.NewDefaultTOTP(secret)
	return stringToNSString(totp.Now())
}

// --- Copy From Gio ---

// stringToNSString converts a Go string to a retained NSString.
func stringToNSString(str string) C.CFTypeRef {
	u16 := utf16.Encode([]rune(str))
	var chars *C.unichar
	if len(u16) > 0 {
		chars = (*C.unichar)(unsafe.Pointer(&u16[0]))
	}
	return C.newNSString(chars, C.NSUInteger(len(u16)))
}

// nsstringToString converts a NSString to a Go string.
func nsstringToString(str C.CFTypeRef) string {
	if str == 0 {
		return ""
	}
	n := C.nsstringLength(str)
	if n == 0 {
		return ""
	}
	chars := make([]uint16, n)
	C.nsstringGetCharacters(str, (*C.unichar)(unsafe.Pointer(&chars[0])), 0, n)
	utf8 := utf16.Decode(chars)
	return string(utf8)
}
