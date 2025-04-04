const std = @import("std");
const RenderQ = @import("../../RenderQueue.zig");
const common = @import("../common.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,
rows: []Row = undefined,
alloc: std.mem.Allocator,
render_q: *RenderQ,

const Self = @This();

const Row = struct {
    cursor: []u8 = undefined,
    content: std.ArrayList(u8) = undefined,
};

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) Self {
    return Self{
        .render_q = render_q,
        .alloc = alloc,
        .position = .{ .col = 1, .row = 1 },
        .dimensions = .{
            .width = common.PONDS_SIDEBAR_SIZE,
            .height = terminal_dimensions.height - 1,
        },
    };
}

pub fn init_first_frame(self: *Self) !void {
    self.rows = try self.alloc.alloc(Row, @intCast(self.dimensions.height));
    const width: usize = @intCast(self.dimensions.width - 2);

    // NOTE: TODO: now, after initiallization we will only have to replace the border with another kind (Normal|Bold|Rounded?)
    // and retain the capacity, which means no additional allocations needed
    var horizontal_border: std.ArrayList(u8) = try .initCapacity(self.alloc, width * common.NormalBorder.HORIZONTAL.len);
    var j: usize = 0;
    while (j < self.dimensions.width - 2) {
        defer j += 1;
        horizontal_border.appendSliceAssumeCapacity(common.NormalBorder.HORIZONTAL);
    }

    const top_border = try std.fmt.allocPrint(
        self.alloc,
        "{s}{s}{s}{s}",
        .{
            common.NormalBorder.TOP_LEFT,
            horizontal_border.items,
            common.NormalBorder.TOP_RIGHT,
            common.RESET_STYLES,
        },
    );

    const bottom_border = try std.fmt.allocPrint(
        self.alloc,
        "{s}{s}{s}{s}",
        .{
            common.NormalBorder.BOTTOM_LEFT,
            horizontal_border.items,
            common.NormalBorder.BOTTOM_RIGHT,
            common.RESET_STYLES,
        },
    );

    const bg_mid = try self.alloc.alloc(u8, width);
    @memset(bg_mid, ' ');
    const bg = try std.fmt.allocPrint(
        self.alloc,
        "{s}{s}{s}{s}",
        .{
            common.NormalBorder.VERTICAL,
            bg_mid,
            common.NormalBorder.VERTICAL,
            common.RESET_STYLES,
        },
    );

    for (self.rows, 0..) |*row, i| {
        row.cursor = try std.fmt.allocPrint(self.alloc, "\x1b[{};{}H", .{ i + 1, self.position.col });
        if (i == 0) {
            row.content = std.ArrayList(u8).fromOwnedSlice(self.alloc, top_border);
        } else if (i == self.rows.len - 1) {
            row.content = std.ArrayList(u8).fromOwnedSlice(self.alloc, bottom_border);
        } else {
            row.content = std.ArrayList(u8).fromOwnedSlice(self.alloc, bg);
        }
    }
}

pub fn render(self: Self) !void {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    for (self.rows) |row| {
        try ponds.writer().print("{s}{s}{s}{s}", .{ row.cursor, common.PRIMARY_THEME.font_color, common.PRIMARY_THEME.background_color, row.content.items });
    }
    const slice = try ponds.toOwnedSlice();
    try self.render_q.add_to_render_q(slice);
}
