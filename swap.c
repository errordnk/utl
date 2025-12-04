#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <sys/syscall.h>

// Константы, которые могут отсутствовать в старых glibc
#ifndef AT_FDCWD
#define AT_FDCWD -100
#endif
#ifndef RENAME_EXCHANGE
#define RENAME_EXCHANGE (1 << 2) // 0x4
#endif

// Прямая обёртка через syscall — работает везде
static int exchange_paths(const char *a, const char *b)
{
    return syscall(__NR_renameat2,
                   AT_FDCWD, a,      // старый путь
                   AT_FDCWD, b,      // новый путь
                   RENAME_EXCHANGE); // флаг обмена
}

int main(int argc, char **argv)
{
    if (argc != 3)
    {
        fprintf(stderr, "Использование: %s <путь1> <путь2>\n", argv[0]);
        fprintf(stderr, "Атомарно меняет местами любые два объекта (файлы, папки, ссылки...)\n");
        return 1;
    }

    const char *src = argv[1];
    const char *dst = argv[2];

    if (exchange_paths(src, dst) == 0)
    {
        // printf("Успешно и атомарно обменяно:\n  %s  ↔  %s\n", src, dst);
        return 0;
    }

    // Человекочитаемые ошибки
    switch (errno)
    {
    case ENOSYS:
    case EINVAL:
        fprintf(stderr, "Ошибка: ядро не поддерживает атомарный обмен (нужен Linux ≥ 3.15)\n");
        break;
    case EXDEV:
        fprintf(stderr, "Ошибка: пути находятся на разных файловых системах\n");
        break;
    case ENOENT:
        fprintf(stderr, "Ошибка: один из путей не существует\n");
        break;
    case EISDIR:
    case ENOTDIR:
        fprintf(stderr, "Ошибка: нельзя обменять файл и директорию одновременно\n");
        break;
    default:
        perror("exchange");
    }
    return 1;
}