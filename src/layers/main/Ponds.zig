const std = @import("std");
const Screen = @import("../../Screen.zig");
const common = @import("../common.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,
rows: []Row = undefined,
terminal_dimensions: *common.Dimensions = undefined,
alloc: std.mem.Allocator,

const Self = @This();

const Row = struct {
    cursor: []u8 = undefined,
    content: std.ArrayList(u8) = undefined,
};

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions) Self {
    return Self{ .alloc = alloc, .position = .{ .col = 1, .row = 1 }, .dimensions = .{ .width = @divFloor(terminal_dimensions.width * 30, 100), .height = terminal_dimensions.height - 1 } };
}

pub fn init_first_frame(self: *Self) !void {
    self.rows = try self.alloc.alloc(Row, @intCast(self.dimensions.height));
    const content: []u8 = try self.alloc.alloc(u8, @intCast(self.dimensions.width));
    @memset(content, ' ');
    // ITEM|
    // ITEM|
    // ****|
    // ****|
    // Might need a border right or border around
    for (self.rows, 0..) |*row, i| {
        // \x1b[<i>;<self.position.>H
        row.cursor = try std.fmt.allocPrint(self.alloc, "\x1b[{};{}H", .{ i + 1, self.position.col });
        row.content = std.ArrayList(u8).fromOwnedSlice(self.alloc, content);
    }
}

pub fn render(self: Self) ![]u8 {
    var ponds: std.ArrayList(u8) = .init(self.alloc);
    for (self.rows) |row| {
        try ponds.writer().print("{s}{s}", .{ row.cursor, row.content.items });
    }
    return ponds.toOwnedSlice();
}
