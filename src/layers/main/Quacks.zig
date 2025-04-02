const common = @import("../common.zig");
const RenderQ = @import("../../RenderQueue.zig");
const std = @import("std");

dimensions: common.Dimensions,
position: common.Position,
rows: []Row = undefined,
render_q: *RenderQ,
alloc: std.mem.Allocator,

const Self = @This();

const Row = struct {
    cursor: []u8 = undefined,
    content: std.ArrayList(u8) = undefined,
};

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) Self {
    return Self{
        .render_q = render_q,
        .alloc = alloc,
        .position = .{
            .row = 1,
            .col = @divFloor(terminal_dimensions.width * 30, 100),
        },
        .dimensions = .{
            .width = @divFloor(terminal_dimensions.width * 70, 100),
            // 6 = 1 (status line) + 2 (top and bottom border of input field) + 3 (lines for actual input)
            .height = terminal_dimensions.height - 6,
        },
    };
}

pub fn init_first_frame(self: *Self) !void {
    self.rows = try self.alloc.alloc(Row, @intCast(self.dimensions.height));
    const content: []u8 = try self.alloc.alloc(u8, @intCast(self.dimensions.width));
    @memset(content, ' ');

    for (self.rows, 0..) |*row, i| {
        // \x1b[<i>;<self.position.>H
        row.cursor = try std.fmt.allocPrint(self.alloc, "\x1b[{};{}H", .{ i + 1, self.position.col });
        row.content = std.ArrayList(u8).fromOwnedSlice(self.alloc, content);
    }
}

pub fn render(self: Self) !void{
    var quacks: std.ArrayList(u8) = .init(self.alloc);
    for (self.rows) |row| {
        try quacks.writer().print("{s}{s}", .{ row.cursor, row.content.items });
    }
    const slice = try quacks.toOwnedSlice();
    try self.render_q.add_to_render_q(slice);
}
