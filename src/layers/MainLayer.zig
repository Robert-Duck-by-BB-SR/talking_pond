const std = @import("std");
const common = @import("common.zig");
const Ponds = @import("main/Ponds.zig");
const Quacks = @import("main/Quacks.zig");
const Insert = @import("main/Insert.zig");
const RenderQueue = @import("../RenderQueue.zig");

ponds: Ponds,
quacks: Quacks,
insert: Insert,
alloc: std.mem.Allocator,
render_queue: *RenderQueue,
active_component: ComponentType = .PONDS_SIDEBAR,

const ComponentType = enum { PONDS_SIDEBAR, QUACKS_CHAT, INPUT_FIELD };

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
    const insert = Insert.create(
        alloc,
        .{
            .row = quacks.dimensions.height + 1,
            .col = quacks.position.col,
        },
        .{
            .height = 5,
            .width = quacks.dimensions.width,
        },
        render_queue,
    );
    return Self{
        .alloc = alloc,
        .render_queue = render_queue,
        .ponds = ponds,
        .quacks = quacks,
        .insert = insert,
    };
}

pub fn render_first_frame(self: *Self) !void {
    try self.ponds.init_first_frame();
    try self.quacks.init_first_frame();
    try self.insert.init_first_frame();
    try self.ponds.render();
    try self.quacks.render();
    try self.insert.render();
}

pub fn handle_current_state(self: *Self, mode: common.MODE, key: u8) !void {
    switch (mode) {
        .NORMAL => try handle_normal(self, key),
        .INSERT => {},
        else => {},
    }
}

fn handle_normal(self: *Self, key: u8) !void {
    switch (key) {
        'P' => {
            try self.switch_active(.PONDS_SIDEBAR);
        },
        'Q' => {
            try self.switch_active(.QUACKS_CHAT);
        },
        'I' => {
            try self.switch_active(.INPUT_FIELD);
        },
        else => {
            switch (self.active_component) {
                .PONDS_SIDEBAR => {
                    try self.ponds.handle_normal(key);
                },
                .QUACKS_CHAT => {},
                .INPUT_FIELD => {},
            }
        },
    }
}

fn switch_active(self: *Self, new_active: ComponentType) !void {
    var old_border: []u8 = undefined;
    var new_border: []u8 = undefined;
    switch (self.active_component) {
        .PONDS_SIDEBAR => {
            self.ponds.is_active = false;
            old_border = self.ponds.border;
        },
        .QUACKS_CHAT => {
            self.quacks.is_active = false;
            old_border = self.quacks.border;
        },
        .INPUT_FIELD => {
            self.insert.is_active = false;
            old_border = self.insert.border;
        },
    }
    switch (new_active) {
        .PONDS_SIDEBAR => {
            self.ponds.is_active = true;
            new_border = self.ponds.border;
        },
        .QUACKS_CHAT => {
            self.quacks.is_active = true;
            new_border = self.quacks.border;
        },
        .INPUT_FIELD => {
            self.quacks.is_active = true;
            new_border = self.insert.border;
        },
    }
    self.active_component = new_active;
    const compiled_old_border = try common.render_border(self.alloc, false, old_border);
    const compiled_new_border = try common.render_border(self.alloc, true, new_border);
    try self.render_queue.add_to_render_q(compiled_old_border);
    try self.render_queue.add_to_render_q(compiled_new_border);
}
