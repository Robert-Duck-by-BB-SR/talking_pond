const std = @import("std");
const builtin = @import("builtin");

pub fn build(b: *std.Build) void {
    const os_option = b.option([]const u8, "os", "os choice") orelse "";
    const available_os = enum { windows, mac, linux, unavailable };
    const os_choice = std.meta.stringToEnum(available_os, os_option) orelse .unavailable;
    const target = switch (os_choice) {
        .mac => std.Target.Os.Tag.macos,
        .linux => std.Target.Os.Tag.linux,
        .windows => std.Target.Os.Tag.windows,
        else => builtin.os.tag,
    };
    const exe_mod = b.createModule(.{
        .root_source_file = b.path("src/main.zig"),
        .target = b.standardTargetOptions(.{ .default_target = .{
            .cpu_arch = .x86_64,
            .os_tag = target,
        } }),
        .optimize = .Debug,
    });
    const exe = b.addExecutable(.{
        .name = "tp_zig",
        .root_module = exe_mod,
    });

    b.installArtifact(exe);

    const run_exe = b.addRunArtifact(exe);

    const run_step = b.step("run", "Run the application");
    run_step.dependOn(&run_exe.step);
}
