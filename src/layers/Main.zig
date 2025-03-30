const std = @import("std");
const common = @import("common.zig");
const Ponds = @import("main/Ponds.zig");
const Quacks = @import("main/Quacks.zig");
const Insert = @import("main/Insert.zig");

ponds: Ponds,
quacks: Quacks,
insert: Insert,
alloc: std.mem.Allocator,


pub fn handle_current_state(mode: *common.MODE, key: u8) void {
    switch (mode) {
        .INSERT => handle_insert(key),
        .NORMAL => handle_normal(key),
        else => {},
    }
}

fn handle_insert(key: u8) void {
    switch (key) {}
}
fn handle_normal(key: u8) void {
    switch (key) {}
}
