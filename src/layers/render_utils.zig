const std = @import("std");
const common = @import("common.zig");
const border = common.theme.BORDER;

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
/// |-- TEXT -----|
pub fn make_border_with_title(alloc: std.mem.Allocator, w: usize, title: []const u8) ![]u8 {
    // 4 for two corner characters and 2 spaces around the title
    const num_of_horizontal_total = w - title.len - 4;
    // 2 for spaces around the title
    const width = num_of_horizontal_total * border.HORIZONTAL.len + title.len + border.TOP_LEFT.len + border.TOP_RIGHT.len + 2;
    const horizontal_border = try alloc.alloc(u8, width);

    // first get all the points for memcopies
    const tlborder_end = border.TOP_LEFT.len;
    const hborder_end = tlborder_end + border.HORIZONTAL.len * 2 + 1; // <- 1 space character to outline the title
    const title_end = hborder_end + title.len + 1; // <- ditto
    const hbwidth = (width - title_end - border.TOP_RIGHT.len);

    @memcpy(horizontal_border[0..tlborder_end], border.TOP_LEFT);
    @memcpy(horizontal_border[tlborder_end .. hborder_end - 1], border.HORIZONTAL ** 2);
    horizontal_border[hborder_end - 1] = ' ';
    @memcpy(horizontal_border[hborder_end .. title_end - 1], title);
    horizontal_border[title_end - 1] = ' ';
    var i = title_end;
    while (i < title_end + hbwidth) {
        defer i += border.HORIZONTAL.len;
        @memcpy(horizontal_border[i .. i + border.HORIZONTAL.len], border.HORIZONTAL);
    }
    @memcpy(horizontal_border[i..], border.TOP_RIGHT);
    return horizontal_border;
}

pub fn make_bottom_border(alloc: std.mem.Allocator, width: usize) ![]u8 {
    const horizontal_border = try alloc.alloc(u8, width);
    const horizonal_border_len = (width - border.BOTTOM_LEFT.len - border.BOTTOM_RIGHT.len);
    const blborder_end = border.BOTTOM_LEFT.len;
    @memcpy(horizontal_border[0..blborder_end], border.BOTTOM_LEFT);
    var i = blborder_end;
    while (i <= horizonal_border_len) {
        defer i += border.HORIZONTAL.len;
        @memcpy(horizontal_border[i .. i + border.HORIZONTAL.len], border.HORIZONTAL);
    }
    @memcpy(horizontal_border[i..], border.BOTTOM_RIGHT);
    return horizontal_border;
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

pub fn render_border(alloc: std.mem.Allocator, is_active: bool, border_slice: []u8) ![]u8 {
    var render_result: std.ArrayList(u8) = .init(alloc);
    try render_result.writer().print("{s}{s}", .{
        if (is_active) .ACTIVE_BORDER else .INACTIVE_ITEM,
        border_slice,
    });
    return render_result.toOwnedSlice();
}
