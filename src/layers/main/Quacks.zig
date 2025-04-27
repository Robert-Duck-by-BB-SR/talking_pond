const common = @import("../common.zig");
const RenderQ = @import("../../RenderQueue.zig");
const std = @import("std");
const render_utils = @import("../render_utils.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

main_allocator: std.mem.Allocator,
render_q: *RenderQ,

content: []Row = undefined,
border: []u8 = undefined,
active_pond: usize = 0,
is_active: bool = false,

quacks_list: std.ArrayList(QuackItem) = undefined,

const QuackItem = struct {
    id: []u8 = undefined,
    message: []const u8 = undefined,
};

const Row = struct {
    cursor: []u8 = undefined,
    content: []u8 = undefined,
};

const Self = @This();

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) !Self {
    const quacks_list: std.ArrayList(QuackItem) = try .initCapacity(
        alloc,
        // 6 = 1 (status line) + 2 (top and bottom border of input field) + 3 (lines for actual input)
        @intCast(terminal_dimensions.height - 6),
    );

    return Self{
        .render_q = render_q,
        .main_allocator = alloc,
        .position = .{
            .row = 1,
            .col = common.PONDS_SIDEBAR_SIZE + 1,
        },
        .dimensions = .{
            .width = terminal_dimensions.width - common.PONDS_SIDEBAR_SIZE - 1,
            // 6 = 1 (status line) + 2 (top and bottom border of input field) + 3 (lines for actual input)
            .height = terminal_dimensions.height - 6,
        },
        .quacks_list = quacks_list,
    };
}

// salty: TODO: oh wait is this an abstraction???
// carrot: naah bro, trust me, one more abstraction
pub fn init_first_frame(self: *Self) !void {
    var arena = std.heap.ArenaAllocator.init(self.main_allocator);
    defer arena.deinit();
    const temp_allocator = arena.allocator();

    self.content = try temp_allocator.alloc(Row, @intCast(self.dimensions.height - 2));
    try self.render_border_with_title("QUACKS", temp_allocator);
    // Background
    for (self.content, 2..) |*row, i| {
        const bg_mid = try self.main_allocator.alloc(u8, @intCast(self.dimensions.width - 2));
        @memset(bg_mid, ' ');
        row.cursor = try std.fmt.allocPrint(
            temp_allocator,
            common.MOVE_CURSOR_TO_POSITION,
            .{ i, self.position.col + 1 },
        );
        row.content = bg_mid;
    }
}

pub fn render_border_with_title(self: *Self, title: []const u8, temp_allocator: std.mem.Allocator) !void {
    const width: usize = @intCast(self.dimensions.width - 2);
    const corners_width = common.theme.BORDER.BOTTOM_LEFT.len + common.theme.BORDER.BOTTOM_RIGHT.len;
    const border_width = width * common.theme.BORDER.HORIZONTAL.len + corners_width;

    // Top border
    const top_border = try render_utils.make_border_with_title(
        temp_allocator,
        @intCast(self.dimensions.width),
        title,
    );
    self.border = try std.fmt.allocPrint(self.main_allocator, "{s}{s}", .{
        try std.fmt.allocPrint(
            temp_allocator,
            common.MOVE_CURSOR_TO_POSITION,
            .{ 1, self.position.col },
        ),
        top_border,
    });

    for (1..@intCast(self.dimensions.height - 1)) |i| {
        self.border = try std.fmt.allocPrint(temp_allocator, "{s}{s}{s}{s}{s}", .{
            self.border,
            try std.fmt.allocPrint(
                temp_allocator,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    i + 1,
                    self.position.col,
                },
            ),
            common.theme.BORDER.VERTICAL,
            try std.fmt.allocPrint(
                temp_allocator,
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
    const bottom_border = try render_utils.make_bottom_border(
        temp_allocator,
        border_width,
    );
    self.border = try std.fmt.allocPrint(
        temp_allocator,
        "{s}{s}{s}{s}",
        .{
            self.border,
            try std.fmt.allocPrint(
                temp_allocator,
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

pub fn fill_content_with_quacks(self: *Self, temporary_alloctor: std.mem.Allocator) !void {
    if (self.quacks_list.items.len == 0) {
        const middle: usize = @intFromFloat(@as(f16, @floatFromInt(self.dimensions.height)) * 0.5);
        const content = try render_utils.render_line_of_text_and_backround(
            temporary_alloctor,
            "**DEAD SILENCE**",
            common.TEXT_POSITION.CENTER,
            @intCast(self.dimensions.width - 2),
        );
        @memcpy(self.content[middle - 2].content[0..content.len], content);
    }
    for (self.quacks_list.items, 0..) |quack, i| {
        const content = try render_utils.render_line_of_text_and_backround(
            temporary_alloctor,
            quack.message,
            common.TEXT_POSITION.LEFT,
            @intCast(self.dimensions.width - 2),
        );
        @memcpy(self.content[i].content[0..content.len], content);
    }
}

fn render_row(self: *Self, row_index: usize) ![]u8 {
    var render_result: std.ArrayList(u8) = .init(self.main_allocator);
    const row = self.content[row_index];
    try render_result.writer().print("{s}{s}{s}", .{
        row.cursor,
        common.INACTIVE_ITEM,
        row.content,
    });
    return render_result.toOwnedSlice();
}

pub fn render(self: *Self) !void {
    var render_result: std.ArrayList(u8) = .init(self.main_allocator);
    try self.fill_content_with_quacks(self.main_allocator);
    for (0..self.content.len) |i| {
        try render_result.writer().print("{s}", .{
            try self.render_row(i),
        });
    }
    const rendered_border = try render_utils.rerender_border(self.main_allocator, self.is_active, self.border);
    try render_result.writer().print("{s}", .{rendered_border});
    const slice = try render_result.toOwnedSlice();
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
