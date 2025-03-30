const std = @import("std");

width: u16,
height: u16,

content: std.ArrayList(u8),

const Self = @This();

pub fn new(alloc: std.mem.Allocator) !Self {
}
