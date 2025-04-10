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

const ACTIVE_ITEM = common.theme.FONT_COLOR ++ common.theme.ACTIVE_BACKGROUND_COLOR;
const INACTIVE_ITEM = common.theme.FONT_COLOR ++ common.theme.BACKGROUND_COLOR;
const BORDER_OFFSET = common.theme.BORDER.VERTICAL.len;

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
    self.rows_to_render = try self.alloc.alloc(Row, @intCast(self.dimensions.height));
    const width: usize = @intCast(self.dimensions.width - 2);

    // NOTE: TODO: now, after initiallization we will only have to replace the border with another kind (Normal|Bold|Rounded?)
    // and retain the capacity, which means no additional allocations needed
    var horizontal_border_list: std.ArrayList(u8) = try .initCapacity(self.alloc, width * common.theme.BORDER.HORIZONTAL.len);
    const top_border = try render_utils.render_border_top_with_title(self.alloc, self.dimensions.width, "PONDS", &horizontal_border_list);
    const bottom_border = try render_utils.render_border_bottom(self.alloc, self.dimensions.width, &horizontal_border_list);
    std.debug.print("{}\n", .{self.rows_to_render.len});

    // Top border
    self.rows_to_render[0].cursor = try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{ 1, self.position.col });
    self.rows_to_render[0].content = std.ArrayList(u8).fromOwnedSlice(self.alloc, top_border);

    // Background
    for (1..self.rows_to_render.len - 1) |i| {
        // TODO: rethink borders part maybe?
        const bg_mid = try self.alloc.alloc(u8, width);
        @memset(bg_mid, ' ');
        const bg = try std.fmt.allocPrint(
            self.alloc,
            "{s}{s}{s}{s}",
            .{
                common.theme.BORDER.VERTICAL,
                bg_mid,
                common.theme.BORDER.VERTICAL,
                common.RESET_STYLES,
            },
        );
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

pub fn remap_content(self: *Self) !void {
    for (self.ponds_list.items, 1..) |item, i| {
        const content = try render_utils.render_line_of_text_and_backround(
            self.alloc,
            item.title,
            @intCast(self.dimensions.width - 2),
        );
        std.debug.print("{s} {d} {d}\n", .{ content, content.len, i });
        @memcpy(self.rows_to_render[i].content.items[BORDER_OFFSET .. content.len + BORDER_OFFSET], content);
        std.debug.print("{s} |||| {s} {d} {d}\n", .{ self.rows_to_render[i].content.items, self.rows_to_render[i].content.items, content.len, i });
    }
}

// fn render_row(self: *Self, row_index: usize) !void { }

pub fn render(self: *Self) !void {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    try self.remap_content();
    for (self.rows_to_render, 0..) |row, i| {
        try ponds.writer().print("{s}{s}{s}{s}{s}", .{
            row.cursor,
            if (self.ponds_list.items.len != 0 and i == self.active_pond) ACTIVE_ITEM else INACTIVE_ITEM,
            row.content.items,
            try std.fmt.allocPrint(self.alloc, common.MOVE_CURSOR_TO_POSITION, .{i + 1, self.dimensions.width - 1}),
            if ((i > 0 and i < self.ponds_list.items.len + 1) and self.ponds_list.items[i - 1].has_update) common.NOTIFICATION_ICON_PATTERN else "",
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
