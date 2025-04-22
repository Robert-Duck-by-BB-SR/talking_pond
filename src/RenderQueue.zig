const std = @import("std");
const CONTENT_TYPE = @import("layers/common.zig").CONTENT_TYPE;

content_queue: std.ArrayList(u8),
status_line_queue: std.ArrayList(u8),
cursor_queue: std.ArrayList(u8),

first_frame: bool = true,

mutex: std.Thread.Mutex = .{},
condition: std.Thread.Condition = .{},
alloc: std.mem.Allocator = undefined,

const Self = @This();

pub fn create(alloc: std.mem.Allocator) Self {
    return Self{
        .alloc = alloc,
        .content_queue = std.ArrayList(u8).init(alloc),
        .status_line_queue = std.ArrayList(u8).init(alloc),
        .cursor_queue = std.ArrayList(u8).init(alloc),
    };
}

pub fn add_to_render_q(self: *Self, line: []const u8, T: CONTENT_TYPE) !void {
    self.mutex.lock();
    defer self.mutex.unlock();
    switch (T) {
        .CONTENT => try self.content_queue.appendSlice(line),
        .CURSOR => try self.cursor_queue.appendSlice(line),
        .STATUS => try self.status_line_queue.appendSlice(line),
    }
}

pub fn sudo_render(self: *Self) void {
    // should there be some logic too?
    self.condition.signal();
}

pub fn render(self: *Self, stdout: std.fs.File.Writer) !void {
    try stdout.print("{s}{s}{s}", .{
        self.content_queue.items,
        self.status_line_queue.items,
        self.cursor_queue.items,
    });
    self.content_queue.clearRetainingCapacity();
    self.status_line_queue.clearRetainingCapacity();
    self.cursor_queue.clearRetainingCapacity();
    self.condition.wait(&self.mutex);
}
