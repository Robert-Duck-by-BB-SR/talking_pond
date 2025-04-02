const std = @import("std");
queue: std.ArrayList(u8),
mutex: std.Thread.Mutex = .{},
condition: std.Thread.Condition = .{},
alloc: std.mem.Allocator = undefined,

const Self = @This();

pub fn create(alloc: std.mem.Allocator) !Self {
    return Self{
        .alloc = alloc,
        .queue = std.ArrayList(u8).init(alloc),
    };
}

pub fn add_to_render_q(self: *Self, line: []u8) !void {
    self.mutex.lock();
    defer self.mutex.unlock();
    try self.queue.appendSlice(line);
    self.condition.signal();
}

pub fn render(self: *Self, stdout: std.fs.File.Writer) !void {
    self.condition.wait(&self.mutex);
    try stdout.print("{s}", .{self.queue.items});
    self.queue.clearAndFree();
}
