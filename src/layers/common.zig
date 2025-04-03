const std = @import("std");

pub const LAYER_TYPE = enum {
    MODAL,
    WINDOW,
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

pub const BG_KEY = "\x1b[48;2;";
pub const INVERT_STYLES = "\x1b[7m";
pub const RED_COLOR = "\x1b[31m;";
pub const RESET_STYLES = "\x1b[0m";

pub const CLEAR_SCREEN = "\x1b[2J";
pub const MOVE_CURSOR_TO_THE_BENINGING = "\x1b[H";
pub const MOVE_CURSOR_TO_POSITION = "\x1b[%d;%dH";
pub const CLEAR_ROW = "\x1b[2K";
pub const HIDDEN_CURSOR = "\x1b[?25l";
pub const VISIBLE_CURSOR = "\x1b[?25h";

pub const Dimensions = struct { width: i16, height: i16 };
pub const Position = struct { row: i16, col: i16 };

pub const PONDS_SIDEBAR_SIZE: i16 = 70;

pub const NormalBorder = struct {
    pub const HORIZONTAL = "─";
    pub const VERTICAL = "│";

    pub const TOP_LEFT = "┌";
    pub const TOP_RIGHT = "┐";
    pub const BOTTOM_LEFT = "└";
    pub const BOTTOM_RIGHT = "┘";
};
