const std = @import("std");

pub fn debug_in_file_or_i_will_smack_your_face(message: []const u8) !void {
    const file = try std.fs.cwd().createFile(
        "junk_file.txt",
        .{ .read = true },
    );
    defer file.close();
    try file.writeAll(message);
}
