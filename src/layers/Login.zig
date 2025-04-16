const std = @import("std");

const common = @import("common.zig");

const HostField = struct {
    input: std.ArrayList(u8),
    const placeholder = "Host";
};

const KeyField = struct {
    input: std.ArrayList(u8),
    const placeholder = "Key";
};

const ErrorLine = struct {
    error_field: std.ArrayList(u8),
};

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
