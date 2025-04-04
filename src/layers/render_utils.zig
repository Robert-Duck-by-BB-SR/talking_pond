const std = @import("std");
const common = @import("common.zig");
const border = common.theme.border;

pub fn generate_border_top(alloc: std.mem.Allocator, width: i16, horizontal_border: *std.ArrayList(u8)) ![]u8 {
    var j: usize = 0;
    while (j < width - 2) {
        defer j += 1;
        horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    }
    return try std.fmt.allocPrint(
        alloc,
        "{s}{s}{s}{s}",
        .{
            border.TOP_LEFT,
            horizontal_border.items,
            border.TOP_RIGHT,
            common.RESET_STYLES,
        },
    );
}

/// Border is generated in style:
/// --TEXT-----
/// That's why we append border.HORIZONTAL in the beninging
pub fn generate_border_top_with_title(alloc: std.mem.Allocator, width: i16, title: []const u8, horizontal_border: *std.ArrayList(u8)) ![]u8 {
    var j: usize = 0;
    horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    horizontal_border.appendSliceAssumeCapacity(title);
    while (j < width - @as(i16, @intCast(title.len)) - 4) {
        defer j += 1;
        horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    }
    return try std.fmt.allocPrint(
        alloc,
        "{s}{s}{s}{s}",
        .{
            border.TOP_LEFT,
            horizontal_border.items,
            border.TOP_RIGHT,
            common.RESET_STYLES,
        },
    );
}

pub fn generate_border_bottom(alloc: std.mem.Allocator, width: i16, horizontal_border: *std.ArrayList(u8)) ![]u8 {
    var j: usize = 0;
    horizontal_border.clearRetainingCapacity();
    while (j < width - 2) {
        defer j += 1;
        horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    }

    return try std.fmt.allocPrint(
        alloc,
        "{s}{s}{s}{s}",
        .{
            border.BOTTOM_LEFT,
            horizontal_border.items,
            border.BOTTOM_RIGHT,
            common.RESET_STYLES,
        },
    );
}
