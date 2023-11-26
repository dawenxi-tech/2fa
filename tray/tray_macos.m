
// +build darwin,!ios

#include <AppKit/AppKit.h>
#include "_cgo_export.h"

const int cellWidth = 200;
const int cellHeight = 60;


@interface CodeView : NSView

@property (strong) NSString * secret;
@property (strong) NSTextField* code;
@property(strong) NSTextField * name;

@end


@implementation CodeView


- (id)init {
    self = [super init];

    self.name = [NSTextField labelWithString:@""];
    self.name.frame = CGRectMake(0, 40, cellWidth, cellHeight-40);
    self.name.wantsLayer = YES;
    [self addSubview:self.name];


    self.code = [NSTextField labelWithString:@""];
    self.code.frame = CGRectMake(0, 0, cellWidth, 40);
    self.code.wantsLayer = YES;
    self.code.layer.backgroundColor = [[NSColor greenColor] CGColor];
    [self addSubview:self.code];

    return self;
}

- (void) mouseDown:(NSEvent *)event {
    NSLog(@"mouse down %@", event);
}

- (void) mouseUp:(NSEvent *)event {
    NSLog(@"mouse up %@", event);
}

@end


@interface CodesView : NSView

@property (nonatomic, strong) NSArray<CodeView *> * codes;

@end

@implementation CodesView


- (void) setCodes:(NSArray<CodeView *> *)codes {
    [self cleanSubview];
    _codes = codes;
    [self layout];
}


- (void) layout {
    int oy = 0;
    for (CodeView *codeView in self.codes) {
        codeView.frame = CGRectMake(0, oy, cellWidth, cellHeight);
        [self addSubview:codeView];
        oy += cellHeight;
    }
}

- (void) cleanSubview {
    for (NSView* view in self.subviews){
        [view removeFromSuperview];
    }
}

- (void) refreshCodes {
    for (CodeView *cv in self.codes) {
        NSString *code = CFBridgingRelease(export_2fa_code((__bridge CFTypeRef)cv.secret));
        cv.code.stringValue = code;
    }
}

@end


@interface PopoverController : NSViewController

@property (nonatomic, strong) NSArray<NSDictionary*>* data;
@property (nonatomic, strong) CodesView *codes;

@end

@implementation PopoverController

- (void) viewDidAppear {
    NSLog(@"PopoverController viewDidAppear");
    [self loadViews];
}


- (void) loadViews {
    if (self.codes != nil) {
        return;
    }
    self.codes = [[CodesView alloc] initWithFrame:CGRectMake(0, 0, cellWidth, 400)];
    NSMutableArray<CodeView*> *codesView = [NSMutableArray array];
	for (NSDictionary *dic in self.data) {
        CodeView *code = [[CodeView alloc] init];
        code.secret = [dic objectForKey:@"secret"];
        code.code.stringValue = [dic objectForKey: @"name"];
        code.name.stringValue = [dic objectForKey: @"name"];
        [codesView addObject:code];
	}
    self.codes.codes = codesView;
    [self.view addSubview:self.codes];

    [self.codes refreshCodes];
}

@end


@interface PopoverManager : NSObject

@property (nonatomic, strong) NSPopover *popover;

+ (void)show:(NSStatusBarButton *) button;

@end

@implementation PopoverManager

+ (void)show:(NSStatusBarButton *) button {
    PopoverController *vc = [[PopoverController alloc] init];
    vc.view = [[NSView alloc] init];
    NSArray *codes = CFBridgingRelease(export_codes());
    vc.data = codes;
    int len = [codes count];
    int height = cellHeight * len;
    
    NSPopover* popover = [[NSPopover alloc] init];
    popover.contentSize = CGSizeMake(cellWidth, height);
    popover.behavior = NSPopoverBehaviorTransient;
    popover.appearance = [NSAppearance appearanceNamed:NSAppearanceNameVibrantLight];
    popover.contentViewController = vc;
    [popover showRelativeToRect:button.bounds ofView:button preferredEdge:NSRectEdgeMaxY];
}


@end

void show_tray(void) {

dispatch_async(dispatch_get_main_queue(), ^{
    NSLog(@"inside dispatch async block main thread from main thread");
	id delegate = [[NSApplication sharedApplication] delegate];
	NSStatusItem* statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
    statusItem.button.title = @"2FA";
	id obj = [PopoverManager class];
    SEL mySelector = @selector(show:);
	statusItem.button.target = obj;
    statusItem.button.action = mySelector;
	statusItem.visible = YES;
	[delegate performSelector:@selector(setStatusItem:) withObject:statusItem];
	NSLog(@"delegate: %@", delegate);
});

}