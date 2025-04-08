const std = @import("std");
const fs = std.fs;
const os = std.os;
const shit_os = std.os.windows;
const posix = std.posix;
const assert = std.debug.assert;
const os_tag = @import("builtin").os.tag;
const common = @import("layers/common.zig");

const RenderQ = @import("RenderQueue.zig");
const ui_common = @import("layers/common.zig");
const Login = @import("layers/Login.zig");
const MainLayer = @import("layers/MainLayer.zig");

// row;col;text
const STATUS_LINE_PATTERN = "\x1b[{};{}H\x1b[2K\x1b[48;2;251;206;44m\x1b[38;2;0;0;0m{s}\x1b[0m";

terminal_dimensions: common.Dimensions = undefined,
exit: bool = false,
render_q: RenderQ,

active_mode: common.MODE = .NORMAL,
active_layer: common.LAYERS = .MAIN,

// layers
main_layer: MainLayer = undefined,

status_line: []u8 = undefined,
status_line_content_len: usize = 0,

alloc: std.mem.Allocator,

const Self = @This();

// const RenderFlags = struct {
//     status_line: bool = true,
//     partial: bool = false,
//     login: bool = false,
//     // NOTE: It's only for testing
//     // later we should add login as first and check if user is logged in
//     main: bool = true,
// };

// pub var render_flags: RenderFlags = .{};
pub fn create(std_out: std.fs.File, alloc: std.mem.Allocator) !Self {
    return Self{
        .alloc = alloc,
        .render_q = .{
            .queue = std.ArrayList(u8).init(alloc),
        },
        .terminal_dimensions = try get_terminal_dimensions(std_out),
    };
}

pub fn create_layers(self: *Self) !void {
    self.main_layer = try MainLayer.create(self.alloc, self.terminal_dimensions, &self.render_q);
}

fn get_terminal_dimensions(std_out: fs.File) !common.Dimensions {
    var dimensions: common.Dimensions = undefined;
    switch (os_tag) {
        .windows => {
            var console_info: shit_os.CONSOLE_SCREEN_BUFFER_INFO = undefined;
            _ = shit_os.kernel32.GetConsoleScreenBufferInfo(std_out.handle, &console_info);
            dimensions = .{ .width = console_info.dwSize.X, .height = console_info.dwSize.Y };
            // self.terminal_dimensions.width = ;
            // self.terminal_dimensions.height = console_info.dwSize.Y;
        },
        .linux, .macos => {
            var win_size: posix.winsize = undefined;
            const res = posix.system.ioctl(std_out.handle, os.linux.T.IOCGWINSZ, @intFromPtr(&win_size));
            if (res != 0) {
                return error.ioctl_return_error_during_getting_linux_dimentions;
            }
            dimensions = .{ .width = @intCast(win_size.col), .height = @intCast(win_size.row) };
        },
        else => return error.UNSUPPORTED_OS,
    }
    assert(dimensions.width != 0 and dimensions.height != 0); // how?
    assert(dimensions.width != std.math.maxInt(i10)); // waytoodank 511 columns is a lot
    assert(dimensions.height < dimensions.width); // we do not support vertical
    return dimensions;
}

/// initialize screen before we start reading from terminal
/// here we check the connection, pick layer to render (login/main)
/// clear screen and finally render first frame
pub fn init_first_frame(self: *Self, stdout: fs.File.Writer) !void {
    // hide cursor
    const hidden_cursor_slice = try std.fmt.allocPrint(self.alloc, "{s}", .{common.HIDDEN_CURSOR});
    try self.render_q.add_to_render_q(hidden_cursor_slice);
    self.status_line = try self.alloc.alloc(u8, @intCast(self.terminal_dimensions.width));
    @memset(self.status_line, ' ');
    try self.render_q.queue.appendSlice(common.CLEAR_SCREEN);
    try self.render_q.queue.appendSlice(common.theme.active_background_color);

    // layer logic here
    try self.main_layer.render_first_frame();
    try self.change_mode(.NORMAL);
    try stdout.print("{s}", .{self.render_q.queue.items});
    self.render_q.queue.clearAndFree();
}

fn render_status_line(self: Self) ![]u8 {
    return std.fmt.allocPrint(self.alloc, STATUS_LINE_PATTERN, .{ self.terminal_dimensions.height, 1, self.status_line });
}

pub fn change_mode(self: *Self, new_mode: common.MODE) !void {
    self.active_mode = new_mode;

    const mode_name = common.MODE_MAP[@intFromEnum(self.active_mode)];
    @memset(self.status_line[0..self.status_line_content_len], ' ');
    @memcpy(self.status_line[0..mode_name.len], mode_name);
    self.status_line_content_len = mode_name.len;

    const status_line_content = try self.render_status_line();
    defer self.alloc.free(status_line_content);
    try self.render_q.add_to_render_q(status_line_content);
}

fn append_to_command(self: *Self, char: u8) !void {
    self.status_line[self.status_line_content_len] = char;
    self.status_line_content_len += 1;
    const status_line_content = try self.render_status_line();
    defer self.alloc.free(status_line_content);
    try self.render_q.add_to_render_q(status_line_content);
}

pub fn read_terminal(self: *Self, std_in: fs.File) !void {
    while (!self.exit) {
        var buf = [1]u8{0};
        const bytes_read = try std_in.read(&buf);
        assert(bytes_read != 0);

        const curr_char = buf[0];

        // FIXME: remove later
        if (curr_char == 3) {
            const result = try std.fmt.allocPrint(self.alloc, "{s}", .{common.VISIBLE_CURSOR});
            try self.render_q.add_to_render_q(result);
            self.exit = true;
            self.render_q.condition.signal();
            return;
        }

        const prev_mode = self.active_mode;
        switch (self.active_mode) {
            .COMMAND => {
                switch (curr_char) {
                    '\r' => {
                        try self.handle_command();
                        try self.change_mode(.NORMAL);
                    },
                    else => {
                        try self.append_to_command(curr_char);
                    },
                }
            },
            // FIXME: MOVE THAT TO LAYER HANDLERS
            .NORMAL => {
                switch (curr_char) {
                    ':' => {
                        try self.change_mode(.COMMAND);
                    },
                    else => {
                        switch (self.active_layer) {
                            .MAIN => {
                                try self.main_layer.handle_current_state(self.active_mode, curr_char);
                            },
                            else => {},
                        }
                    },
                }
            },
            else => {},
        }
        if (prev_mode != self.active_mode) {
            // TODO deal with the mode change
        }
    }
}

fn handle_command(self: *Self) !void {
    const command = common.KNOWN_COMMANDS.get(self.status_line[0..self.status_line_content_len]);
    std.debug.print("COMMAND: {any} vs ITEMS: {s} vs AVAILABLE: {any}\n", .{ command, self.status_line, common.KNOWN_COMMANDS });
    if (command) |real_command| switch (real_command) {
        .QUIT => {
            const result = try std.fmt.allocPrint(self.alloc, "{s}", .{common.VISIBLE_CURSOR});
            try self.render_q.add_to_render_q(result);
            self.exit = true;
        },
        .NEW_CONVERSATION => {
            std.debug.print("wtf\n", .{});
        },
    };
}
