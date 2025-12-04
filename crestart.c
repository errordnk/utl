#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

int main(void) {
    FILE *f = fopen("/etc/hostname", "r");
    if (!f) {
        perror("fopen");
        return 1;
    }

    int hostnum;
    if (fscanf(f, "%d", &hostnum) != 1) {
        fclose(f);
        return 1;
    }
    fclose(f);

    // Get current time
    time_t now = time(NULL);
    struct tm *tm_info = localtime(&now);
    if (!tm_info) {
        return 1;
    }

    int current_hour = tm_info->tm_hour;
    int current_min = tm_info->tm_min;

    int should_restart = 0;

    if (hostnum < 50) {
        // Even hours: 0,2,4,6,8,10,12,14,16,18,20,22
        // Minute = hostnum
        if (current_hour % 2 == 0 && current_min == hostnum) {
            should_restart = 1;
        }
    } else {
        // Odd hours: 1,3,5,7,9,11,13,15,17,19,21,23
        // Minute = hostnum - 50
        int target_min = hostnum - 50;
        if (current_hour % 2 == 1 && current_min == target_min) {
            should_restart = 1;
        }
    }

    if (should_restart) {
        execl("/usr/bin/systemctl", "systemctl", "restart", "bot", NULL);
        perror("execl");
        return 1;
    }

    return 0;
}
