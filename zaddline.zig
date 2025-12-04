// addline.zig
// Добавляет строку в файл только если её там ещё нет (точное совпадение)

const std = @import("std");
const os = std.os;
const fs = std.fs;
const mem = std.mem;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    if (args.len != 3) {
        const progname = fs.path.basename(args[0]);
        std.debug.print("Использование: {s} <файл> \"строка\"\n", .{progname});
        std.debug.print("Пример: {s} /etc/environment \"MYVAR=hello\"\n", .{progname});
        return error.InvalidArgs;
    }

    const filename = args[1];
    const line_to_add = args[2];

    // Открываем файл с эксклюзивной блокировкой (чтобы не было гонки)
    const file = try fs.cwd().openFile(filename, .{ .mode = .read_write });
    defer file.close();

    try os.flock(file.handle, os.LOCK.EX);

    // Читаем весь файл в память (для простоты и скорости)
    const content = try file.readToEndAlloc(allocator, 1024 * 1024 * 10); // макс 10 МБ
    defer allocator.free(content);

    // Проверяем, есть ли точная строка (с \n или в конце файла)
    const needle = line_to_add ++ "\n";
    if (mem.indexOf(u8, content, needle)) |_| {
        std.debug.print("Строка уже есть в файле: {s}\n", .{filename});
        return;
    }
    // Также проверяем случай, если строка в конце файла без \n
    if (mem.endsWith(u8, content, line_to_add)) {
        std.debug.print("Строка уже есть в файле (в конце без \\n): {s}\n", .{filename});
        return;
    }

    // Перемещаем указатель в конец и дописываем
    try file.seekTo(try file.getEndPos());
    try file.writer().print("{s}\n", .{line_to_add});

    std.debug.print("Строка успешно добавлена: {s}\n", .{filename});
}
