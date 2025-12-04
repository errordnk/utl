// addline.c
// Добавляет строку в конец файла, только если её там ещё нет

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>

#define BUFFER_SIZE 4096

// Проверяет, есть ли строка (с завершающим \n) в файле
int line_exists(const char *filename, const char *line)
{
    FILE *fp = fopen(filename, "r");
    if (!fp)
        return 0; // файл не существует → строки точно нет

    size_t len = strlen(line);
    char buffer[BUFFER_SIZE];

    while (fgets(buffer, sizeof(buffer), fp))
    {
        // Убираем \n в конце буфера, если есть
        char *nl = strchr(buffer, '\n');
        if (nl)
            *nl = '\0';

        if (strcmp(buffer, line) == 0)
        {
            fclose(fp);
            return 1; // найдено точное совпадение
        }
    }

    fclose(fp);
    return 0;
}

// Добавляет строку + \n в конец файла
int append_line(const char *filename, const char *line)
{
    FILE *fp = fopen(filename, "a");
    if (!fp)
    {
        perror("fopen");
        return -1;
    }
    fprintf(fp, "%s\n", line);
    if (ferror(fp))
    {
        fclose(fp);
        return -1;
    }
    fclose(fp);
    return 0;
}

int main(int argc, char *argv[])
{
    if (argc != 3)
    {
        fprintf(stderr, "Использование: %s <файл> \"строка\"\n", argv[0]);
        fprintf(stderr, "Пример: %s /etc/fstab \"tmpfs /tmp tmpfs defaults,noatime 0 0\"\n", argv[0]);
        return 1;
    }

    const char *filename = argv[1];
    const char *line = argv[2];

    if (line_exists(filename, line))
    {
        printf("Строка уже присутствует в файле: %s\n", filename);
        return 0;
    }

    if (append_line(filename, line) == 0)
    {
        printf("Строка успешно добавлена в %s\n", filename);
        return 0;
    }
    else
    {
        fprintf(stderr, "Не удалось добавить строку в %s\n", filename);
        return 1;
    }
}