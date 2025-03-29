const std = @import("std");
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;

const Screen = @import("Screen.zig");
const terminal = @import("terminal.zig");
const server = @import("tp_server.zig");

pub fn main() !void {
    const std_out = std.io.getStdOut();

    const std_in = std.io.getStdIn();
    const stdout = std_out.writer();

    var debug_allocator = std.heap.DebugAllocator(.{}).init;
    var screen = try Screen.new(debug_allocator.allocator());
    // TODO: REMOVE AFTER CONFIRMING IT WORKS
    // TODO: HOW THE HELL WE NEED TO COMFIRM THAT?
    screen.active_mode = .COMMAND;
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
        if (screen.staying_alive) break;
        try stdout.print("\x1b[48;2;25;60;80m{s}\x1b[0m\n", .{screen.render_q.items});
        screen.render_q.clearAndFree();
    }
}
