const std = @import("std");
const common = @import("common.zig");

pub fn message_author_rizzling(temporary_alloctor: std.mem.Allocator, author: []const u8, message: []const u8) ![]const u8 {
    if (message[0] == ' ') {
        return message;
    }

    var splitted_message = std.mem.splitScalar(u8, message, ' ');
    var result: std.ArrayList(u8) = .init(temporary_alloctor);

    var index: u8 = 0;
    while (splitted_message.next()) |part| {
        switch (index) {
            // Time
            0 => {
                try result.writer().print("{s} ", .{part});
            },
            // Author
            1 => {
                // remove : from the author name
                if (std.mem.eql(u8, author, part[0 .. part.len - 1])) {
                    try result.writer().print("{s}{s}{s} ", .{ common.theme.ACTIVE_FONT_COLOR, part, common.INACTIVE_ITEM });
                } else {
                    try result.writer().print("{s} ", .{part});
                }
            },
            // Message
            else => {
                try result.writer().print("{s} ", .{part});
            },
        }
        index += 1;
    }

    return result.toOwnedSlice();
}
