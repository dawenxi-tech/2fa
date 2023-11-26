
// +build darwin,!ios

#include <AppKit/AppKit.h>
#include "_cgo_export.h"

const int cellWidth = 200;
const int cellHeight = 60;
const int cellSpace = 10;
const int paddingHorizontal = 10;

@protocol PopoverManagerDelegate <NSObject>
- (void) closePopover;
@end

@interface VerticalCenterCell : NSTextFieldCell
@end

@implementation VerticalCenterCell
- (NSRect)drawingRectForBounds:(NSRect)theRect{
    NSRect newRect = [super drawingRectForBounds:theRect];
    NSSize textSize = [self cellSizeForBounds:theRect];
    float heightDelta = newRect.size.height - textSize.height;
    if (heightDelta > 0){
        newRect.size.height -= heightDelta;
        newRect.origin.y += round(heightDelta / 2);
    }
    return newRect;
}
- (void)selectWithFrame:(NSRect)aRect inView:(NSView *)controlView editor:(NSText *)textObj delegate:(nullable id)anObject start:(NSInteger)selStart length:(NSInteger)selLength{
    aRect = [self drawingRectForBounds:aRect];
    [super selectWithFrame:aRect inView:controlView editor:textObj delegate:anObject start:selStart length:selLength];
}
@end

@interface CodeView : NSView

@property (strong) NSString * secret;
@property (strong) NSTextField* code;
@property(strong) NSTextField * name;
@property (nonatomic, weak) id<PopoverManagerDelegate> popover;

@end


@implementation CodeView


- (id)init {
    self = [super init];

    self.name = [[NSTextField alloc] init];
    self.name.frame = CGRectMake(10, 40, cellWidth-20 - 2 * paddingHorizontal, cellHeight-40);
    VerticalCenterCell *nameCell = [[VerticalCenterCell alloc] init];
    nameCell.stringValue = @"";
    nameCell.editable = NO;
    nameCell.scrollable = NO;
    nameCell.alignment = NSTextAlignmentLeft;
    nameCell.bordered = NO;
    nameCell.textColor = [NSColor grayColor];
    self.name.cell = nameCell;
    [self addSubview:self.name];

    self.code = [[NSTextField alloc] init];
    self.code.frame = NSMakeRect(0, 0, cellWidth-paddingHorizontal * 2, 40);
    VerticalCenterCell *cell = [[VerticalCenterCell alloc] init];
    cell.stringValue = @"";
    cell.editable = NO;
    cell.scrollable = NO;
    cell.alignment = NSTextAlignmentCenter;
    cell.bordered = NO;
    cell.font = [NSFont monospacedSystemFontOfSize:28 weight:600];
    self.code.cell = cell;
    self.code.wantsLayer = YES;
    self.code.usesSingleLineMode = YES;

    [self addSubview:self.code];

    self.wantsLayer = YES;
    self.layer.backgroundColor = [[NSColor whiteColor] CGColor];

    return self;
}

- (void) mouseDown:(NSEvent *)event {
    // todo Material Ripple effect
    // NSLog(@"mouse down %@", event);
    code_on_click((__bridge CFTypeRef)self.secret);
    if (self.popover != nil) {
        [self.popover closePopover];
    }
}

- (void) mouseUp:(NSEvent *)event {
    // todo Material Ripple effect
    // NSLog(@"mouse up %@", event);
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
    int oy = cellSpace / 2;
    for (CodeView *codeView in self.codes) {
        codeView.frame = CGRectMake(paddingHorizontal, oy, cellWidth-2*paddingHorizontal, cellHeight);
        [self addSubview:codeView];
        oy += cellHeight + cellSpace;
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


@interface PopoverController : NSViewController <PopoverManagerDelegate>

@property (nonatomic, strong) NSArray<NSDictionary*>* data;
@property (nonatomic, strong) CodesView *codes;
@property (nonatomic, weak) NSPopover *popover;

@end

@implementation PopoverController

- (void) viewDidAppear {
//     NSLog(@"PopoverController viewDidAppear");
    [super viewDidAppear];
    [self loadViews];
}

- (void) viewDidDisappear {
//     NSLog(@"PopoverController viewDidDisappear");
    [super viewDidDisappear];
}

- (void) loadViews {
    if (self.codes != nil) {
        return;
    }
    int len = [self.data count];
    int height = (cellHeight+cellSpace) * len;

    self.codes = [[CodesView alloc] initWithFrame:CGRectMake(0, 0, cellWidth, height)];
    NSMutableArray<CodeView*> *codesView = [NSMutableArray array];
	for (NSDictionary *dic in self.data) {
        CodeView *code = [[CodeView alloc] init];
        code.popover = self;
        code.secret = [dic objectForKey:@"secret"];
        code.name.stringValue = [dic objectForKey: @"name"];
        [codesView addObject:code];
	}
    self.codes.codes = codesView;
    [self.view addSubview:self.codes];

    [self.codes refreshCodes];
}

- (void) closePopover {
    if (self.popover) {
        [self.popover close];
    }
}

@end


@interface PopoverManager : NSObject

@property (nonatomic, strong) NSPopover *popover;

+ (void)show:(NSStatusBarButton *) button;

@end

@implementation PopoverManager

+ (void)show:(NSStatusBarButton *) button {
    NSArray *codes = CFBridgingRelease(export_codes());
    int len = [codes count];
    int height = (cellHeight+cellSpace) * len;

    PopoverController *vc = [[PopoverController alloc] init];
    vc.view = [[NSView alloc] init];
    vc.data = codes;

    NSPopover* popover = [[NSPopover alloc] init];
    popover.contentSize = CGSizeMake(cellWidth, height);
    popover.behavior = NSPopoverBehaviorTransient;
    popover.appearance = [NSAppearance appearanceNamed:NSAppearanceNameVibrantLight];
    popover.contentViewController = vc;
    vc.popover = popover;
    [popover showRelativeToRect:button.bounds ofView:button preferredEdge:NSRectEdgeMaxY];

    // auto close popover when click out bounds;
    [[[[popover contentViewController] view] window] makeKeyWindow];
}


@end

void show_tray(void) {

dispatch_async(dispatch_get_main_queue(), ^{
//     NSLog(@"inside dispatch async block main thread from main thread");
	id delegate = [[NSApplication sharedApplication] delegate];
	NSStatusItem* statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
    statusItem.button.title = @"2FA";
	id obj = [PopoverManager class];
    SEL mySelector = @selector(show:);
	statusItem.button.target = obj;
    statusItem.button.action = mySelector;
	statusItem.visible = YES;
	[delegate performSelector:@selector(setStatusItem:) withObject:statusItem];
// 	NSLog(@"delegate: %@", delegate);
});

}