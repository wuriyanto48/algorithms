#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int min(int a, int b, int c);
int word_distance(char* a, char* b);

int main() {
    printf("%d \n", min(1, 0, 2));
    
    char* a = "wury";
    char* b = "lury";
    
    int r = word_distance(a, b);
    printf("%d \n", r);
    return 0;
}

int min(int a, int b, int c) 
{
    if (b < a)
        a = b;
    
    if (c < a)
        return c;
    
    return a;
}

int word_distance(char* a, char* b) 
{
    int a_len = strlen(a)+1;
    int b_len = strlen(b)+1;
    
    int d[a_len][b_len];
    
    for (int i = 0; i <= a_len-1; i++)
        d[i][0] = i;
     
    for (int j = 0; j <= b_len-1; j++)
        d[0][j] = j;
    
    int i = 0;
    int j = 0;
    
    for (i = 0; i < a_len-1; i++) {
        for (j = 0; j < b_len-1; j++) {
            if (a[i] == b[j]) {
                d[i+1][j+1] = d[(i+1)-1][(j+1)-1];
            } else {
                d[i+1][j+1] = min(d[(i+1)-1][j+1],
					d[i+1][(j+1)-1],
					d[(i+1)-1][(j+1)-1]) + 1;
            }
        }
    }
    
    return d[i][j];
}