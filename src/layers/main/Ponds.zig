const std = @import("std");
const Screen = @import("../../Screen.zig");
const common = @import("../common.zig");

dimensions: common.Dimensions = undefined,
bg_line: []u8 = undefined,
screen: *Screen = undefined,
alloc: std.mem.Allocator,

const Self = @This();

pub fn init(alloc: std.mem.Allocator, screen: *Screen) !Self {
    return Self{
        .alloc = alloc,
        .screen = screen,
    };
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

//     // ITS A TOTAL BS, STILL LEARNING ZIG ^_^
//     // self.bg_line = try self.alloc.alloc(u8, @intCast(self.dimensions.width));
//     // @memset(self.bg_line, ' ');
//     // self.screen.add_to_render_q(std.fmt.allocPrint(self.alloc, BG_PATTERN, .{ 1, 1, self.bg_line }));
//     // for (1..self.dimensions.height) |i| {
//     // self.screen.add_to_render_q(std.fmt.allocPrint(self.alloc, BG_PATTERN, .{ i, 1, self.bg_line }));
//     // }
// }
