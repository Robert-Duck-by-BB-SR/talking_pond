const std = @import("std");
const fs = std.fs;
const os = std.os;
const shit_os = std.os.windows;
const posix = std.posix;
const assert = std.debug.assert;
const os_tag = @import("builtin").os.tag;
const ui_common = @import("../internal/layers/common.zig");

const TerminalDimensions = struct { width: i16, height: i16 };
mutex: std.Thread.Mutex = .{},
condition: std.Thread.Condition = .{},
exit: bool = false,
render_q: std.ArrayList(u8),
terminal_dimensions: TerminalDimensions = undefined,

active_mode: ui_common.MODE = .NORMAL,
active_layer: ui_common.LAYERS = .LOGIN,
status_line: StatusLine = .{},

const Self = @This();

const StatusLine = struct {
    state: std.ArrayList(u8) = undefined,

    const InnerSelf = @This();

    fn display(self: *InnerSelf) *std.ArrayList(u8) {
        return &self.state;
    }
};

const COMMANDS = enum {
    QUIT,
    NEW_CONVERSATION,
};

var known_commands: std.StringHashMap(COMMANDS) = undefined;

pub fn new(alloc: std.mem.Allocator) !Self {
    // initialize known commands
    known_commands = std.StringHashMap(COMMANDS).init(alloc);

    try known_commands.put(":q", .QUIT);
    try known_commands.put(":new", .NEW_CONVERSATION);

    return Self{
        .render_q = std.ArrayList(u8).init(alloc),
        .status_line = .{
            .state = std.ArrayList(u8).init(alloc),
        },
    };
}

pub fn destroy(self: *Self) void {
    self.render_q.deinit();
    self.status_line.state.deinit();
}

pub fn add_to_render_q(self: *Self, line: *std.ArrayList(u8)) !void {
    self.mutex.lock();
    defer self.mutex.unlock();
    const slice = try line.toOwnedSlice();
    try self.render_q.appendSlice(slice);
    self.condition.signal();
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

pub fn read_terminal(self: *Self, std_in: std.fs.File) !void {
    while (true) {
        var buf = [1]u8{0};
        _ = try std_in.read(&buf);

        if (buf[0] == 3) {
            self.exit = true;
            return;
        }

        const prev_mode = self.active_mode;
        switch (self.active_mode) {
            .COMMAND => {
                switch (buf[0]) {
                    '\r' => {
                        self.handle_command();
                        self.active_mode = .NORMAL;
                        self.status_line.state.clearAndFree();
                        try self.status_line.state.appendSlice("NORMAL");
                    },
                    else => {
                        try self.status_line.state.append(buf[0]);
                    },
                }
                try self.add_to_render_q(self.status_line.display());
            },
            else => {},
        }
        if (prev_mode != self.active_mode) {
            // TODO deal with the mode change
        }
    }
}

fn handle_command(self: *Self) void {
    const command = self.status_line.state.toOwnedSlice() catch unreachable;
    if (known_commands.get(command)) |real_command| switch (real_command) {
        .QUIT => self.exit = true,
        .NEW_CONVERSATION => {},
    };
}


const testing = std.testing;
test "test quitting" {
    const cwd = fs.cwd();
    const file = cwd.createFile("test_quitting_std_in", .{});
}
