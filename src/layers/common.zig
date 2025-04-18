const std = @import("std");

pub const LAYER_TYPE = enum {
    MODAL,
    WINDOW,
};

pub const CONTENT_TYPE = enum {
    CONTENT,
    STATUS,
    CURSOR,
};

pub const LAYERS = enum {
    LOGIN,
    MAIN,
};

pub const MODE = enum {
    NORMAL,
    INSERT,
    COMMAND,

    // DON'T TOUCH IT IS USED IN MODE MAP AS THE LENGTH AND THEREFORE SHOULD BE THE LAST ITEM
    MODES_COUNT,
};

pub const ComponentType = enum { PONDS_SIDEBAR, QUACKS_CHAT, INPUT_FIELD };

pub const MODE_MAP = [@intFromEnum(MODE.MODES_COUNT)][]const u8{
    "NORMAL",
    "INSERT",
    ":",
};

const COMMANDS = enum {
    QUIT,
    NEW_CONVERSATION,
};

pub const KNOWN_COMMANDS = std.StaticStringMap(COMMANDS).initComptime(.{
    .{ ":q", .QUIT },
    .{ ":new", .NEW_CONVERSATION },
});

pub const INVERT_STYLES = "\x1b[7m";
pub const RED_COLOR = "\x1b[31m;";
pub const RESET_STYLES = "\x1b[0m";

pub const CLEAR_SCREEN = "\x1b[2J";
pub const MOVE_CURSOR_TO_THE_BENINGING = "\x1b[H";
pub const MOVE_CURSOR_TO_POSITION = "\x1b[{};{}H";
pub const CLEAR_ROW = "\x1b[2K";
pub const HIDDEN_CURSOR = "\x1b[?25l";
pub const VISIBLE_CURSOR = "\x1b[?25h";

pub const Dimensions = struct { width: i16, height: i16 };
pub const Position = struct { row: i16, col: i16 };

pub const PONDS_SIDEBAR_SIZE: i16 = 35;

pub const NormalBorder = struct {
    pub const HORIZONTAL = "─";
    pub const VERTICAL = "│";
    pub const TOP_LEFT = "┌";
    pub const TOP_RIGHT = "┐";
    pub const BOTTOM_LEFT = "└";
    pub const BOTTOM_RIGHT = "┘";
};

pub const NOTIFICATION_ICON = "\u{25FC}";

pub const theme = struct {
    pub const FONT_COLOR = "\x1b[38;2;192;192;192m";
    pub const ACTIVE_FONT_COLOR = "\x1b[38;2;112;255;112m";
    pub const BACKGROUND_COLOR = "\x1b[48;2;25;25;25m";
    pub const ACTIVE_BACKGROUND_COLOR = "\x1b[48;2;155;100;0m";
    pub const BORDER = NormalBorder;
};

pub const NOTIFICATION_ICON_PATTERN = theme.ACTIVE_FONT_COLOR ++ NOTIFICATION_ICON ++ RESET_STYLES;


pub const ACTIVE_ITEM = theme.FONT_COLOR ++ theme.ACTIVE_BACKGROUND_COLOR;
pub const INACTIVE_ITEM = theme.FONT_COLOR ++ theme.BACKGROUND_COLOR;
pub const ACTIVE_BORDER = theme.ACTIVE_FONT_COLOR ++ theme.BACKGROUND_COLOR;

pub fn render_border(alloc: std.mem.Allocator, is_active: bool, border: []u8) ![]u8 {
    var ponds: std.ArrayList(u8) = .init(alloc);
    try ponds.writer().print("{s}{s}", .{
        if (is_active) ACTIVE_BORDER else INACTIVE_ITEM,
        border,
    });
    return ponds.toOwnedSlice();
}
