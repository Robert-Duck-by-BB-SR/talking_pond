const std = @import("std");
const common = @import("common.zig");
const render_utils = @import("render_utils.zig");
const Ponds = @import("main/Ponds.zig");
const Quacks = @import("main/Quacks.zig");
const Insert = @import("main/Insert.zig");
const RenderQueue = @import("../RenderQueue.zig");

ponds: Ponds,
quacks: Quacks,
insert: Insert,
alloc: std.mem.Allocator,
render_queue: *RenderQueue,
active_component: common.ComponentType = .PONDS_SIDEBAR,

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

pub fn handle_current_state(self: *Self, mode: *common.MODE, key: u8) !void {
    switch (mode.*) {
        .NORMAL => try handle_normal(self, mode, key),
        .INSERT => {
            if (self.active_component == .INPUT_FIELD) {
                try self.insert.handle_insert(mode, key);
            }
        },
        else => {},
    }
}

fn handle_normal(self: *Self, mode: *common.MODE, key: u8) !void {
    var new_active = self.active_component;

    switch (self.active_component) {
        .PONDS_SIDEBAR => {
            try self.ponds.handle_normal(mode, key, &new_active);
        },
        .QUACKS_CHAT => {
            try self.quacks.handle_normal(mode, key, &new_active);
        },
        .INPUT_FIELD => {
            try self.insert.handle_normal(mode, key, &new_active);
        },
    }

    if (new_active != self.active_component) {
        try self.switch_active(new_active);
    }
}

fn switch_active(self: *Self, new_active: common.ComponentType) !void {
    var arena = std.heap.ArenaAllocator.init(self.alloc);
    const temporary_allocator = arena.allocator();
    defer arena.deinit();

    var old_border: []u8 = undefined;
    var updated_old_border: []u8 = undefined;
    var updated_new_border: []u8 = undefined;

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
            try self.render_queue.add_to_render_q(common.HIDDEN_CURSOR, .CURSOR);
        },
    }
    switch (new_active) {
        .PONDS_SIDEBAR => {
            self.ponds.is_active = true;
            updated_old_border = try render_utils.rerender_border(temporary_allocator, false, old_border);
            updated_new_border = try render_utils.rerender_border(temporary_allocator, true, self.ponds.border);
        },
        .QUACKS_CHAT => {
            self.quacks.is_active = true;
            updated_old_border = try render_utils.rerender_border(temporary_allocator, false, old_border);
            // TODO: handle unselected pond
            try self.quacks.render_border_with_title(self.ponds.get_active_pond_title(), temporary_allocator);
            updated_new_border = try render_utils.rerender_border(temporary_allocator, true, self.quacks.border);
        },
        .INPUT_FIELD => {
            self.quacks.is_active = true;
            try self.render_queue.add_to_render_q(common.VISIBLE_CURSOR, .CURSOR);
            try self.render_queue.add_to_render_q(self.insert.render_current_virtual_cursor(), .CURSOR);
            updated_old_border = try render_utils.rerender_border(temporary_allocator, false, old_border);
            updated_new_border = try render_utils.rerender_border(temporary_allocator, true, self.insert.border);
        },
    }
    self.active_component = new_active;

    try self.render_queue.add_to_render_q(updated_old_border, .CONTENT);
    try self.render_queue.add_to_render_q(updated_new_border, .CONTENT);
    self.render_queue.sudo_render();
}
