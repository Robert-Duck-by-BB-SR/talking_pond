const std = @import("std");
const common = @import("../common.zig");
const RenderQ = @import("../../RenderQueue.zig");
const render_utils = @import("../render_utils.zig");
const formatting_utils = @import("../formatting_utils.zig");
const debug = @import("../../debug/debugger.zig");

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

main_allocator: std.mem.Allocator,
render_q: *RenderQ,

rows: []Row = undefined,
border: []u8 = undefined,
active_quack: usize = 0,
is_active: bool = false,

quacks_list: std.ArrayList(QuackItem) = undefined,
visible_quacks_list: std.ArrayList(VisibleQuackItem) = undefined,

const QuackItem = struct {
    // id: []u8 = undefined,
    time: []const u8,
    author: []const u8,
    message: []const u8 = undefined,
};

const VisibleQuackItem = struct {
    id: usize,
    lines_will_take: u8,
};

const Row = struct {
    cursor: []u8 = undefined,
    content: []u8 = undefined,
};

const Self = @This();

pub fn create(alloc: std.mem.Allocator, terminal_dimensions: common.Dimensions, render_q: *RenderQ) !Self {
    var quacks_list: std.ArrayList(QuackItem) = try .initCapacity(
        alloc,
        // 8 = 1 (status line) + 2 (top and bottom border of input field) + 5 (lines for actual input)
        @intCast(terminal_dimensions.height - 8),
    );

    const quack_one: QuackItem = .{ .time = "69:42", .author = "bob", .message = "Yo bitches, whatup" };
    // const quack_two: QuackItem = .{ .time = "69:52", .author = "kakashi", .message = "Whatup homie" };
    // const quack_three: QuackItem = .{ .time = "69:69", .author = "bob", .message = "Bro what the fuck are you talking about? Btw, babagi?" };
    // const quack_four: QuackItem = .{ .time = "69:69", .author = "kakashi", .message = "With a capital G" };
    const quack_five: QuackItem = .{ .time = "69:69", .author = "bibi", .message = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum." };

    for (0..@intCast(terminal_dimensions.height - 10)) |i| {
        _ = i;
        try quacks_list.append(quack_one);
    }
    try quacks_list.append(quack_five);

    return Self{
        .render_q = render_q,
        .main_allocator = alloc,
        .position = .{
            .row = 1,
            .col = common.PONDS_SIDEBAR_SIZE + 1,
        },
        .dimensions = .{
            .width = terminal_dimensions.width - common.PONDS_SIDEBAR_SIZE - 1,
            // 6 = 1 (status line) + 5 (lines for actual input)
            .height = terminal_dimensions.height - 6,
        },
        .quacks_list = quacks_list,
        .visible_quacks_list = try .initCapacity(
            alloc,
            // 6 = 1 (status line) + 5 (lines for actual input)
            @intCast(terminal_dimensions.height - 6),
        ),
    };
}

// salty: TODO: oh wait is this an abstraction???
// carrot: naah bro, trust me, one more abstraction and it will be definetely the best code
pub fn init_first_frame(self: *Self) !void {
    var arena = std.heap.ArenaAllocator.init(self.main_allocator);
    defer arena.deinit();
    const temp_allocator = arena.allocator();

    self.rows = try temp_allocator.alloc(Row, @intCast(self.dimensions.height - 2));
    try self.render_border_with_title("QUACKS", temp_allocator);
    // Background
    for (self.rows, 2..) |*row, i| {
        const bg_mid = try self.main_allocator.alloc(u8, @intCast(self.dimensions.width - 2));
        @memset(bg_mid, ' ');
        row.cursor = try std.fmt.allocPrint(
            temp_allocator,
            common.MOVE_CURSOR_TO_POSITION,
            .{ i, self.position.col + 1 },
        );
        row.content = bg_mid;
    }
}

pub fn render_border_with_title(self: *Self, title: []const u8, temp_allocator: std.mem.Allocator) !void {
    const width: usize = @intCast(self.dimensions.width - 2);
    const corners_width = common.theme.BORDER.BOTTOM_LEFT.len + common.theme.BORDER.BOTTOM_RIGHT.len;
    const border_width = width * common.theme.BORDER.HORIZONTAL.len + corners_width;

    // Top border
    const top_border = try render_utils.make_border_with_title(
        temp_allocator,
        @intCast(self.dimensions.width),
        title,
    );
    self.border = try std.fmt.allocPrint(self.main_allocator, "{s}{s}", .{
        try std.fmt.allocPrint(
            temp_allocator,
            common.MOVE_CURSOR_TO_POSITION,
            .{ 1, self.position.col },
        ),
        top_border,
    });

    for (1..@intCast(self.dimensions.height - 1)) |i| {
        self.border = try std.fmt.allocPrint(temp_allocator, "{s}{s}{s}{s}{s}", .{
            self.border,
            try std.fmt.allocPrint(
                temp_allocator,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    i + 1,
                    self.position.col,
                },
            ),
            common.theme.BORDER.VERTICAL,
            try std.fmt.allocPrint(
                temp_allocator,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    i + 1,
                    self.position.col + self.dimensions.width - 1,
                },
            ),
            common.theme.BORDER.VERTICAL,
        });
    }

    // Bottom border
    const bottom_border = try render_utils.make_bottom_border(
        temp_allocator,
        border_width,
    );
    self.border = try std.fmt.allocPrint(
        temp_allocator,
        "{s}{s}{s}{s}",
        .{
            self.border,
            try std.fmt.allocPrint(
                temp_allocator,
                common.MOVE_CURSOR_TO_POSITION,
                .{
                    self.dimensions.height,
                    self.position.col,
                },
            ),
            bottom_border,
            common.RESET_STYLES,
        },
    );
}

pub fn fill_content_with_quacks(self: *Self, temporary_alloctor: std.mem.Allocator) !void {
    if (self.quacks_list.items.len == 0) {
        const content = try render_utils.render_line_of_text_and_backround(
            temporary_alloctor,
            "**DEAD SILENCE**",
            common.TEXT_POSITION.CENTER,
            @intCast(self.dimensions.width - 2),
        );
        const middle_index: usize = @intFromFloat(@as(f16, @floatFromInt(self.dimensions.height)) * 0.5);
        @memcpy(self.rows[middle_index - 2].content[0..content.len], content);
        return;
    }

    // Structure of line: [12:00] author: message -> 7 (time), 1 (space), author.len + 1 (:) + 1(space) + message.len (slice)
    // Total = 11 + author.len + message.len (slice)
    var lines_taken: u8 = 0;
    // used to move other messages down if there is a multiple line message
    var multiple_lines_spacing: u8 = 0;
    for (self.quacks_list.items[0..], 0..) |quack, quack_id| {
        const all_content_lines = try render_utils.render_multiple_lines_with_background(
            temporary_alloctor,
            try std.fmt.allocPrint(temporary_alloctor, " [{s}] {s}: {s}", .{
                quack.time,
                quack.author,
                quack.message,
            }),
            @intCast(self.dimensions.width - 2),
            9,
        );

        if (all_content_lines.items.len > 1) {
            var lines_counter: u8 = 0;
            for (all_content_lines.items) |value| {
                // NOTE: WHAT IF HALF OF THE MESSAGE WILL BE THERE?
                if (lines_taken == self.dimensions.height - 2) {
                    break;
                }
                @memcpy(self.rows[quack_id + multiple_lines_spacing].content[0..value.len], value);
                // last line for multiple line message does not need a spacing
                if (multiple_lines_spacing != all_content_lines.items.len - 1) {
                    multiple_lines_spacing += 1;
                }
                lines_counter += 1;
                lines_taken += 1;
            }
            try self.visible_quacks_list.append(VisibleQuackItem{
                // id which will be corresponding to id in quacks_list
                .id = quack_id,
                .lines_will_take = lines_counter,
            });
        } else {
            if (lines_taken == self.dimensions.height - 2) {
                break;
            }
            @memcpy(self.rows[quack_id + multiple_lines_spacing].content[0..all_content_lines.items[0].len], all_content_lines.items[0]);
            try self.visible_quacks_list.append(VisibleQuackItem{
                .id = quack_id,
                .lines_will_take = 1,
            });
            lines_taken += 1;
        }
    }
}

fn render_row(self: *Self, temporary_alloctor: std.mem.Allocator, row_index: usize) ![]u8 {
    //-----
    // GOOFING AROUND

    self.set_highlight_styles(row_index);

    //-----
    var render_result: std.ArrayList(u8) = .init(self.main_allocator);
    const row = self.rows[row_index];
    try render_result.writer().print("{s}{s}{s}", .{
        row.cursor,
        common.INACTIVE_ITEM,
        try formatting_utils.message_author_rizzling(temporary_alloctor, "bob", row.content),
    });
    return render_result.toOwnedSlice();
}

fn set_highlight_styles(self: *Self, content_row_index: usize) void {
    if (self.active_quack == content_row_index) {
        self.rows[content_row_index].content[0] = '>';
    } else {
        self.rows[content_row_index].content[0] = ' ';
    }
}

pub fn render(self: *Self) !void {
    var arena = std.heap.ArenaAllocator.init(self.main_allocator);
    defer arena.deinit();
    const temp_allocator = arena.allocator();
    var render_result: std.ArrayList(u8) = .init(self.main_allocator);
    try self.fill_content_with_quacks(self.main_allocator);
    for (0..self.rows.len) |i| {
        try render_result.writer().print("{s}", .{
            try self.render_row(temp_allocator, i),
        });
    }
    const rendered_border = try render_utils.rerender_border(self.main_allocator, self.is_active, self.border);
    try render_result.writer().print("{s}", .{rendered_border});
    const slice = try render_result.toOwnedSlice();
    try self.render_q.add_to_render_q(slice, .CONTENT);
    self.render_q.sudo_render();
}

pub fn handle_normal(self: *Self, mode: *common.MODE, key: u8, new_active: *common.ComponentType) !void {
    switch (key) {
        'j' => {
            if (self.active_quack == self.dimensions.height - 3) {
                return;
            }
            var arena = std.heap.ArenaAllocator.init(self.main_allocator);
            defer arena.deinit();
            const allocator = arena.allocator();

            const prev_quack = self.active_quack;
            self.active_quack += 1;

            const result = try std.fmt.allocPrint(allocator, "{s}{s}", .{
                try self.render_row(allocator, prev_quack),
                try self.render_row(allocator, self.active_quack),
            });

            try self.render_q.add_to_render_q(
                result,
                .CONTENT,
            );
            self.render_q.sudo_render();
        },
        'k' => {
            if (self.active_quack == 0) {
                return;
            }
            var arena = std.heap.ArenaAllocator.init(self.main_allocator);
            defer arena.deinit();
            const allocator = arena.allocator();

            const prev_quack = self.active_quack;
            self.active_quack -= 1;

            const result = try std.fmt.allocPrint(allocator, "{s}{s}", .{
                try self.render_row(allocator, prev_quack),
                try self.render_row(allocator, self.active_quack),
            });

            try self.render_q.add_to_render_q(
                result,
                .CONTENT,
            );
            self.render_q.sudo_render();
        },
        'C', 'P' => {
            new_active.* = .PONDS_SIDEBAR;
        },
        'I' => {
            new_active.* = .INPUT_FIELD;
        },
        ':' => {
            mode.* = .COMMAND;
        },
        else => {},
    }
}
