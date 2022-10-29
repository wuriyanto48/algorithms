# Wuriyanto 2022
# is_valid_cc
def is_valid_cc(cc: str) -> bool:
    def is_digit(d: int) -> bool:
        if (d < 48) or (d > 57):
            return False
        return True
    len_of_cc = len(cc)
    digit_checker = int(cc[len_of_cc-1:])

    parity = (len_of_cc-2)%2
    sum = 0

    for i in range(0, len_of_cc-1):
        if not is_digit(ord(cc[i])):
            return False
        digit = int(cc[i])
        if (i & 1) == parity:
            digit = digit * 2
        if digit > 9:
            digit = digit - 9
        sum = sum + digit
        
    return (10 - (sum % 10 )) == digit_checker

if __name__ == '__main__':
    print(is_valid_cc('79927398713')) # should return True
    print(is_valid_cc('5101865470958946')) # should return True
    print(is_valid_cc('5351851532935566')) # should return True
    print(is_valid_cc('5476362313068858')) # should return True
    print(is_valid_cc('0812238448827277')) # should return False