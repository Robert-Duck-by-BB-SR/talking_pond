const std = @import("std");
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;

const Screen = @import("Screen.zig");
const Ponds = @import("layers/main/Ponds.zig");
const terminal = @import("terminal.zig");
const server = @import("tp_server.zig");

pub fn main() !void {
    const std_out = std.io.getStdOut();

    const std_in = std.io.getStdIn();
    const writer = std_out.writer();

    var gpa = std.heap.GeneralPurposeAllocator(.{}).init;
    var arena = std.heap.ArenaAllocator.init(gpa.allocator());
    defer arena.deinit();

    // NOTE: im leaving it just in case arena will do another two page stack trace
    // var allocator = std.heap.DebugAllocator(.{}).init;
    var screen = Screen.create(arena.allocator());
    try screen.create_layers(std_out);

    // NOTE: we don't want that to happend while debbuging
    // defer writer.print("\x1b[2J", .{}) catch unreachable;
    try screen.init_first_frame();

    var termos = try terminal.get_termos_with_tea();
    try terminal.start_raw_mode(std_in, std_out, &termos);

    defer terminal.restore_terminal(std_in, std_out, writer, termos) catch unreachable;

    const render_thread = try std.Thread.spawn(.{}, Screen.read_terminal, .{ &screen, std_in });
    defer render_thread.join();

    screen.render_q.mutex.lock();
    defer screen.render_q.mutex.unlock();
    while (!screen.exit) {
        try screen.render_q.render(writer);
    }
}
