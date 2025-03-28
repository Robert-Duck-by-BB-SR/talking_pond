const std = @import("std");
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;

const Screen = @import("internal/Screen.zig");
const terminal = @import("internal/terminal.zig");
const server = @import("internal/tp_server.zig");

pub fn main() !void {
    const std_out = std.io.getStdOut();

    const std_in = std.io.getStdIn();
    const stdout = std_out.writer();

    var gpa = std.heap.DebugAllocator(.{}).init;
    var screen = try Screen.new(gpa.allocator());
    defer screen.destroy();

    try screen.get_terminal_dimensions(std_out);
    try stdout.print("{}\n", .{screen.terminal_dimensions});

    var termos = try terminal.get_termos_with_tea();
    try terminal.start_raw_mode(std_in, std_out, &termos);
    defer terminal.restore_terminal(std_in, std_out, termos);

    const render_thread = try std.Thread.spawn(.{}, Screen.read_terminal, .{ &screen, std_in });
    defer render_thread.join();

    screen.mutex.lock();
    defer screen.mutex.unlock();
    while (screen.render_q.items.len == 0) {
        screen.condition.wait(&screen.mutex);
        if (screen.exit) break;
        try stdout.print("\x1b[48;2;25;60;80m{s}\x1b[0m\n", .{screen.render_q.items});
        screen.render_q.clearAndFree();
    }
}
