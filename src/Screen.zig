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

const TerminalDimensions = struct { width: i16, height: i16 };

/// row;col;text
const STATUS_LINE = "\x1b[{};{}H\x1b[2K\x1b[48;2;251;206;44m\x1b[38;2;192;192;192m{s}\x1b[0m";

mutex: std.Thread.Mutex = .{},
condition: std.Thread.Condition = .{},
staying_alive: bool = true,
render_q: std.ArrayList(u8),
terminal_dimensions: TerminalDimensions = undefined,

active_mode: common.MODE = .NORMAL,
active_layer: common.LAYERS = .LOGIN,
status_line: StatusLine = .{},
alloc: std.mem.Allocator,

const Self = @This();

const COMMANDS = enum {
    QUIT,
    NEW_CONVERSATION,
};

// TODO: make it static string map
var known_commands: std.StringHashMap(common.COMMANDS) = undefined;

const RenderFlags = struct {
    status_line: bool = true,
    partial: bool = false,
    login: bool = false,
    main: bool = false,
};

pub var ready_to_render: RenderFlags = .{};
pub fn new(alloc: std.mem.Allocator) !Self {
    // initialize known commands
    known_commands = std.StringHashMap(common.COMMANDS).init(alloc);

    try known_commands.put(":q", .QUIT);
    try known_commands.put(":new", .NEW_CONVERSATION);

    return Self{
        .alloc = alloc,
        .render_q = std.ArrayList(u8).init(alloc),
        .status_line = std.ArrayList(u8).init(alloc),
    };
}

pub fn destroy(self: *Self) void {
    self.render_q.deinit();
    self.status_line.deinit();
}

pub fn add_to_render_q(self: *Self, line: []u8) !void {
    self.mutex.lock();
    defer self.mutex.unlock();
    try self.render_q.appendSlice(line);
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
                        self.active_mode = .NORMAL;
                        self.status_line.clearAndFree();
                        try self.status_line.appendSlice("NORMAL");
                    },
                    else => {
                        try self.status_line.append(curr_char);
                    },
                }
                try self.add_to_render_q(self.status_line.items);
            },
            .NORMAL => {
                switch (curr_char) {
                    ':' => {
                        self.active_mode = .COMMAND;
                        self.status_line.clearAndFree();
                        try self.status_line.append(':');
                    },
                    else => {},
                }
            },
            .NORMAL => {
                switch (buf[0]) {
                    ':' => {
                        self.active_mode = .COMMAND;
                        self.status_line.state.clearAndFree();
                        try self.status_line.state.append(':');
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
    const items = self.status_line.items;
    const command = known_commands.get(items);
    std.debug.print("COMMAND: {any} vs ITEMS: {s} vs AVAILABLE: {any}\n", .{ command, items, known_commands });
    if (command) |real_command| switch (real_command) {
        .QUIT => {
            self.staying_alive = true;
        },
        .NEW_CONVERSATION => {},
    };
}

pub fn render(self: *Self, stdout: fs.File.Writer) !void {
    if (!ready_to_render.partial) {
        ready_to_render.partial = true;
        try stdout.print("\x1b[2J", .{});
    }

    try stdout.print("{s}", .{self.render_q.items});

    if (ready_to_render.status_line) {
        defer ready_to_render.status_line = false;
        const status_line = try self.render_status_line();
        try stdout.print("{s}", .{status_line});
    }
    self.render_q.clearAndFree();
}

fn render_status_line(self: Self) ![]u8 {
    return std.fmt.allocPrint(self.alloc, STATUS_LINE, .{ self.terminal_dimensions.height, 0, self.status_line.items });
}

pub fn render_available() bool {
    return ready_to_render.login or ready_to_render.main or ready_to_render.status_line;
}
