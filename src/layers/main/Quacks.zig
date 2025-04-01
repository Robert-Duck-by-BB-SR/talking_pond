const common = @import("../common.zig");
const std = @import("std");

dimensions: common.Dimensions,
position: common.Position,
rows: []Row = undefined,
alloc: std.mem.Allocator,

const Self = @This();

const Row = struct {
    cursor: []u8 = undefined,
    content: std.ArrayList(u8) = undefined,
};

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions) Self {
    return Self{
        .alloc = alloc,
        .position = .{
            .row = 0,
            .col = @divFloor(terminal_dimensions.width * 35, 100),
        },
        .dimensions = .{
            .width = @divFloor(terminal_dimensions.width * 65, 100),
            .height = terminal_dimensions.height - 5,
        },
    };
}

pub fn init_first_frame(self: *Self) !void {
    self.rows = try self.alloc.alloc(Row, @intCast(self.dimensions.height));
    for (self.rows, 0..) |*row, i| {
        // \x1b[<i>;<self.position.>H
        row.cursor = try std.fmt.allocPrint(self.alloc, "\x1b[{};{}H", .{ i + 1, self.position.col });
        const content: []u8 = try self.alloc.alloc(u8, @intCast(self.dimensions.width));
        @memset(content, ' ');
        row.content = std.ArrayList(u8).fromOwnedSlice(self.alloc, content);
    }
}

pub fn render(self: Self) ![]u8 {
    var quacks: std.ArrayList(u8) = .init(self.alloc);
    for (self.rows) |row| {
        try quacks.writer().print("{s}{s}", .{ row.cursor, row.content.items });
    }
    return quacks.toOwnedSlice();
}
