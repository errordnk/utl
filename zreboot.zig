const std = @import("std");

pub fn main() !void {
    const file = std.fs.openFileAbsolute("/proc/loadavg", .{}) catch return error.CannotOpenFile;
    defer file.close();

    var buf: [64]u8 = undefined;
    const bytes_read = file.read(&buf) catch return error.CannotReadFile;

    var it = std.mem.splitScalar(u8, buf[0..bytes_read], ' ');
    const la1_str = it.next() orelse return error.InvalidFormat;

    const la1 = std.fmt.parseFloat(f32, la1_str) catch return error.InvalidFloat;

    if (la1 > 600.0) {
        const argv = [_][]const u8{ "/sbin/reboot", "reboot" };
        return std.process.execve(std.heap.page_allocator, &argv, null);
    }
}
