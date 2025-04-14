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

const Components = enum { PONDS_SIDEBAR, QUACKS_CHAT, INPUT_FIELD };
var active_component: Components = .PONDS_SIDEBAR;

const Self = @This();

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_queue: *RenderQueue) !Self {
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
    return Self{
        .alloc = alloc,
        .render_queue = render_queue,
        .ponds = ponds,
        .quacks = quacks,
    };
}

pub fn render_first_frame(self: *Self) !void {
    try self.ponds.init_first_frame();
    try self.quacks.init_first_frame();
    try self.ponds.render();
    try self.quacks.render();
}

pub fn handle_current_state(self: *Self, mode: common.MODE, key: u8) !void {
    switch (mode) {
        .NORMAL => try handle_normal(self, key),
        .INSERT => handle_insert(key),
        else => {},
    }
}

fn handle_normal(self: *Self, key: u8) !void {
    switch (key) {
        'S' => {
            active_component = .PONDS_SIDEBAR;
        },
        'C' => {
            active_component = .QUACKS_CHAT;
        },
        'I' => {
            active_component = .INPUT_FIELD;
        },
        else => {
            switch (active_component) {
                .PONDS_SIDEBAR => {
                    try self.ponds.handle_normal(key);
                },
                .QUACKS_CHAT => {},
                .INPUT_FIELD => {},
            }
        },
    }
}

fn handle_insert(key: u8) void {
    switch (key) {
        else => {},
    }
}
