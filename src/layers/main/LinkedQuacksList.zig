const std = @import("std");

pub fn LinkedQuaksList() type {
    return struct {
        const Self = @This();

        const Node = struct {
            // quack line
            id: usize,
            line: []const u8,
            prev: ?*Node = null,
            next: ?*Node = null,
        };

        head: ?*Node,
        tail: ?*Node,
        length: u8,
        allocator: std.mem.Allocator,

        pub fn new(allocator: std.mem.Allocator) Self {
            return .{ .length = 0, .head = null, .tail = null, .allocator = allocator };
        }

        pub fn reverse_insert(self: *Self, value: []const u8) !void {
            var node = try self.allocator.create(Node);
            node.line = value;
            if (self.head == null) {
                self.head = node;
                self.tail = node;
                self.head.?.prev = self.tail;
                self.tail.?.next = self.head;
            } else {
                self.tail.?.prev = node;
                node.next = self.tail;
                self.tail = node;
            }
            self.length += 1;
        }

        pub fn get_len(self: *Self) u8 {
            return self.length;
        }
    };
}
