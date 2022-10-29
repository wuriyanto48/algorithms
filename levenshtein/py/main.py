def word_distance(t1: str, t2: str, in_case_sensitive = False) -> int:
    t1_len = len(t1) + 1
    t2_len = len(t2) + 1

    if in_case_sensitive:
        t1 = t1.lower()
        t2 = t2.lower()

    def min(a: int, b: int, c: int) -> int:
        if b < a: a = b
        if c < a: return c
        return a
    
    d = [[0]*t1_len for _ in range(t2_len)]
    for i in range(t1_len):
        d[0][i] = i
    
    for i in range(t2_len):
        d[i][0] = i

    for y in range(t2_len - 1):
        for x in range(t1_len - 1):
            if t2[y] == t1[x]:
                d[y + 1][x + 1] = d[(y + 1) - 1][(x + 1) - 1]
            else:
                d[y + 1][x + 1] = min(d[(y + 1) - 1][x + 1], 
                    d[y + 1][(x + 1) - 1], 
                    d[(y + 1) -1][(x + 1) - 1]) + 1
    return d[t2_len-1][t1_len-1]

if __name__ == '__main__':
    a = 'baik'
    b = 'baik'

    r = word_distance(a, b, in_case_sensitive = True)
    print(r)