
// +build darwin,!ios

#include <AppKit/AppKit.h>
#include <QuartzCore/QuartzCore.h>
#include "_cgo_export.h"

const int cellWidth = 200;
const int cellHeight = 60;
const int cellSpace = 10;
const int paddingHorizontal = 10;
const int iconSize = 20;

// copy from https://github.com/cool8jay/VerticalCenterField/blob/master/VerticalCenterField/ViewController.m
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

@property (nonatomic, strong)   NSView *ripple;

@end


@implementation CodeView


- (id)init {
    self = [super init];

    NSView *ripple = [[NSView alloc] init];
    ripple.wantsLayer = YES;
    ripple.layer.backgroundColor = [NSColor grayColor].CGColor;
    ripple.alphaValue = 0.0f;
    self.ripple = ripple;
    [self addSubview:ripple];

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
    code_on_click((__bridge CFTypeRef)self.secret);

    CGFloat duration = 0.5;
    NSPoint p = event.locationInWindow;
    CGFloat size = self.frame.size.width * 3;
    CGFloat x = p.x - self.frame.origin.x;
    CGFloat y = p.y - self.frame.origin.y;
    x = -size*0.1/2 + x;
    y = -size*0.1/2 + y;
    CGFloat tx = -size/2+x;
    CGFloat ty = -size/2+y;
    CAMediaTimingFunctionName fn = kCAMediaTimingFunctionLinear;

    self.ripple.layer.cornerRadius = size * 0.5f;
    self.ripple.frame = CGRectMake(tx, ty, size, size);

    CABasicAnimation *alphaAnimation =  [CABasicAnimation animationWithKeyPath:@"opacity"];
    alphaAnimation.fromValue = [NSNumber numberWithFloat:0.8];
    alphaAnimation.toValue = [NSNumber numberWithFloat:0.0];
    alphaAnimation.duration = duration;
    alphaAnimation.timingFunction = [CAMediaTimingFunction functionWithName:fn];


    CABasicAnimation *scaleAnimation = [CABasicAnimation animationWithKeyPath:@"transform.scale"];
    scaleAnimation.fromValue = [NSNumber numberWithFloat:0.1];
    scaleAnimation.toValue = [NSNumber numberWithFloat:1];
    scaleAnimation.duration = duration;
    scaleAnimation.timingFunction = [CAMediaTimingFunction functionWithName:fn];

    CABasicAnimation *posYAnimation = [CABasicAnimation animationWithKeyPath:@"position.y"];
    posYAnimation.fromValue = [NSNumber numberWithFloat:y];
    posYAnimation.toValue = [NSNumber numberWithFloat:ty];
    posYAnimation.duration = duration;
    posYAnimation.timingFunction = [CAMediaTimingFunction functionWithName:fn];

    CABasicAnimation *posXAnimation = [CABasicAnimation animationWithKeyPath:@"position.x"];
    posXAnimation.fromValue = [NSNumber numberWithFloat:x];
    posXAnimation.toValue = [NSNumber numberWithFloat:tx];
    posXAnimation.duration = duration;
    posXAnimation.timingFunction = [CAMediaTimingFunction functionWithName:fn];

    [self.ripple.layer addAnimation:posXAnimation forKey:@"posXAnimation"];
    [self.ripple.layer addAnimation:posYAnimation forKey:@"posYAnimation"];
    [self.ripple.layer addAnimation:scaleAnimation forKey:@"scaleAnimation"];
    [self.ripple.layer addAnimation:alphaAnimation forKey:@"alphaAnimation"];
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


@interface PopoverController : NSViewController

@property (nonatomic, strong) NSArray<NSDictionary*>* data;
@property (nonatomic, strong) CodesView *codes;
@property (nonatomic, weak) NSPopover *popover;

@property (nonatomic, strong) NSButton *settingsBtn;
@property (nonatomic, strong) NSButton *dashboardBtn;
@property (nonatomic, strong) NSButton *quitBtn;

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
    int height = (cellHeight+cellSpace) * len + iconSize + cellSpace;

    self.codes = [[CodesView alloc] initWithFrame:CGRectMake(0, 0, cellWidth, height)];
    NSMutableArray<CodeView*> *codesView = [NSMutableArray array];
	for (NSDictionary *dic in self.data) {
        CodeView *code = [[CodeView alloc] init];
        code.secret = [dic objectForKey:@"secret"];
        code.name.stringValue = [dic objectForKey: @"name"];
        [codesView addObject:code];
	}
    self.codes.codes = codesView;
    [self.view addSubview:self.codes];

    [self.codes refreshCodes];
    [self refreshButton];
}

- (void) refreshButton {
    [self ensureButton];

    int len = [self.data count];
    int oy = (cellHeight+cellSpace) * len;

    self.settingsBtn.frame = CGRectMake(paddingHorizontal, oy, iconSize, iconSize);
    self.dashboardBtn.frame = CGRectMake(paddingHorizontal + iconSize + paddingHorizontal, oy, iconSize, iconSize);
    self.quitBtn.frame = CGRectMake(cellWidth-iconSize-paddingHorizontal, oy, iconSize, iconSize);
}

- (void) ensureButton {
    if (self.settingsBtn ==  nil) {
        NSImage *settingsIcon = CFBridgingRelease(export_settings_icon());
        settingsIcon.size = CGSizeMake(iconSize, iconSize);

        self.settingsBtn = [[NSButton alloc] init];
        self.settingsBtn.image = settingsIcon;
        self.settingsBtn.wantsLayer = YES;
        self.settingsBtn.bordered = NO;
        self.settingsBtn.target = self;
        self.settingsBtn.action = @selector(buttonClick:);
        self.settingsBtn.layer.backgroundColor = [[NSColor clearColor] CGColor];
        [self.view addSubview: self.settingsBtn];

        NSImage *dashboardIcon = CFBridgingRelease(export_dashboard_icon());
        dashboardIcon.size = CGSizeMake(iconSize, iconSize);

        self.dashboardBtn = [[NSButton alloc] init];
        self.dashboardBtn.image = dashboardIcon;
        self.dashboardBtn.wantsLayer = YES;
        self.dashboardBtn.bordered = NO;
        self.dashboardBtn.target = self;
        self.dashboardBtn.action = @selector(buttonClick:);
        self.dashboardBtn.layer.backgroundColor = [[NSColor clearColor] CGColor];
        [self.view addSubview: self.dashboardBtn];

        NSImage *quitIcon = CFBridgingRelease(export_quit_icon());
        quitIcon.size = CGSizeMake(iconSize, iconSize);

        self.quitBtn = [[NSButton alloc] init];
        self.quitBtn.image = quitIcon;
        self.quitBtn.wantsLayer = YES;
        self.quitBtn.bordered = NO;
        self.quitBtn.target = self;
        self.quitBtn.action = @selector(buttonClick:);
        self.quitBtn.layer.backgroundColor = [[NSColor clearColor] CGColor];
        [self.view addSubview: self.quitBtn];
    }
}

- (void) buttonClick:(NSButton *) button {
    // NSLog(@"settingsBtnClicked");
    if (self.settingsBtn == button) {
        tray_button_on_click(1);
    } else if (self.dashboardBtn == button) {
        tray_button_on_click(2);
    } else {
        tray_button_on_click(3);
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
    int height = (cellHeight+cellSpace) * len + iconSize + cellSpace;

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
	id delegate = [[NSApplication sharedApplication] delegate];
	NSStatusItem* statusItem = [delegate performSelector:@selector(statusItem)];
	if (statusItem == nil) {
	    NSStatusItem* statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
        NSImage *image = CFBridgingRelease(export_icon_data());
        [image setSize:NSMakeSize(16, 16)];
        image.template = true;
        statusItem.button.image = image;
    	id obj = [PopoverManager class];
        SEL mySelector = @selector(show:);
    	statusItem.button.target = obj;
        statusItem.button.action = mySelector;
    	[delegate performSelector:@selector(setStatusItem:) withObject:statusItem];
	}
    statusItem.visible = YES;
});

}

void dismiss_tray(void) {

dispatch_async(dispatch_get_main_queue(), ^{
	id delegate = [[NSApplication sharedApplication] delegate];
	NSStatusItem* statusItem = [delegate performSelector:@selector(statusItem)];
    if (statusItem == nil) {
        NSStatusItem* statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
        NSImage *image = CFBridgingRelease(export_icon_data());
        [image setSize:NSMakeSize(16, 16)];
        image.template = true;
        statusItem.button.image = image;
        id obj = [PopoverManager class];
        SEL mySelector = @selector(show:);
       	statusItem.button.target = obj;
        statusItem.button.action = mySelector;
       	[delegate performSelector:@selector(setStatusItem:) withObject:statusItem];
    }
    statusItem.visible = NO;
});

}
