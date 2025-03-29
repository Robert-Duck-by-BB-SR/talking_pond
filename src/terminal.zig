const std = @import("std");
const fs = std.fs;
pub const shit_os = std.os.windows;
pub const posix = std.posix;
const os_tag = @import("builtin").os.tag;

const ENABLE_LINE_INPUT: u32 = 0x2;
const ENABLE_ECHO_INPUT: u32 = 0x4;
const ENABLE_PROCESSED_INPUT: u32 = 0x1;
const ENABLE_WINDOW_INPUT: u32 = 0x8;

pub const OldState = union {
    pub const TerminalDimensions = struct { width: i16, height: i16 };
    win: struct {
        std_out: shit_os.DWORD,
        std_in: shit_os.DWORD,
    },
    posix: struct {
        std_in: posix.termios,
    },
};

// Get terminal old state before running the application
pub fn get_termos_with_tea() !OldState {
    return switch (os_tag) {
        .windows => OldState{ .win = .{ .std_in = undefined, .std_out = undefined } },
        .linux, .macos => OldState{ .posix = .{ .std_in = undefined } },
        else => {
            return error.UNSUPPORTED_OS;
        },
    };
}

pub fn start_raw_mode(std_in: fs.File, std_out: fs.File, termos: *OldState) !void {
    switch (os_tag) {
        .windows => {
            var old_stdin: shit_os.DWORD = undefined;
            _ = shit_os.kernel32.GetConsoleMode(std_in.handle, &old_stdin);
            var raw_mode = old_stdin & ~(ENABLE_LINE_INPUT |
                ENABLE_ECHO_INPUT |
                ENABLE_PROCESSED_INPUT);

            raw_mode |= ENABLE_WINDOW_INPUT;

            _ = shit_os.kernel32.SetConsoleMode(std_in.handle, raw_mode);

            _ = shit_os.kernel32.GetConsoleMode(std_out.handle, &termos.win.std_out);
            _ = shit_os.kernel32.SetConsoleMode(std_out.handle, termos.win.std_out | shit_os.ENABLE_VIRTUAL_TERMINAL_PROCESSING);
        },
        .linux => {
            termos.posix.std_in = try posix.tcgetattr(std_in.handle);
            var raw = termos.posix.std_in;
            raw.lflag.ECHO = false;
            raw.lflag.ICANON = false;
            raw.lflag.ISIG = false;
            raw.lflag.IEXTEN = false;
            raw.iflag.ICRNL = false;
            raw.iflag.IXON = false;
            raw.oflag.OPOST = false;
            try posix.tcsetattr(
                std_in.handle,
                .FLUSH,
                raw,
            );
        },
        else => return error.UNSUPPORTED_OS,
    }
}

pub fn restore_terminal(std_in: fs.File, std_out: fs.File, termos: OldState) void {
    switch (os_tag) {
        .windows => {
            _ = shit_os.kernel32.SetConsoleMode(std_in.handle, termos.win.std_in);
            _ = shit_os.kernel32.SetConsoleMode(std_out.handle, termos.win.std_out);
        },
        .linux => {
            posix.tcsetattr(
                std_in.handle,
                .FLUSH,
                termos.posix.std_in,
            ) catch unreachable;
        },
        else => {},
    }
}
