//go:build darwin && !ios
// +build darwin,!ios

package tray

/*
#cgo CFLAGS: -Werror -Wno-deprecated-declarations  -DDARWIN -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework WebKit -framework AppKit -framework QuartzCore

#include <AppKit/AppKit.h>

@interface Popover : NSObject
+(void)show:(id) obj;

@end

@implementation Popover

+(void)show:(NSStatusBarButton *) button
{
    NSViewController *vc = [[NSViewController alloc] init];
    vc.view = [[NSView alloc] init];
NSPopover* popover = [[NSPopover alloc]init];
popover.behavior = NSPopoverBehaviorTransient;
popover.appearance = [NSAppearance appearanceNamed:NSAppearanceNameVibrantLight];
popover.contentViewController = vc;
[popover showRelativeToRect:button.bounds ofView:button preferredEdge:NSRectEdgeMaxY];
}

@end

static void show_tray(void) {
dispatch_async(dispatch_get_main_queue(), ^{
	id delegate = [[NSApplication sharedApplication] delegate];
	NSStatusItem* statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
    statusItem.button.title = @"2FA";
	id obj = [Popover class];
    SEL mySelector = @selector(show:);
	statusItem.button.target = obj;
    statusItem.button.action = mySelector;
	statusItem.visible = YES;
 	[delegate performSelector:@selector(setStatusItem:) withObject:statusItem];
});
}

*/
import "C"

func show_tray() {
	C.show_tray()
}
