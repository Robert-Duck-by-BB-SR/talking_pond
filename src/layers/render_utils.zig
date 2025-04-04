const std = @import("std");
const common = @import("common.zig");
const border = common.theme.border;

pub fn generate_top_border_with_title(width: i16, title: []const u8, horizontal_border: *std.ArrayList(u8)) void {
    var j: usize = 0;
    horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    horizontal_border.appendSliceAssumeCapacity(title);
    while (j < width - @as(i16, @intCast(title.len)) - 3) {
        defer j += 1;
        horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    }
}

pub fn generate_bottom_border(width: i16, horizontal_border: *std.ArrayList(u8)) void {
    var j: usize = 0;
    horizontal_border.clearRetainingCapacity();
    horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    while (j < width - 3) {
        defer j += 1;
        horizontal_border.appendSliceAssumeCapacity(border.HORIZONTAL);
    }
}
