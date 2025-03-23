const std = @import("std");
const windows = std.os.windows;
const os_tag = @import("builtin").os.tag;

const ENABLE_LINE_INPUT: u32 = 0x2;
const ENABLE_ECHO_INPUT: u32 = 0x4;
const ENABLE_PROCESSED_INPUT: u32 = 0x1;
const ENABLE_WINDOW_INPUT: u32 = 0x8;

fn start_raw_mode(std_in: std.io.File, std_out: std.io.File) !void {
    switch (os_tag) {
        .windows => {
            var old_stdin: windows.DWORD = undefined;
            _ = windows.kernel32.GetConsoleMode(std_in, &old_stdin);
            var raw_mode = old_stdin & ~(ENABLE_LINE_INPUT |
                ENABLE_ECHO_INPUT |
                ENABLE_PROCESSED_INPUT);

            raw_mode |= ENABLE_WINDOW_INPUT;

            _ = windows.kernel32.SetConsoleMode(std_in, raw_mode);

            var old_stdout: windows.DWORD = undefined;
            _ = windows.kernel32.GetConsoleMode(std_out, &old_stdout);
            _ = windows.kernel32.SetConsoleMode(std_out, old_stdout | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING);
        },
        .linux => {},
        else => return error.UNSUPPORTED_OS,
    }
}

fn restore_terminal(std_in: std.io.File, std_out: std.io.File) void {
    switch (os_tag) {
        .windows => {
            _ = windows.kernel32.SetConsoleMode(std_in, old_stdin);
        },
        .linux => {},
        else => {},
    }
}

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const std_in = std.io.getStdIn();

    if (@import("builtin").os.tag == .windows) {

        // var console_info: windows.CONSOLE_SCREEN_BUFFER_INFO = undefined;

        // _ = windows.kernel32.GetConsoleScreenBufferInfo(std_out, &console_info);
        // try stdout.print("{}\n", .{console_info});

        var exit = false;

        while (!exit) {
            var buf = [1]u8{0};
            _ = try windows.ReadFile(std_in, &buf, null);

            switch (buf[0]) {
                3 => exit = true,
                '\n' => try stdout.print("\n", .{}),
                else => try stdout.print("\x1b[48;2;25;60;80m{c}\x1b[0m", .{buf[0]}),
            }
        }
    } else {}
}
