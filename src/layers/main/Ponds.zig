const std = @import("std");
const RenderQ = @import("../../RenderQueue.zig");
const common = @import("../common.zig");
const render_utils = @import("../render_utils.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

main_allocator: std.mem.Allocator = undefined,
temporary_allocator: std.mem.Allocator = undefined,
render_q: *RenderQ,

rows: []Row = undefined,
border: []u8 = undefined,
ponds_list: std.ArrayList(PondItem) = undefined,
active_pond_index: usize = 0,
sliding_window_move_by: u8 = 0,
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
    var ponds_list: std.ArrayList(PondItem) = try .initCapacity(
        parent_allocator,
        @intCast(terminal_dimensions.height - 2),
    );

    const pond_item_one: PondItem = .{ .title = "YAPPING IS BACK", .has_update = true };
    const pond_item_two: PondItem = .{ .title = "HELL YEAH", .has_update = true };
    const pond_item_three: PondItem = .{ .title = "Babagi with a capital G", .has_update = false };
    const pond_item_four: PondItem = .{ .title = "GITGOOD / fix skill issue (same thing)", .has_update = true };
    const pond_item_last_one: PondItem = .{ .title = "LASTx1", .has_update = true };
    const pond_item_last_two: PondItem = .{ .title = "LASTx2", .has_update = false };

    try ponds_list.append(pond_item_one);
    try ponds_list.append(pond_item_two);
    try ponds_list.append(pond_item_three);
    try ponds_list.append(pond_item_four);
    try ponds_list.append(pond_item_one);
    try ponds_list.append(pond_item_two);
    try ponds_list.append(pond_item_three);
    try ponds_list.append(pond_item_four);
    try ponds_list.append(pond_item_one);
    try ponds_list.append(pond_item_two);
    try ponds_list.append(pond_item_two);
    try ponds_list.append(pond_item_one);
    try ponds_list.append(pond_item_two);
    try ponds_list.append(pond_item_three);
    try ponds_list.append(pond_item_four);
    try ponds_list.append(pond_item_one);
    try ponds_list.append(pond_item_two);
    try ponds_list.append(pond_item_three);
    try ponds_list.append(pond_item_four);
    try ponds_list.append(pond_item_last_one);
    try ponds_list.append(pond_item_last_two);
    try ponds_list.append(pond_item_one);

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
    self.rows = try temporary_alloctor.alloc(Row, @intCast(self.dimensions.height - 2));
    for (self.rows, 2..) |*row, i| {
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
        @memcpy(self.rows[middle - 2].content[0..content.len], content);
        return;
    }

    if (self.rows.len >= self.ponds_list.items.len) {
        for (self.ponds_list.items, 0..) |pond, i| {
            const line = try render_utils.render_line_of_text_and_backround(
                temporary_alloctor,
                pond.title,
                common.TEXT_POSITION.LEFT,
                @intCast(self.dimensions.width - 2),
            );
            @memcpy(self.rows[i].content[0..line.len], line);
        }
    } else {
        // Reached the end end of ponds_list
        if (self.active_pond_index > self.ponds_list.items.len - 1) {
            return;
        }

        // Refilling content for scrolling
        const slice = self.ponds_list.items[self.sliding_window_move_by .. self.rows.len + self.sliding_window_move_by];
        for (slice, 0..) |pond, i| {
            const line = try render_utils.render_line_of_text_and_backround(
                temporary_alloctor,
                pond.title,
                common.TEXT_POSITION.LEFT,
                @intCast(self.dimensions.width - 2),
            );
            @memcpy(self.rows[i].content[0..line.len], line);
        }
    }
}

pub fn get_active_pond_title(self: *Self) []const u8 {
    return self.ponds_list.items[self.active_pond_index].title;
}

fn render_pond_item(self: *Self, content_row_index: usize, allocator: std.mem.Allocator) ![]u8 {
    var render_result: std.ArrayList(u8) = .init(allocator);
    const row = self.rows[content_row_index];
    try render_result.writer().print("{s}{s}{s}{s}{s}", .{
        row.cursor,
        self.set_highlight_styles(content_row_index),
        row.content,
        try std.fmt.allocPrint(allocator, common.MOVE_CURSOR_TO_POSITION, .{ content_row_index + 2, self.dimensions.width - 1 }),
        self.render_pond_notification_icon(content_row_index + self.sliding_window_move_by),
    });
    return render_result.toOwnedSlice();
}

fn set_highlight_styles(self: *Self, content_row_index: usize) []const u8 {
    if (self.ponds_list.items.len != 0 and content_row_index + self.sliding_window_move_by == self.active_pond_index) {
        return common.ACTIVE_ITEM;
    }
    return common.INACTIVE_ITEM;
}

fn render_pond_notification_icon(self: *Self, pond_index: usize) []const u8 {
    if (pond_index < self.ponds_list.items.len and self.ponds_list.items[pond_index].has_update) {
        return common.NOTIFICATION_ICON_PATTERN;
    }
    return "";
}

pub fn render(self: *Self) !void {
    var arena = std.heap.ArenaAllocator.init(self.main_allocator);
    defer arena.deinit();
    const temporary_allocator = arena.allocator();
    var render_result: std.ArrayList(u8) = .init(temporary_allocator);

    try self.fill_content_with_ponds(temporary_allocator);
    const rendered_ponds = try self.render_ponds(temporary_allocator);

    try render_result.writer().print("{s}", .{rendered_ponds});
    const rendered_border = try render_utils.rerender_border(temporary_allocator, self.is_active, self.border);
    try render_result.writer().print("{s}", .{rendered_border});
    const slice = try render_result.toOwnedSlice();
    try self.render_q.add_to_render_q(slice, .CONTENT);
    self.render_q.sudo_render();
}

fn render_ponds(self: *Self, temporary_allocator: std.mem.Allocator) ![]u8 {
    var render_result: std.ArrayList(u8) = .init(temporary_allocator);
    for (0..self.rows.len) |i| {
        try render_result.writer().print("{s}", .{
            try self.render_pond_item(i, temporary_allocator),
        });
    }

    return render_result.toOwnedSlice();
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
            if (self.active_pond_index == self.ponds_list.items.len - 1) {
                return;
            }

            var arena = std.heap.ArenaAllocator.init(self.main_allocator);
            defer arena.deinit();
            const allocator = arena.allocator();

            // Scroll if we need to show ponds which are out of bounds of content
            // Then we need to refill content again with new ponds
            if (self.active_pond_index - self.sliding_window_move_by == self.rows.len - 1) {
                // TODO:" inroduce this feature back
                // self.active_pond_index = wrapi(self.active_pond_index, 1, self.ponds_list.items.len);
                self.active_pond_index += 1;
                self.sliding_window_move_by += 1;
                try self.fill_content_with_ponds(allocator);
                const rendered_ponds = try self.render_ponds(allocator);
                try self.render_q.add_to_render_q(
                    rendered_ponds,
                    .CONTENT,
                );
                self.render_q.sudo_render();
            } else {
                const prev_pond_index = self.active_pond_index - self.sliding_window_move_by;
                // TODO:" inroduce this feature back
                // self.active_pond_index = wrapi(self.active_pond_index, 1, self.ponds_list.items.len);
                self.active_pond_index += 1;
                const result = try std.fmt.allocPrint(allocator, "{s}{s}", .{
                    try self.render_pond_item(prev_pond_index, allocator),
                    try self.render_pond_item(self.active_pond_index - self.sliding_window_move_by, allocator),
                });
                try self.render_q.add_to_render_q(
                    result,
                    .CONTENT,
                );
                self.render_q.sudo_render();
            }
        },
        'k' => {
            if (self.ponds_list.items.len == 0) {
                return;
            }
            if (self.active_pond_index == 0) {
                return;
            }

            var arena = std.heap.ArenaAllocator.init(self.main_allocator);
            defer arena.deinit();
            const allocator = arena.allocator();

            // Move sliding window
            if (self.active_pond_index - self.sliding_window_move_by == 0) {
                self.active_pond_index -= 1;
                self.sliding_window_move_by -= 1;
                try self.fill_content_with_ponds(allocator);

                const rendered_ponds = try self.render_ponds(allocator);
                try self.render_q.add_to_render_q(
                    rendered_ponds,
                    .CONTENT,
                );
                self.render_q.sudo_render();
            } else {
                const prev_pond = self.active_pond_index - self.sliding_window_move_by;
                // TODO:" inroduce this feature back
                self.active_pond_index -= 1;
                const old_pond = try self.render_pond_item(prev_pond, allocator);
                const new_pond = try self.render_pond_item(self.active_pond_index - self.sliding_window_move_by, allocator);
                const result = try std.fmt.allocPrint(allocator, "{s}{s}", .{ old_pond, new_pond });
                try self.render_q.add_to_render_q(
                    result,
                    .CONTENT,
                );

                self.render_q.sudo_render();
            }
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
