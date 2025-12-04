#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(void) {
    FILE *f = fopen("/proc/loadavg", "r");
    if (!f) {
        perror("fopen");
        return 1;
    }

    float la1;
    if (fscanf(f, "%f", &la1) != 1) {
        fclose(f);
        return 1;
    }
    fclose(f);

    if (la1 > 600.0) {
        execl("/sbin/reboot", "reboot", NULL);
        perror("execl");
        return 1;
    }

    return 0;
}
