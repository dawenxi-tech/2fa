//go:build darwin && !ios
// +build darwin,!ios

package tray

/*
#cgo CFLAGS: -Werror -Wno-deprecated-declarations  -DDARWIN -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework WebKit -framework AppKit -framework QuartzCore

#include <AppKit/AppKit.h>

__attribute__ ((visibility ("hidden"))) void show_tray(void);
__attribute__ ((visibility ("hidden"))) void dismiss_tray(void);

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

static CFTypeRef nsimageWithData(const char* iconBytes, int length) {
@autoreleasepool {
    NSData* buffer = [NSData dataWithBytes:iconBytes length:length];
    NSImage *image = [[NSImage alloc] initWithData:buffer];
	return CFBridgingRetain(image);
}
}

static void windowMakeKeyAndOrderFront() {
dispatch_async(dispatch_get_main_queue(), ^{
	//NSLog(@"windows -> %@", [NSApp windows]);
    NSWindow *win = [NSApp mainWindow];
	//NSLog(@"win -> %@", win);
	[win makeKeyAndOrderFront:nil];
});
}

*/
import "C"
import (
	_ "embed"
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/xlzd/gotp"
	"golang.design/x/clipboard"
	"unicode/utf16"
	"unsafe"
)

//go:embed settings.png
var settingsIcon []byte

//go:embed dashboard.png
var dashboardIcon []byte

func show_tray() {
	C.show_tray()
}

func dismiss_tray() {
	C.dismiss_tray()
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

//export code_on_click
func code_on_click(str C.CFTypeRef) {
	secret := nsstringToString(str)
	totp := gotp.NewDefaultTOTP(secret)
	code := totp.Now()
	clipboard.Write(clipboard.FmtText, []byte(code))
}

//export export_settings_icon
func export_settings_icon() C.CFTypeRef {
	cstr := (*C.char)(unsafe.Pointer(&settingsIcon[0]))
	return C.nsimageWithData(cstr, (C.int)(len(settingsIcon)))
}

//export export_dashboard_icon
func export_dashboard_icon() C.CFTypeRef {
	cstr := (*C.char)(unsafe.Pointer(&dashboardIcon[0]))
	return C.nsimageWithData(cstr, (C.int)(len(dashboardIcon)))
}

//export tray_button_on_click
func tray_button_on_click(typ C.int) {
	var t = int(typ)
	sendEvent(EventType(t))
}

func bring_window_to_front() {
	C.windowMakeKeyAndOrderFront()
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
