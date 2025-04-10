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
active_pond: usize = 1,

const ACTIVE_ITEM = common.theme.active_font_color ++ common.theme.active_background_color;
const INACTIVE_ITEM = common.theme.font_color ++ common.theme.background_color;

const Row = struct {
    cursor: []u8 = undefined,
    content: std.ArrayList(u8) = undefined,
};

const PondItem = struct {
    id: []u8 = undefined,
    has_update: bool = false,
    title: []const u8 = undefined,
};

const Self = @This();

fn wrapi(index: usize, max: usize) usize {
    if (index == 0) {
        return max;
    } else if (index > max) {
        return 1;
    } else {
        return index;
    }
}

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) !Self {
    // NOTE: -2 accounts for borders
    const height: usize = @intCast(terminal_dimensions.height);
    var ponds_list: std.ArrayList(PondItem) = try .initCapacity(alloc, height - 2);

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

    // Background
    for (1..self.rows_to_render.len - 2) |i| {
        self.rows_to_render[i].cursor = try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ i + 1, self.position.col });
        self.rows_to_render[i].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, bg);
    }

    // Bottom border
    self.rows_to_render[self.rows_to_render.len - 1].cursor = try std.fmt.allocPrint(
        self.alloc,
        common.MOVE_CURSOR_TO_POSITION,
        .{
            self.rows_to_render.len,
            self.position.col,
        },
    );
    self.rows_to_render[self.rows_to_render.len - 1].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, bottom_border);
}

fn render_pond_item(alloc: std.mem.Allocator, text: []const u8, has_update: bool, width: usize) ![]u8 {
    const result = try std.fmt.allocPrint(
        alloc,
        "{s}{s}",
        .{
            try render_utils.render_line_of_text_and_backround(alloc, text, width),
            if (has_update) common.theme.notification_ico_pattern else "",
        },
    );
    return result;
}

pub fn remap_content(self: *Self) !void {
    for (self.ponds_list.items, 1..) |item, i| {
        const content = try render_pond_item(
            self.alloc,
            item.title,
            item.has_update,
            @intCast(self.dimensions.width),
        );
        std.debug.print("{s} {d}\n", .{content, content.len});
        @memcpy(self.rows_to_render[i].content.items[1..content.len + 1], content);
    }
}


fn render_row(self: *Self, row_index: usize) !void {
    const row = self.rows_to_render[row_index];
    const result = try std.fmt.allocPrint(self.alloc, "{s}{s}{s}{s}", .{ row.cursor, common.theme.font_color, common.theme.background_color, row.content.items });
    try self.render_q.add_to_render_q(result);
}

pub fn render(self: *Self) !void {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    try self.remap_content();
    for (self.rows_to_render, 0..) |row, i| {
        try ponds.writer().print("{s}{s}{s}", .{
            row.cursor,
            if (i == self.active_row) ACTIVE_ITEM else INACTIVE_ITEM,
            row.content.items,
        });
    }
    const slice = try ponds.toOwnedSlice();
    try self.render_q.add_to_render_q(slice);
}

// pub fn handle_normal(self: *Self, key: u8) !void {
//     switch (key) {
//         'j' => {
//             if (self.active_pond + 1 < self.ponds_list.items.len) {
//                 self.active_pond += 1;
//                 // TODO: REMOVE ME AND MAKE GLOBAL
//                 const width: usize = @intCast(self.dimensions.width - 2);
//                 // make me a part of a function
//                 const title = try render_active_pond_item(self.alloc, self.ponds_list.items[active_pond].title, width);
//                 self.rows_to_render[active_pond + 1].cursor = try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ active_pond + 2, self.position.col });
//                 self.rows_to_render[active_pond + 1].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, title);
//                 try render_row(self, active_pond + 1);
//             }
//         },
//         'k' => {
//             if (active_pond - 1 >= 0) {
//                 active_pond -= 1;
//             }
//         },
//         else => {},
//     }
// }
