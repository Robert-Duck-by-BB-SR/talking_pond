const std = @import("std");
const RenderQ = @import("../../RenderQueue.zig");
const common = @import("../common.zig");
const render_utils = @import("../render_utils.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

alloc: std.mem.Allocator,
render_q: *RenderQ,

rows_to_render: []Row = undefined,
ponds_list: std.ArrayList(PondItem) = undefined,

const Row = struct {
    cursor: []u8 = undefined,
    content: std.ArrayList(u8) = undefined,
};

const PondItem = struct {
    ducks_count: u8,
    has_update: bool,
    title: []const u8 = undefined,
};

const duck_count_label = "Ducks: ";
const notification_ico = "\u{25FC}";

const Self = @This();

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) !Self {
    // Note: min capacity for ArrayList is total height / 2 because pond item is 2 rows height
    var ponds_list: std.ArrayList(PondItem) = try .initCapacity(alloc, @intFromFloat(@as(f16, @floatFromInt(terminal_dimensions.height)) * 0.5));
    const pond_item_one: PondItem = .{ .title = "YAPPING IS BACK", .ducks_count = 3, .has_update = false };
    const pond_item_two: PondItem = .{ .title = "HELL YEAH", .ducks_count = 2, .has_update = true };
    try ponds_list.append(pond_item_one);
    try ponds_list.append(pond_item_two);

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
    self.rows_to_render = try self.alloc.alloc(Row, @intCast(self.dimensions.height));
    const width: usize = @intCast(self.dimensions.width - 2);

    // NOTE: TODO: now, after initiallization we will only have to replace the border with another kind (Normal|Bold|Rounded?)
    // and retain the capacity, which means no additional allocations needed
    var horizontal_border_list: std.ArrayList(u8) = try .initCapacity(self.alloc, width * common.theme.border.HORIZONTAL.len);
    const top_border = try render_utils.render_border_top_with_title(self.alloc, self.dimensions.width, "PONDS", &horizontal_border_list);
    const bottom_border = try render_utils.render_border_bottom(self.alloc, self.dimensions.width, &horizontal_border_list);

    const bg_mid = try self.alloc.alloc(u8, width);
    @memset(bg_mid, ' ');
    const bg = try std.fmt.allocPrint(
        self.alloc,
        "{s}{s}{s}{s}",
        .{
            common.theme.border.VERTICAL,
            bg_mid,
            common.theme.border.VERTICAL,
            common.RESET_STYLES,
        },
    );

    // Top border
    self.rows_to_render[0].cursor = try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ 1, self.position.col });
    self.rows_to_render[0].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, top_border);

    // Sidebar items
    var arr_index: u8 = 1;
    var row_index: u8 = 2;
    for (self.ponds_list.items) |item| {
        var title: []u8 = undefined;
        if (item.has_update) {
            title = try std.fmt.allocPrint(
                self.alloc,
                "{s}{s}{s}{s}{s}{s}{s}{s}",
                .{
                    common.theme.border.VERTICAL,
                    try render_utils.render_line_of_text_and_backround(self.alloc, item.title, width - 1),
                    common.theme.active_font_color,
                    notification_ico,
                    common.theme.background_color,
                    common.theme.font_color,
                    common.theme.border.VERTICAL,
                    common.RESET_STYLES,
                },
            );
        } else {
            title = try std.fmt.allocPrint(
                self.alloc,
                "{s}{s}{s}{s}",
                .{
                    common.theme.border.VERTICAL,
                    try render_utils.render_line_of_text_and_backround(self.alloc, item.title, width),
                    common.theme.border.VERTICAL,
                    common.RESET_STYLES,
                },
            );
        }
        self.rows_to_render[arr_index].cursor = try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ row_index, self.position.col });
        self.rows_to_render[arr_index].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, title);

        row_index += 1;
        arr_index += 1;
    }

    // Background
    for (arr_index..self.rows_to_render.len - 1) |i| {
        self.rows_to_render[i].cursor = try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ row_index, self.position.col });
        self.rows_to_render[i].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, bg);
        row_index += 1;
    }

    // Bottom border
    self.rows_to_render[self.rows_to_render.len - 1].cursor = try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ self.rows_to_render.len, self.position.col });
    self.rows_to_render[self.rows_to_render.len - 1].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, bottom_border);
}

pub fn render(self: Self) !void {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    for (self.rows_to_render) |row| {
        try ponds.writer().print("{s}{s}{s}{s}", .{ row.cursor, common.theme.font_color, common.theme.background_color, row.content.items });
    }
    const slice = try ponds.toOwnedSlice();
    try self.render_q.add_to_render_q(slice);
}
