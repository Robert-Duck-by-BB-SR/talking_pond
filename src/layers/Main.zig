const std = @import("std");
const common = @import("common.zig");
const Ponds = @import("main/Ponds.zig");
const Quacks = @import("main/Quacks.zig");
const Insert = @import("main/Insert.zig");
const RenderQueue = @import("../RenderQueue.zig");

ponds: Ponds,
quacks: Quacks,
// insert: Insert,
alloc: std.mem.Allocator,
render_queue: *RenderQueue,

const Self = @This();

pub fn create(alloc: std.mem.Allocator, render_queue: *RenderQueue, terminal_dimensions: common.Dimensions) !Self {
    const ponds = try Ponds.create(
        alloc,
        terminal_dimensions,
        render_queue,
    );
    const quacks = Quacks.create(
        alloc,
        terminal_dimensions,
        render_queue,
    );
    return Self{ .alloc = alloc, .render_queue = render_queue, .ponds = ponds, .quacks = quacks };
}

pub fn render_first_frame(self: *Self) !void {
    try self.ponds.init_first_frame();
    try self.quacks.init_first_frame();
    try self.ponds.render();
    try self.quacks.render();
}

pub fn handle_current_state(mode: *common.MODE, key: u8) void {
    switch (mode) {
        .INSERT => handle_insert(key),
        .NORMAL => handle_normal(key),
        else => {},
    }
}

fn handle_insert(key: u8) void {
    switch (key) {}
}

pub fn handle_normal(key: u8) void {
    switch (key) {
        'C' => {
            std.debug.print("CONVERAASD");
        },
        'S' => {
            std.debug.print("BOOOOOOBS");
        },
        else => {},
    }
}

pub fn init_first_frame() ![]u8 {}
