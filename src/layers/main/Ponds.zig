const std = @import("std");
const RenderQ = @import("../../RenderQueue.zig");
const common = @import("../common.zig");
const render_utils = @import("../render_utils.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

main_allocator: std.mem.Allocator = undefined,
temporary_allocator: std.mem.Allocator = undefined,
render_q: *RenderQ,

content: []Row = undefined,
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

pub fn create(parent_allocator: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) !Self {
    // NOTE: -2 accounts for borders
    const ponds_list: std.ArrayList(PondItem) = try .initCapacity(
        parent_allocator,
        @intCast(terminal_dimensions.height - 2),
    );

    // const pond_item_one: PondItem = .{ .title = "YAPPING IS BACK", .has_update = false };
    // const pond_item_two: PondItem = .{ .title = "HELL YEAH", .has_update = true };
    // const pond_item_three: PondItem = .{ .title = "Babagi with a capital G", .has_update = false };
    // const pond_item_four: PondItem = .{ .title = "GITGOOD / fix skill issue (same thing)", .has_update = true };
    //
    // try ponds_list.append(pond_item_one);
    // try ponds_list.append(pond_item_two);
    // try ponds_list.append(pond_item_three);
    // try ponds_list.append(pond_item_four);

    return Self{
        .render_q = render_q,
        .main_allocator = parent_allocator,
        .position = .{ .col = 1, .row = 1 },
        .dimensions = .{
            .width = common.PONDS_SIDEBAR_SIZE,
            .height = terminal_dimensions.height - 1,
        },
        .ponds_list = ponds_list,
    };
}

/// Renders borders and background
pub fn init_first_frame(self: *Self) !void {
    var arena = std.heap.ArenaAllocator.init(self.main_allocator);
    defer arena.deinit();
    const temporary_alloctor = arena.allocator();

    try self.render_border(temporary_alloctor);

    // Background
    self.content = try temporary_alloctor.alloc(Row, @intCast(self.dimensions.height - 2));
    for (self.content, 2..) |*row, i| {
        const bg_mid = try self.main_allocator.alloc(u8, @intCast(self.dimensions.width - 2));
        @memset(bg_mid, ' ');
        row.cursor = try std.fmt.allocPrint(
            temporary_alloctor,
            common.MOVE_CURSOR_TO_POSITION,
            .{ i, self.position.col + 1 },
        );
        row.content = bg_mid;
    }
}

fn render_border(self: *Self, temporary_allocator: std.mem.Allocator) !void {
    const width: usize = @intCast(self.dimensions.width - 2);
    const corners_width = common.theme.BORDER.BOTTOM_LEFT.len + common.theme.BORDER.BOTTOM_RIGHT.len;
    const border_width = width * common.theme.BORDER.HORIZONTAL.len + corners_width;

    // Top border
    const top_border = try render_utils.make_border_with_title(
        temporary_allocator,
        @intCast(self.dimensions.width),
        "PONDS",
    );

    self.border = try std.fmt.allocPrint(temporary_allocator, "{s}{s}", .{
        try std.fmt.allocPrint(
            temporary_allocator,
            common.MOVE_CURSOR_TO_POSITION,
            .{ 1, self.position.col },
        ),
        top_border,
    });

    for (1..@intCast(self.dimensions.height - 1)) |i| {
        self.border = try std.fmt.allocPrint(temporary_allocator, "{s}{s}{s}{s}{s}", .{
            self.border,
            try std.fmt.allocPrint(
                temporary_allocator,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    i + 1,
                    self.position.col,
                },
            ),
            common.theme.BORDER.VERTICAL,
            try std.fmt.allocPrint(
                temporary_allocator,
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
    const bottom_border = try render_utils.make_bottom_border(
        temporary_allocator,
        border_width,
    );

    self.border = try std.fmt.allocPrint(
        temporary_allocator,
        "{s}{s}{s}{s}",
        .{
            self.border,
            try std.fmt.allocPrint(
                temporary_allocator,
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

pub fn fill_content_with_ponds(self: *Self, temporary_alloctor: std.mem.Allocator) !void {
    if (self.ponds_list.items.len == 0) {
        const middle: usize = @intFromFloat(@as(f16, @floatFromInt(self.dimensions.height)) * 0.5);
        const content = try render_utils.render_line_of_text_and_backround(
            temporary_alloctor,
            "DRY LAND",
            common.TEXT_POSITION.CENTER,
            @intCast(self.dimensions.width - 2),
        );
        @memcpy(self.content[middle - 2].content[0..content.len], content);
    }
    for (self.ponds_list.items, 0..) |pond, i| {
        const content = try render_utils.render_line_of_text_and_backround(
            temporary_alloctor,
            pond.title,
            common.TEXT_POSITION.LEFT,
            @intCast(self.dimensions.width - 2),
        );
        @memcpy(self.content[i].content[0..content.len], content);
    }
}

pub fn get_active_pond_title(self: *Self) []const u8 {
    return self.ponds_list.items[self.active_pond].title;
}

fn render_pond_item(self: *Self, row_index: usize, allocator: std.mem.Allocator) ![]u8 {
    var render_result: std.ArrayList(u8) = .init(allocator);
    const row = self.content[row_index];
    try render_result.writer().print("{s}{s}{s}{s}{s}", .{
        row.cursor,
        if (self.ponds_list.items.len != 0 and row_index == self.active_pond) common.ACTIVE_ITEM else common.INACTIVE_ITEM,
        row.content,
        try std.fmt.allocPrint(allocator, common.MOVE_CURSOR_TO_POSITION, .{ row_index + 2, self.dimensions.width - 1 }),
        if (row_index < self.ponds_list.items.len and
            self.ponds_list.items[row_index].has_update) common.NOTIFICATION_ICON_PATTERN else "",
    });
    return render_result.toOwnedSlice();
}

pub fn render(self: *Self) !void {
    var arena = std.heap.ArenaAllocator.init(self.main_allocator);
    defer arena.deinit();
    const temporary_allocator = arena.allocator();
    var render_result: std.ArrayList(u8) = .init(temporary_allocator);
    try self.fill_content_with_ponds(temporary_allocator);
    for (0..self.content.len) |i| {
        try render_result.writer().print("{s}", .{
            try self.render_pond_item(i, temporary_allocator),
        });
    }
    const rendered_border = try render_utils.rerender_border(temporary_allocator, self.is_active, self.border);
    try render_result.writer().print("{s}", .{rendered_border});
    const slice = try render_result.toOwnedSlice();
    try self.render_q.add_to_render_q(slice, .CONTENT);
    self.render_q.sudo_render();
}

pub fn handle_normal(
    self: *Self,
    mode: *common.MODE,
    key: u8,
    new_active: *common.ComponentType,
) !void {
    switch (key) {
        'j' => {
            if (self.ponds_list.items.len == 0) {
                return;
            }
            var arena = std.heap.ArenaAllocator.init(self.main_allocator);
            defer arena.deinit();
            const allocator = arena.allocator();
            const prev_pond = self.active_pond;
            self.active_pond = wrapi(self.active_pond, 1, self.ponds_list.items.len);
            const old_pond = try self.render_pond_item(prev_pond, allocator);
            const new_pond = try self.render_pond_item(self.active_pond, allocator);
            const result = try std.fmt.allocPrint(allocator, "{s}{s}", .{ old_pond, new_pond });
            try self.render_q.add_to_render_q(
                result,
                .CONTENT,
            );
            self.render_q.sudo_render();
        },
        'k' => {
            if (self.ponds_list.items.len == 0) {
                return;
            }
            var arena = std.heap.ArenaAllocator.init(self.main_allocator);
            defer arena.deinit();
            const allocator = arena.allocator();
            const prev_pond = self.active_pond;
            self.active_pond = wrapi(self.active_pond, -1, self.ponds_list.items.len);
            const old_pond = try self.render_pond_item(prev_pond, allocator);
            const new_pond = try self.render_pond_item(self.active_pond, allocator);
            const result = try std.fmt.allocPrint(allocator, "{s}{s}", .{ old_pond, new_pond });
            try self.render_q.add_to_render_q(
                result,
                .CONTENT,
            );
            self.render_q.sudo_render();
        },
        'M', 'Q' => {
            new_active.* = .QUACKS_CHAT;
        },
        'I' => {
            new_active.* = .INPUT_FIELD;
        },
        ':' => {
            mode.* = .COMMAND;
        },
        13 => {
            if (self.ponds_list.items.len == 0) {
                return;
            }
            new_active.* = .QUACKS_CHAT;
        },
        else => {},
    }
}
