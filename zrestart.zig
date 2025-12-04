const std = @import("std");

pub fn main() !void {
    const file = std.fs.openFileAbsolute("/etc/hostname", .{}) catch return error.CannotOpenFile;
    defer file.close();

    var buf: [64]u8 = undefined;
    const bytes_read = file.read(&buf) catch return error.CannotReadFile;

    // Parse hostname as integer
    const hostname_str = std.mem.trim(u8, buf[0..bytes_read], &std.ascii.whitespace);
    const hostnum = std.fmt.parseInt(i32, hostname_str, 10) catch return error.InvalidHostname;

    // Get current time using C time()
    const c = @cImport({
        @cInclude("time.h");
    });
    const now = c.time(null);
    const tm_ptr = c.localtime(&now);
    const current_hour: i32 = @intCast(tm_ptr.*.tm_hour);
    const current_min: i32 = @intCast(tm_ptr.*.tm_min);

    var should_restart = false;

    if (hostnum < 50) {
        // Even hours: 0,2,4,6,8,10,12,14,16,18,20,22
        // Minute = hostnum
        if (@mod(current_hour, 2) == 0 and current_min == hostnum) {
            should_restart = true;
        }
    } else {
        // Odd hours: 1,3,5,7,9,11,13,15,17,19,21,23
        // Minute = hostnum - 50
        const target_min = hostnum - 50;
        if (@mod(current_hour, 2) == 1 and current_min == target_min) {
            should_restart = true;
        }
    }

    if (should_restart) {
        const argv = [_][]const u8{ "/usr/bin/systemctl", "systemctl", "restart", "bot" };
        return std.process.execve(std.heap.page_allocator, &argv, null);
    }
}
