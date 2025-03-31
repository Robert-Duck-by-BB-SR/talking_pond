const std = @import("std");
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;

const Screen = @import("Screen.zig");
const Ponds = @import("./layers/main/Ponds.zig");
const terminal = @import("terminal.zig");
const server = @import("tp_server.zig");

pub fn main() !void {
    const std_out = std.io.getStdOut();

    const std_in = std.io.getStdIn();
    const stdout = std_out.writer();

    var debug_allocator = std.heap.DebugAllocator(.{}).init;
    defer {
        switch (debug_allocator.deinit()) {
            .ok => {},
            .leak => {
                _ = debug_allocator.detectLeaks();
            },
        }
    }
    var screen = try Screen.new(debug_allocator.allocator());
    defer screen.destroy();

    try screen.get_terminal_dimensions(std_out);

    // NOTE: we don't want that to happend while debbuging
    // defer stdout.print("\x1b[2J", .{}) catch unreachable;
    try screen.init_first_frame(stdout);

    var termos = try terminal.get_termos_with_tea();
    try terminal.start_raw_mode(std_in, std_out, &termos);

    var ponds = try Ponds.new(debug_allocator.allocator(), &screen);
    try ponds.render_test();
    defer screen.destroy();

    defer terminal.restore_terminal(std_in, std_out, termos);

    const render_thread = try std.Thread.spawn(.{}, Screen.read_terminal, .{ &screen, std_in });
    defer render_thread.join();

    screen.mutex.lock();
    defer screen.mutex.unlock();
    while (!screen.exit) {
        screen.condition.wait(&screen.mutex);
        try stdout.print("{s}", .{screen.render_q.items});
        screen.render_q.clearAndFree();
    }
}
