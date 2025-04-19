const common = @import("../common.zig");
const RenderQ = @import("../../RenderQueue.zig");
const std = @import("std");
const render_utils = @import("../render_utils.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

alloc: std.mem.Allocator,
render_q: *RenderQ,

rows_to_render: []Row = undefined,
border: []u8 = undefined,
active_pond: usize = 0,
is_active: bool = false,

virtual_cursor: common.Position = .{ .col = 0, .row = 0 },
full_content: std.ArrayList(u8) = undefined,

const Row = struct {
    cursor: []u8 = undefined,
    content: []u8 = undefined,
};

const Self = @This();

pub fn create(alloc: std.mem.Allocator, position: common.Position, dimensions: common.Dimensions, render_q: *RenderQ) Self {
    return Self{
        .render_q = render_q,
        .alloc = alloc,
        .dimensions = dimensions,
        .position = position,
        .full_content = .init(alloc),
    };
}

pub fn init_first_frame(self: *Self) !void {
    self.rows_to_render = try self.alloc.alloc(Row, @intCast(self.dimensions.height - 2));
    const width: usize = @intCast(self.dimensions.width - 2);

    const corners_width = common.theme.BORDER.BOTTOM_LEFT.len + common.theme.BORDER.BOTTOM_RIGHT.len;
    const border_width = width * common.theme.BORDER.HORIZONTAL.len + corners_width;

    const top_border = try render_utils.make_border_with_title(
        self.alloc,
        @intCast(self.dimensions.width),
        "Quack here...",
    );

    const bottom_border = try render_utils.make_bottom_border(
        self.alloc,
        border_width,
    );

    // Top border
    self.border = try std.fmt.allocPrint(self.alloc, "{s}{s}", .{
        try std.fmt.allocPrint(
            self.alloc,
            common.MOVE_CURSOR_TO_POSITION,
            .{ self.position.row, self.position.col },
        ),
        top_border,
    });

    // Background
    const r: usize = @intCast(self.position.row);
    for (self.rows_to_render, 1..) |*row, i| {
        const bg_mid = try self.alloc.alloc(u8, width);
        @memset(bg_mid, ' ');
        row.cursor = try std.fmt.allocPrint(
            self.alloc,
            common.MOVE_CURSOR_TO_POSITION,
            .{ r + i, self.position.col + 1 },
        );
        row.content = bg_mid;
    }

    const row: usize = @intCast(self.position.row);

    for (1..@intCast(self.dimensions.height - 1)) |i| {
        self.border = try std.fmt.allocPrint(self.alloc, "{s}{s}{s}{s}{s}", .{
            self.border,
            try std.fmt.allocPrint(
                self.alloc,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    row + i,
                    self.position.col,
                },
            ),
            common.theme.BORDER.VERTICAL,
            try std.fmt.allocPrint(
                self.alloc,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    row + i,
                    self.position.col + self.dimensions.width - 1,
                },
            ),
            common.theme.BORDER.VERTICAL,
        });
    }

    // Bottom border
    self.border = try std.fmt.allocPrint(
        self.alloc,
        "{s}{s}{s}{s}",
        .{
            self.border,
            try std.fmt.allocPrint(
                self.alloc,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    self.position.row + self.dimensions.height - 1,
                    self.position.col,
                },
            ),
            bottom_border,
            common.RESET_STYLES,
        },
    );
}

fn render_row(self: *Self, row_index: usize) ![]u8 {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    const row = self.rows_to_render[row_index];
    try ponds.writer().print("{s}{s}{s}", .{
        row.cursor,
        common.INACTIVE_ITEM,
        row.content,
    });
    return ponds.toOwnedSlice();
}

fn remap_content(self: *Self) !void {
    // TODO: replace this with actual logic for remap
    @memcpy(self.rows_to_render[0].content[0..self.full_content.items.len], self.full_content.items);
}

pub fn render_current_virtual_cursor(self: *Self) []u8 {
    const actual_position = .{
        self.position.row + self.virtual_cursor.row + 1,
        self.position.col + self.virtual_cursor.col + 1,
    };
    return std.fmt.allocPrint(
        self.alloc,
        common.MOVE_CURSOR_TO_POSITION,
        .{ actual_position[0], actual_position[1] },
    ) catch "";
}

pub fn render(self: *Self) !void {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    try self.remap_content();
    for (0..self.rows_to_render.len) |i| {
        try ponds.writer().print("{s}", .{
            try self.render_row(i),
        });
    }
    const rendered_border = try common.render_border(self.alloc, self.is_active, self.border);
    try ponds.writer().print("{s}", .{rendered_border});
    if (self.is_active) {
        try ponds.writer().print("{s}", .{self.render_current_virtual_cursor()});
    }
    const slice = try ponds.toOwnedSlice();
    try self.render_q.add_to_render_q(slice, .CONTENT);
    self.render_q.sudo_render();
}

pub fn handle_normal(self: *Self, mode: *common.MODE, key: u8, new_active: *common.ComponentType) !void {
    switch (key) {
        'M' => {
            new_active.* = .QUACKS_CHAT;
        },
        'C' => {
            new_active.* = .PONDS_SIDEBAR;
        },
        ':' => {
            mode.* = .COMMAND;
        },
        'a' => {
            if (self.full_content.items.len > 0) {
                self.virtual_cursor.col += 1;
            }
            try self.render_q.add_to_render_q(self.render_current_virtual_cursor(), .CURSOR);
            self.render_q.sudo_render();
            mode.* = .INSERT;
        },
        else => {},
    }
}

pub fn handle_insert(self: *Self, mode: *common.MODE, key: u8) !void {
    switch (key) {
        3 => {
            mode.* = .NORMAL;
            if (self.virtual_cursor.col != 0) {
                self.virtual_cursor.col -= 1;
            }
            try self.render_q.add_to_render_q(self.render_current_virtual_cursor(), .CURSOR);
        },
        else => {
            try self.full_content.append(key);
            try self.remap_content();
            const row = try self.render_row(0);
            try self.render_q.add_to_render_q(row, .CONTENT);
            self.virtual_cursor.col += 1;
            try self.render_q.add_to_render_q(self.render_current_virtual_cursor(), .CURSOR);
            self.render_q.sudo_render();
        },
    }
}
