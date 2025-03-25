const std = @import("std");
const dd = @import("internal/dd.zig");
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;
const assert = std.debug.assert;

const ENABLE_LINE_INPUT: u32 = 0x2;
const ENABLE_ECHO_INPUT: u32 = 0x4;
const ENABLE_PROCESSED_INPUT: u32 = 0x1;
const ENABLE_WINDOW_INPUT: u32 = 0x8;

const OldState = union {
    win: struct {
        std_out: shit_os.DWORD,
        std_in: shit_os.DWORD,
    },
    posix: struct {
        std_in: posix.termios,
    },
};

const Screen = struct {
    mutex: std.Thread.Mutex = .{},
    condition: std.Thread.Condition = .{},
    exit: bool = false,
    render_q: std.ArrayList(u8),

    const Self = @This();

    fn new(alloc: std.mem.Allocator) !Self {
        return Self{ .render_q = std.ArrayList(u8).init(alloc) };
    }

    fn add_to_render_q(self: *Self, line: []u8) !void {
        self.mutex.lock();
        defer self.mutex.unlock();
        try self.render_q.appendSlice(line);
        self.condition.signal();
    }
};

pub fn read_terminal(self: *Screen, std_in: std.fs.File) !void {
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

pub fn receive(screen: *Screen) !void {
    var fuck_text = "\x1b[38;2;255;0;0mFUCK".*;
    std.Thread.sleep(1_000_000_000);
    try screen.add_to_render_q(&fuck_text);
    std.Thread.sleep(1_000_000_000);
    try screen.add_to_render_q(&fuck_text);
}

pub fn main() !void {
    const std_out = std.io.getStdOut();

    const std_in = std.io.getStdIn();
    const stdout = std_out.writer();

    var terminal_dimensions: dd.TerminalDimensions = undefined;

    try dd.get_terminal_dimensions(
        std_out,
        &terminal_dimensions,
    );
    try stdout.print("{}\n", .{terminal_dimensions});

    // TODO: why does this not try?
    var termos = dd.get_termos_with_tea();
    try dd.start_raw_mode(std_in, std_out, &termos);
    defer dd.restore_terminal(std_in, std_out, termos);

    var gpa = std.heap.DebugAllocator(.{}).init;

    var screen = try Screen.new(gpa.allocator());
    defer screen.render_q.deinit();

    const render_thread = try std.Thread.spawn(.{}, read_terminal, .{ &screen, std_in });
    defer render_thread.join();

    const receive_thread = try std.Thread.spawn(.{}, receive, .{&screen});
    defer receive_thread.join();

    screen.mutex.lock();
    defer screen.mutex.unlock();
    while (screen.render_q.items.len == 0) {
        screen.condition.wait(&screen.mutex);
        if (screen.exit) break;
        try stdout.print("\x1b[48;2;25;60;80m{s}\x1b[0m\n", .{screen.render_q.items});
        screen.render_q.clearAndFree();
    }
}
