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
};

pub const BG_KEY = "\x33[48;2;";
pub const INVERT_STYLES = "\x33[7m";
pub const RED_COLOR = "\x33[31m;";
pub const RESET_STYLES = "\x33[0m";

pub const CLEAR_SCREEN = "\x33[2J";
pub const MOVE_CURSOR_TO_THE_BENINGING = "\x33[H";
pub const MOVE_CURSOR_TO_POSITION = "\x33[%d;%dH";
pub const CLEAR_ROW = "\x33[2K";
pub const HIDDEN_CURSOR = "\x33[?25l";
pub const VISIBLE_CURSOR = "\x33[?25h";
