const std = @import("std");
const Screen = @import("../../Screen.zig");
const common = @import("../common.zig");

dimensions: common.Dimensions = undefined,
bg_line: []u8 = undefined,
terminal_dimensions: *common.Dimensions = undefined,
screen: *Screen = undefined,
alloc: std.mem.Allocator,

const Self = @This();

pub fn create(alloc: std.mem.Allocator, screen: *Screen) !Self {
    return Self{
        .alloc = alloc,
        .screen = screen,
    };
}

pub fn init(self: *Self) !void {
    self.dimensions.width = @intFromFloat(@as(f16, @floatFromInt(self.screen.terminal_dimensions.width)) * 0.3);
    self.dimensions.height = self.screen.terminal_dimensions.height - 1;
    self.bg_line = try self.allocator.alloc(u8, @intCast(self.dimensions.width));
}

const BG_PATTERN = "\x1b[{};{}H\x1b[2K\x1b[48;2;251;206;44m\x1b[38;2;0;0;0m{s}\x1b[0m";

// pub fn render_test(self: *Self) ![]u8 {
//     // ITEM|
//     // ITEM|
//     // ****|
//     // ****|
//     // Might need a border right or border around
//     // self.height = self.screen.terminal_dimensions.height - 1;
//     // self.width = self.screen.terminal_dimensions.width * 0.35;
//     self.dimensions.width = 10;
//     self.dimensions.height = 10;
// }
