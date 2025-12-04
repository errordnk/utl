// swap.zig
// Атомарный обмен любых двух объектов файловой системы (файлы, папки, ссылки и т.д.)
// Требуется Linux с ядром ≥ 3.15

const std = @import("std");
const os = std.os;
const linux = os.linux;
const builtin = @import("builtin");

// Константы, которых может не быть в старых версиях Zig
const AT_FDCWD = -100;
const RENAME_EXCHANGE = 1 << 2; // 0x4

fn syscall5(num: usize, arg1: usize, arg2: usize, arg3: usize, arg4: usize, arg5: usize) usize {
    return switch (builtin.cpu.arch) {
        .x86_64 => asm volatile ("syscall"
            : [ret] "={rax}" (-> usize),
            : [num] "{rax}" (num),
              [a1] "{rdi}" (arg1),
              [a2] "{rsi}" (arg2),
              [a3] "{rdx}" (arg3),
              [a4] "{r10}" (arg4),
              [a5] "{r8}" (arg5),
            : .{ .rcx = true, .r11 = true, .memory = true }),
        .aarch64 => asm volatile ("svc #0"
            : [ret] "={x0}" (-> usize),
            : [num] "{x8}" (num),
              [a1] "{x0}" (arg1),
              [a2] "{x1}" (arg2),
              [a3] "{x2}" (arg3),
              [a4] "{x3}" (arg4),
              [a5] "{x4}" (arg5),
            : .{ .memory = true }),
        else => @compileError("Архитектура не поддерживается"),
    };
}

fn renameat2_exchange(oldpath: []const u8, newpath: []const u8) !void {
    const rc = syscall5(
        linux.SYS.renameat2,
        @intCast(AT_FDCWD),
        @intFromPtr(oldpath.ptr),
        @intCast(AT_FDCWD),
        @intFromPtr(newpath.ptr),
        RENAME_EXCHANGE,
    );

    switch (@as(isize, @bitCast(rc))) {
        0 => return,
        -@as(isize, @intFromEnum(os.E.INVAL)), -@as(isize, @intFromEnum(os.E.NOSYS)) => return error.KernelTooOld,
        -@as(isize, @intFromEnum(os.E.XDEV)) => return error.CrossDevice,
        -@as(isize, @intFromEnum(os.E.NOENT)) => return error.NotFound,
        -@as(isize, @intFromEnum(os.E.ISDIR)), -@as(isize, @intFromEnum(os.E.NOTDIR)) => return error.TypeMismatch,
        else => return os.unexpectedErrno(@intCast(-rc)),
    }
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    if (args.len != 3) {
        std.debug.print("Использование: {s} <путь1> <путь2>\n", .{args[0]});
        std.debug.print("Атомарно меняет местами любые два объекта (файлы, папки, ссылки...)\n", .{});
        std.process.exit(1);
    }

    const a = args[1];
    const b = args[2];

    renameat2_exchange(a, b) catch |err| {
        switch (err) {
            error.KernelTooOld => {
                std.debug.print("Ошибка: ядро не поддерживает атомарный обмен (нужен Linux ≥ 3.15)\n", .{});
            },
            error.CrossDevice => {
                std.debug.print("Ошибка: объекты находятся на разных файловых системах\n", .{});
            },
            error.NotFound => {
                std.debug.print("Ошибка: один из путей не существует\n", .{});
            },
            error.TypeMismatch => {
                std.debug.print("Ошибка: нельзя обменять файл и директорию одновременно\n", .{});
            },
            else => {
                std.debug.print("Ошибка: {s}\n", .{@errorName(err)});
            },
        }
        std.process.exit(1);
    };

    std.debug.print("Успешно и атомарно обменяно:\n  {s}  ↔  {s}\n", .{ a, b });
}
