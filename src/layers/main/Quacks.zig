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

const Row = struct {
    cursor: []u8 = undefined,
    content: []u8 = undefined,
};

const Self = @This();

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) Self {
    return Self{
        .render_q = render_q,
        .alloc = alloc,
        .position = .{
            .row = 1,
            .col = common.PONDS_SIDEBAR_SIZE + 1,
        },
        .dimensions = .{
            .width = terminal_dimensions.width - common.PONDS_SIDEBAR_SIZE - 1,
            // 6 = 1 (status line) + 2 (top and bottom border of input field) + 3 (lines for actual input)
            .height = terminal_dimensions.height - 6,
        },
    };
}

// salty: TODO: oh wait is this an abstraction???
// carrot: naah bro, trust me, one more abstraction
pub fn render_first_frame(self: *Self, title: []const u8) !void {
    self.rows_to_render = try self.alloc.alloc(Row, @intCast(self.dimensions.height - 2));
    const width: usize = @intCast(self.dimensions.width - 2);

    const corners_width = common.theme.BORDER.BOTTOM_LEFT.len + common.theme.BORDER.BOTTOM_RIGHT.len;
    const border_width = width * common.theme.BORDER.HORIZONTAL.len + corners_width;

    const top_border = try render_utils.make_border_with_title(
        self.alloc,
        @intCast(self.dimensions.width),
        title,
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
            .{ 1, self.position.col },
        ),
        top_border,
    });

    // Background
    for (self.rows_to_render, 2..) |*row, i| {
        const bg_mid = try self.alloc.alloc(u8, width);
        @memset(bg_mid, ' ');
        row.cursor = try std.fmt.allocPrint(
            self.alloc,
            common.MOVE_CURSOR_TO_POSITION,
            .{ i, self.position.col + 1 },
        );
        row.content = bg_mid;
    }

    for (1..@intCast(self.dimensions.height - 1)) |i| {
        self.border = try std.fmt.allocPrint(self.alloc, "{s}{s}{s}{s}{s}", .{
            self.border,
            try std.fmt.allocPrint(
                self.alloc,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    i + 1,
                    self.position.col,
                },
            ),
            common.theme.BORDER.VERTICAL,
            try std.fmt.allocPrint(
                self.alloc,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    i + 1,
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
                    self.dimensions.height,
                    self.position.col,
                },
            ),
            bottom_border,
            common.RESET_STYLES,
        },
    );
}

// pub fn remap_content(self: *Self) !void {
//     for (self.ponds_list.items, 0..) |pond, i| {
//         const content = try render_utils.render_line_of_text_and_backround(
//             self.alloc,
//             pond.title,
//             @intCast(self.dimensions.width - 2),
//         );
//         @memcpy(self.rows_to_render[i].content[0..content.len], content);
//     }
// }

fn render_row(self: *Self, row_index: usize) ![]u8 {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    const row = self.rows_to_render[row_index];
    try ponds.writer().print("{s}{s}{s}", .{
        row.cursor,
        common.INACTIVE_ITEM,
        // if (self.ponds_list.items.len != 0 and row_index == self.active_pond) common.ACTIVE_ITEM else common.INACTIVE_ITEM,
        row.content,
    });
    return ponds.toOwnedSlice();
}

pub fn render(self: *Self) !void {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    // try self.remap_content();
    for (0..self.rows_to_render.len) |i| {
        try ponds.writer().print("{s}", .{
            try self.render_row(i),
        });
    }
    const rendered_border = try render_utils.render_border(self.alloc, self.is_active, self.border);
    try ponds.writer().print("{s}", .{rendered_border});
    const slice = try ponds.toOwnedSlice();
    try self.render_q.add_to_render_q(slice, .CONTENT);
    self.render_q.sudo_render();
}

pub fn handle_normal(_: *Self, mode: *common.MODE, key: u8, new_active: *common.ComponentType) !void {
    switch (key) {
        'C', 'P' => {
            new_active.* = .PONDS_SIDEBAR;
        },
        'I' => {
            new_active.* = .INPUT_FIELD;
        },
        ':' => {
            mode.* = .COMMAND;
        },
        else => {},
    }
}
