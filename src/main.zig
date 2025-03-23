const std = @import("std");
const dd = @import("internal/dd.zig");
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;

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

fn start_raw_mode(std_in: std.fs.File, std_out: std.fs.File, termos: *OldState) !void {
    switch (os_tag) {
        .windows => {
            var old_stdin: shit_os.DWORD = undefined;
            _ = shit_os.kernel32.GetConsoleMode(std_in.handle, &old_stdin);
            var raw_mode = old_stdin & ~(ENABLE_LINE_INPUT |
                ENABLE_ECHO_INPUT |
                ENABLE_PROCESSED_INPUT);

            raw_mode |= ENABLE_WINDOW_INPUT;

            _ = shit_os.kernel32.SetConsoleMode(std_in.handle, raw_mode);

            _ = shit_os.kernel32.GetConsoleMode(std_out.handle, &termos.win.std_out);
            _ = shit_os.kernel32.SetConsoleMode(std_out.handle, termos.win.std_out | shit_os.ENABLE_VIRTUAL_TERMINAL_PROCESSING);
        },
        .linux => {
            termos.posix.std_in = try posix.tcgetattr(std_in.handle);
            var raw = termos.posix.std_in;
            raw.lflag.ECHO = false;
            raw.lflag.ICANON = false;
            raw.lflag.ISIG = false;
            raw.lflag.IEXTEN = false;
            raw.iflag.ICRNL = false;
            raw.iflag.IXON = false;
            raw.oflag.OPOST = false;
            try posix.tcsetattr(
                std_in.handle,
                .FLUSH,
                raw,
            );
        },
        else => return error.UNSUPPORTED_OS,
    }
}

fn restore_terminal(std_in: std.fs.File, std_out: std.fs.File, termos: OldState) void {
    switch (os_tag) {
        .windows => {
            _ = shit_os.kernel32.SetConsoleMode(std_in.handle, termos.win.std_in);
            _ = shit_os.kernel32.SetConsoleMode(std_out.handle, termos.win.std_out);
        },
        .linux => {
            posix.tcsetattr(
                std_in.handle,
                .FLUSH,
                termos.posix.std_in,
            ) catch unreachable;
        },
        else => {},
    }
}

pub const TerminalDimensions = struct { width: i16, height: i16 };

fn get_terminal_dimensions(std_out: std.fs.File, terminal_dimensions: *TerminalDimensions) !void {
    switch (os_tag) {
        .windows => {
            var console_info: shit_os.CONSOLE_SCREEN_BUFFER_INFO = undefined;
            _ = shit_os.kernel32.GetConsoleScreenBufferInfo(std_out.handle, &console_info);
            terminal_dimensions.width = console_info.dwSize.X;
            terminal_dimensions.height = console_info.dwSize.Y;
        },
        .linux, .macos => {
            var win_size: std.posix.winsize = undefined;

            const res = posix.system.ioctl(std_out.handle, std.os.linux.T.IOCGWINSZ, @intFromPtr(&win_size));
            if (res != 0) {
                return error.ioctl_return_error_during_getting_linux_dimentions;
            }
            terminal_dimensions.width = @intCast(win_size.col);
            terminal_dimensions.height = @intCast(win_size.row);
        },
        else => return error.UNSUPPORTED_OS,
    }
}
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

    var terminal_dimensions: TerminalDimensions = undefined;

    try get_terminal_dimensions(
        std_out,
        &terminal_dimensions,
    );
    try stdout.print("{}\n", .{terminal_dimensions});
    var termos = switch (os_tag) {
        .windows => OldState{ .win = .{ .std_in = undefined, .std_out = undefined } },
        .linux, .macos => OldState{ .posix = .{ .std_in = undefined } },
        else => {
            return error.UNSUPPORTED_OS;
        },
    };
    try start_raw_mode(std_in, std_out, &termos);
    defer restore_terminal(std_in, std_out, termos);

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
