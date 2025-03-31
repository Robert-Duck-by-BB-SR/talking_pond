const std = @import("std");
const fs = std.fs;
const os = std.os;
const shit_os = std.os.windows;
const posix = std.posix;
const assert = std.debug.assert;
const os_tag = @import("builtin").os.tag;
const common = @import("layers/common.zig");

const ui_common = @import("../internal/layers/common.zig");
const Login = @import("../internal/layers/Login.zig");

// row;col;text
const STATUS_LINE_PATTERN = "\x1b[{};{}H\x1b[2K\x1b[48;2;251;206;44m\x1b[38;2;0;0;0m{s}\x1b[0m";

mutex: std.Thread.Mutex = .{},
condition: std.Thread.Condition = .{},
exit: bool = false,
render_q: std.ArrayList(u8),
terminal_dimensions: common.Dimensions = undefined,

active_mode: common.MODE = .NORMAL,
active_layer: common.LAYERS = .LOGIN,
status_line: []u8 = undefined,
status_line_content_len: usize = 0,
alloc: std.mem.Allocator,

const Self = @This();

const RenderFlags = struct {
    status_line: bool = true,
    partial: bool = false,
    login: bool = false,
    main: bool = false,
};

pub var ready_to_render: RenderFlags = .{};
pub fn new(alloc: std.mem.Allocator) !Self {
    return Self{
        .alloc = alloc,
        .render_q = std.ArrayList(u8).init(alloc),
    };
}

pub fn destroy(self: *Self) void {
    self.render_q.deinit();
    self.alloc.free(self.status_line);
}

pub fn init_first_frame(self: *Self, stdout: fs.File.Writer) !void {
    // initialize screen before we start reading from terminal
    // here we check the connection, pick layer to render (login/main)
    // clear screen and finally render first frame
    self.status_line = try self.alloc.alloc(u8, @intCast(self.terminal_dimensions.width));
    @memset(self.status_line, ' ');
    try self.render_q.appendSlice("\x1b[2J");
    try self.change_mode(.NORMAL);
    try stdout.print("{s}", .{self.render_q.items});
    self.render_q.clearAndFree();
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
    try self.add_to_render_q(status_line_content);
}

pub fn add_to_render_q(self: *Self, line: []u8) !void {
    self.mutex.lock();
    defer self.mutex.unlock();
    try self.render_q.appendSlice(line);
    self.condition.signal();
}

fn append_to_command(self: *Self, char: u8) !void {
    self.status_line[self.status_line_content_len] = char;
    self.status_line_content_len += 1;
    const status_line_content = try self.render_status_line();
    defer self.alloc.free(status_line_content);
    try self.add_to_render_q(status_line_content);
}

pub fn get_terminal_dimensions(self: *Self, std_out: fs.File) !void {
    switch (os_tag) {
        .windows => {
            var console_info: shit_os.CONSOLE_SCREEN_BUFFER_INFO = undefined;
            _ = shit_os.kernel32.GetConsoleScreenBufferInfo(std_out.handle, &console_info);
            self.terminal_dimensions.width = console_info.dwSize.X;
            self.terminal_dimensions.height = console_info.dwSize.Y;
        },
        .linux, .macos => {
            var win_size: posix.winsize = undefined;

            const res = posix.system.ioctl(std_out.handle, os.linux.T.IOCGWINSZ, @intFromPtr(&win_size));
            if (res != 0) {
                return error.ioctl_return_error_during_getting_linux_dimentions;
            }
            self.terminal_dimensions.width = @intCast(win_size.col);
            self.terminal_dimensions.height = @intCast(win_size.row);
        },
        else => return error.UNSUPPORTED_OS,
    }
    assert(self.terminal_dimensions.width != 0 and self.terminal_dimensions.height != 0); // how?
    assert(self.terminal_dimensions.width != std.math.maxInt(i10)); // waytoodank 511 columns is a lot
    assert(self.terminal_dimensions.height < self.terminal_dimensions.width); // we do not support vertical

}

pub fn read_terminal(self: *Self, std_in: fs.File) !void {
    while (!self.exit) {
        var buf = [1]u8{0};
        const bytes_read = try std_in.read(&buf);
        assert(bytes_read != 0);

        const curr_char = buf[0];

        // FIXME: remove later
        if (curr_char == 3) {
            self.exit = true;
            self.condition.signal();
            return;
        }

        const prev_mode = self.active_mode;
        switch (self.active_mode) {
            .COMMAND => {
                switch (curr_char) {
                    '\r' => {
                        self.handle_command();
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
                    else => {},
                }
            },
            else => {},
        }
        if (prev_mode != self.active_mode) {
            // TODO deal with the mode change
        }
    }
}

fn handle_command(self: *Self) void {
    const command = common.KNOWN_COMMANDS.get(self.status_line[0..self.status_line_content_len]);
    std.debug.print("COMMAND: {any} vs ITEMS: {s} vs AVAILABLE: {any}\n", .{ command, self.status_line, common.KNOWN_COMMANDS });
    if (command) |real_command| switch (real_command) {
        .QUIT => {
            self.exit = true;
        },
        .NEW_CONVERSATION => {
            std.debug.print("wtf\n", .{});
        },
    };
}
