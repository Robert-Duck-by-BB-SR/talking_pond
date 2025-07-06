const std = @import("std");
const common = @import("../common.zig");
const RenderQ = @import("../../RenderQueue.zig");
const render_utils = @import("../render_utils.zig");
const formatting_utils = @import("../formatting_utils.zig");
const debug = @import("../../debug/debugger.zig");
const linked_quacks_list = @import("./LinkedQuacksList.zig").LinkedQuaksList;

dimensions: common.Dimensions = undefined,
position: common.Position = undefined,

main_allocator: std.mem.Allocator,
render_q: *RenderQ,

rows: []Row = undefined,
border: []u8 = undefined,
active_quack: usize = 0,
sliding_window_move_by: u8 = 0,
is_active: bool = false,

quacks_list: std.ArrayList(QuackItem) = undefined,

const QuackItem = struct {
    // id: []u8 = undefined,
    time: []const u8,
    author: []const u8,
    message: []const u8 = undefined,
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

    const quack_first: QuackItem = .{ .time = "69:42", .author = "bob", .message = "first" };
    const quack_last: QuackItem = .{ .time = "69:42", .author = "bob", .message = "last" };
    const quack_one: QuackItem = .{ .time = "69:42", .author = "bob", .message = "1" };
    const quack_two: QuackItem = .{ .time = "69:52", .author = "kakashi", .message = "2" };
    const quack_three: QuackItem = .{ .time = "69:69", .author = "bob", .message = "3" };
    // const quack_four: QuackItem = .{ .time = "69:69", .author = "kakashi", .message = "With a capital G" };
    const quack_five: QuackItem = .{ .time = "69:69", .author = "bibi", .message = "I start here, and is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including version" };

    try quacks_list.append(quack_first);
    try quacks_list.append(quack_one);
    try quacks_list.append(quack_two);
    try quacks_list.append(quack_three);
    try quacks_list.append(quack_one);
    try quacks_list.append(quack_one);
    try quacks_list.append(quack_two);
    try quacks_list.append(quack_three);
    try quacks_list.append(quack_two);
    try quacks_list.append(quack_three);
    try quacks_list.append(quack_five);
    try quacks_list.append(quack_five);
    try quacks_list.append(quack_last);

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
            .height = terminal_dimensions.height - 5,
        },
        .quacks_list = quacks_list,
    };
}

// salty: TODO: oh wait is this an abstraction???
// carrot: naah bro, trust me, one more abstraction and it will be definetely the best code
pub fn init_first_frame(self: *Self) !void {
    var arena = std.heap.ArenaAllocator.init(self.main_allocator);
    defer arena.deinit();
    const temp_allocator = arena.allocator();

    self.rows = try temp_allocator.alloc(Row, @intCast(self.dimensions.height));
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

    // Message list structure:
    // Suppose user got 6 messages (one, two, three, four, five, your mama)
    // And suppose that his screen is small as his schlong so it can fit only 4 messages
    // So the client will display only (three, four, five, your mama), where 'your mama' will be at the bottom
    // Structure of line: [12:00] author: message -> 7 (time), 1 (space), author.len + 1 (:) + 1(space) + message.len (slice)
    // Total = 11 + author.len + message.len (slice)

    const capacity: usize = @intCast(self.rows.len);
    var first_slice_index: usize = 0;
    if (capacity < self.quacks_list.items.len - 1) {
        first_slice_index = self.quacks_list.items.len - 1 - capacity;
    }

    const quacks_slice = self.quacks_list.items[first_slice_index..self.quacks_list.items.len];

    const max_width = @as(usize, @intCast(self.dimensions.width - 2));
    var lines_left_to_fill: usize = capacity;
    var linked_quacks = linked_quacks_list().new(temporary_alloctor);
    var i: usize = quacks_slice.len - 1;

    while (true) {
        if (lines_left_to_fill == 0) {
            break;
        }
        const prepared_message = try std.fmt.allocPrint(temporary_alloctor, " [{s}] {s}: {s}", .{
            quacks_slice[i].time,
            quacks_slice[i].author,
            quacks_slice[i].message,
        });

        // TODO: FIX LAST CHARACTER
        const all_content_lines = try render_utils.render_multiple_lines_with_background(
            temporary_alloctor,
            prepared_message,
            max_width,
            9,
        );

        var mult_lines_to_use: usize = 0;
        if (all_content_lines.len > lines_left_to_fill) {
            mult_lines_to_use = lines_left_to_fill;
            lines_left_to_fill = 0;
        } else {
            mult_lines_to_use = all_content_lines.len;
            lines_left_to_fill -= all_content_lines.len;
        }

        if (all_content_lines.len == 1) {
            try linked_quacks.reverse_insert(all_content_lines[0]);
        } else{
            // var lines_counter: usize = 0; 
            var index: usize = all_content_lines.len - 1;
            while (true) {
                const line = all_content_lines[index];
                try linked_quacks.reverse_insert(line);
                
                // if (index == 0 or lines_counter == mult_lines_to_use) {
                //     break;
                // }
                if (index == 0) {
                    break;
                }
                index -= 1;
                // lines_counter += 1;
            }
        }
        if (i == 0) {
            break;
        }

        i -= 1;
    }

    // IF SPACE LEFT - render from top to bottom
    // IF NO SPACE LEFT - render from bottom to top
    if (lines_left_to_fill == 0){
        var row_id: usize = self.rows.len - 1;
        var quack_node = linked_quacks.head;

        while (quack_node) |node| {
            @memcpy(self.rows[row_id].content[0..node.line.len], node.line);

            quack_node = node.prev;
            if (row_id == 0) {
                break;
            }
            row_id -= 1;
        }
    } else{
        var row_id: usize = 0;
        const max_row: usize = linked_quacks.get_len();
        var quack_node = linked_quacks.tail;

        while (quack_node) |node| {
            if (row_id == max_row or row_id == self.rows.len) {
                break;
            }
            @memcpy(self.rows[row_id].content[0..node.line.len], node.line);

            quack_node = node.next;
            row_id += 1;
        }
    }
}

fn render_row(self: *Self, temporary_alloctor: std.mem.Allocator, row_index: usize) ![]u8 {
    var render_result: std.ArrayList(u8) = .init(self.main_allocator);
    const row = self.rows[row_index];
    self.set_highlight_styles(row_index);
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
