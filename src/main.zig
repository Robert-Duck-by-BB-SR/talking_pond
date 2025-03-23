const std = @import("std");
const shit_os = std.os.windows;
const posix = std.posix;
const os_tag = @import("builtin").os.tag;

const ENABLE_LINE_INPUT: u32 = 0x2;
const ENABLE_ECHO_INPUT: u32 = 0x4;
const ENABLE_PROCESSED_INPUT: u32 = 0x1;
const ENABLE_WINDOW_INPUT: u32 = 0x8;

fn start_raw_mode(std_in: std.fs.File, std_out: std.fs.File) !void {
    switch (os_tag) {
        .windows => {
            var old_stdin: shit_os.DWORD = undefined;
            _ = shit_os.kernel32.GetConsoleMode(std_in, &old_stdin);
            var raw_mode = old_stdin & ~(ENABLE_LINE_INPUT |
                ENABLE_ECHO_INPUT |
                ENABLE_PROCESSED_INPUT);

            raw_mode |= ENABLE_WINDOW_INPUT;

            _ = shit_os.kernel32.SetConsoleMode(std_in, raw_mode);

            var old_stdout: shit_os.DWORD = undefined;
            _ = shit_os.kernel32.GetConsoleMode(std_out, &old_stdout);
            _ = shit_os.kernel32.SetConsoleMode(std_out, old_stdout | shit_os.ENABLE_VIRTUAL_TERMINAL_PROCESSING);
        },
        .linux => {},
        else => return error.UNSUPPORTED_OS,
    }
}

fn restore_terminal(std_in: std.fs.File, old_stdin: ?shit_os.DWORD) void {
    switch (os_tag) {
        .windows => {
            _ = shit_os.kernel32.SetConsoleMode(std_in, old_stdin.?);
        },
        .linux => {},
        else => {},
    }
}

fn read_terminal(std_in: std.fs.File, stdout: std.fs.File.Writer) !bool {
    var buf = [1]u8{0};
    switch (os_tag) {
        .windows => _ = try shit_os.ReadFile(std_in, &buf, null),
        .linux => {},
        else => {},
    }
    switch (buf[0]) {
        3 => return true,
        '\n' => try stdout.print("\n", .{}),
        else => try stdout.print("\x1b[48;2;25;60;80m{c}\x1b[0m", .{buf[0]}),
    }
    return false;
}

pub const TerminalDimensions = struct { width: u16, height: u16 };
pub const UnixWinSize = struct { row: u16 = 0, col: u16 = 0, xpixel: u16 = 0, ypixel: u16 = 0 };

fn get_terminal_dimentions(std_out: std.fs.File) !TerminalDimensions {
    const terminal_dimensions: TerminalDimensions = undefined;
    switch (os_tag) {
        .windows => {
            terminal_dimensions = .{
                .width = shit_os.CONSOLE_SCREEN_BUFFER_INFO.dwSize.X,
                .height = shit_os.CONSOLE_SCREEN_BUFFER_INFO.dwSize.Y,
            };
        },
        .linux, .macos => {
            const win_size: UnixWinSize = .{};
            const res = std.os.linux.ioctl(@intCast(std_out.handle), std.os.linux.T.IOCGWINSZ, @intFromPtr(&win_size));
            if (res != 0) {
                return error.ioctl_return_error_during_getting_linux_dimentions;
            }
            terminal_dimensions = .{
                .width = win_size.col,
                .height = win_size.row,
            };
        },
        else => return error.UNSUPPORTED_OS,
    }
}

pub fn main() !void {
    const std_out = std.io.getStdOut();

    const std_in = std.io.getStdIn();
    const stdout = std_out.writer();

    get_terminal_dimentions(std_out);

    var exit = false;

    while (!exit) {
        exit = try read_terminal(std_in, stdout);
    }
}
