const std = @import("std");
const RenderQ = @import("../../RenderQueue.zig");
const common = @import("../common.zig");
const render_utils = @import("../render_utils.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

alloc: std.mem.Allocator,
render_q: *RenderQ,

rows_to_render: []Row = undefined,
border: []u8 = undefined,
ponds_list: std.ArrayList(PondItem) = undefined,
active_pond: usize = 0,
is_active: bool = true,

const Row = struct {
    cursor: []u8 = undefined,
    content: []u8 = undefined,
};

const PondItem = struct {
    id: []u8 = undefined,
    has_update: bool = false,
    title: []const u8 = undefined,
};

const Self = @This();

fn wrapi(index: usize, direction: isize, max: usize) usize {
    if (direction == -1 and index == 0) {
        return max - 1;
    } else if (direction == 1 and index == max - 1) {
        return 0;
    } else {
        const s_index: isize = @intCast(index);
        return @intCast(s_index + direction);
    }
}

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) !Self {
    // NOTE: -2 accounts for borders
    var ponds_list: std.ArrayList(PondItem) = try .initCapacity(
        alloc,
        @intCast(terminal_dimensions.height - 2),
    );

    const pond_item_one: PondItem = .{ .title = "YAPPING IS BACK", .has_update = false };
    const pond_item_two: PondItem = .{ .title = "HELL YEAH", .has_update = true };
    const pond_item_three: PondItem = .{ .title = "Babagi with a capital G", .has_update = false };
    const pond_item_four: PondItem = .{ .title = "GITGOOD / fix skill issue (same thing)", .has_update = true };

    try ponds_list.append(pond_item_one);
    try ponds_list.append(pond_item_two);
    try ponds_list.append(pond_item_three);
    try ponds_list.append(pond_item_four);

    return Self{
        .render_q = render_q,
        .alloc = alloc,
        .position = .{ .col = 1, .row = 1 },
        .dimensions = .{
            .width = common.PONDS_SIDEBAR_SIZE,
            .height = terminal_dimensions.height - 1,
        },
        .ponds_list = ponds_list,
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
        "PONDS",
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
                    self.dimensions.width,
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

pub fn remap_content(self: *Self) !void {
    for (self.ponds_list.items, 0..) |pond, i| {
        const content = try render_utils.render_line_of_text_and_backround(
            self.alloc,
            pond.title,
            @intCast(self.dimensions.width - 2),
        );
        @memcpy(self.rows_to_render[i].content[0..content.len], content);
    }
}

fn render_row(self: *Self, row_index: usize) ![]u8 {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    const row = self.rows_to_render[row_index];
    try ponds.writer().print("{s}{s}{s}{s}{s}", .{
        row.cursor,
        if (self.ponds_list.items.len != 0 and row_index == self.active_pond) common.ACTIVE_ITEM else common.INACTIVE_ITEM,
        row.content,
        try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ row_index + 2, self.dimensions.width - 1 }),
        if (row_index < self.ponds_list.items.len and
            self.ponds_list.items[row_index].has_update) common.NOTIFICATION_ICON_PATTERN else "",
    });
    return ponds.toOwnedSlice();
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
    const slice = try ponds.toOwnedSlice();
    try self.render_q.add_to_render_q(slice);
}

pub fn handle_normal(self: *Self, key: u8) !void {
    switch (key) {
        'j' => {
            const prev_pond = self.active_pond;
            self.active_pond = wrapi(self.active_pond, 1, self.ponds_list.items.len);
            const old_pond = try self.render_row(prev_pond);
            const new_pond = try self.render_row(self.active_pond);
            try self.render_q.add_to_render_q(try std.fmt.allocPrint(self.alloc, "{s}{s}", .{ old_pond, new_pond }));
        },
        'k' => {
            const prev_pond = self.active_pond;
            self.active_pond = wrapi(self.active_pond, -1, self.ponds_list.items.len);
            const old_pond = try self.render_row(prev_pond);
            const new_pond = try self.render_row(self.active_pond);
            try self.render_q.add_to_render_q(try std.fmt.allocPrint(self.alloc, "{s}{s}", .{ old_pond, new_pond }));
        },
        else => {},
    }
}
