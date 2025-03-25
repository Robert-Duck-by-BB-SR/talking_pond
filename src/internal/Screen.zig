const std = @import("std");
const fs = std.fs;
const os = std.os;
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;

const assert = std.debug.assert;

const TerminalDimensions = struct { width: i16, height: i16 };
mutex: std.Thread.Mutex = .{},
condition: std.Thread.Condition = .{},
exit: bool = false,
render_q: std.ArrayList(u8),
terminal_dimensions: TerminalDimensions = undefined,

const Self = @This();

pub fn new(alloc: std.mem.Allocator) !Self {
    return Self{ .render_q = std.ArrayList(u8).init(alloc) };
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

pub fn read_terminal(self: *Self, std_in: std.fs.File) !void {
    loop: while (true) {
        var buf = [1]u8{0};
        _ = try std_in.read(&buf);
        switch (buf[0]) {
            3 => {
                self.mutex.lock();
                defer self.mutex.unlock();
                self.exit = true;
                self.condition.signal();
                break :loop;
            },
            else => try self.add_to_render_q(&buf),
        }
    }
}
