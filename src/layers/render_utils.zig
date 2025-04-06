const std = @import("std");
const common = @import("common.zig");
const border = common.theme.border;

pub fn render_border_top(alloc: std.mem.Allocator, width: i16, horizontal_border: *std.ArrayList(u8)) ![]u8 {
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
pub fn render_border_top_with_title(alloc: std.mem.Allocator, width: i16, title: []const u8, horizontal_border: *std.ArrayList(u8)) ![]u8 {
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

pub fn render_border_bottom(alloc: std.mem.Allocator, width: i16, horizontal_border: *std.ArrayList(u8)) ![]u8 {
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

pub fn render_line_of_text_and_backround(alloc: std.mem.Allocator, text: []const u8, width: usize) ![]u8 {
    var result: []u8 = undefined;
    if (text.len >= width) {
        result = try render_truncated_line_of_text_and_backround(alloc, text, width);
    } else {
        const bg_mid = try alloc.alloc(u8, width - text.len);
        @memset(bg_mid, ' ');
        result = try std.fmt.allocPrint(
            alloc,
            "{s}{s}",
            .{
                text,
                bg_mid,
            },
        );
    }

    return result;
}

fn render_truncated_line_of_text_and_backround(alloc: std.mem.Allocator, text: []const u8, width: usize) ![]u8 {
    const bg_len = 1;
    var truncated_text: []const u8 = undefined;
    const to_remove = text.len - width;
    // 3 for ... and 1 for bg
    truncated_text = text[0 .. text.len - to_remove - 4];
    const bg_mid = try alloc.alloc(u8, bg_len);
    @memset(bg_mid, ' ');
    const result = try std.fmt.allocPrint(
        alloc,
        "{s}{s}{s}",
        .{
            truncated_text,
            "...",
            bg_mid,
        },
    );

    return result;
}
