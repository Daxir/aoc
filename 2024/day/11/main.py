from functools import lru_cache

def read_input(file):
  with open(file, 'r') as f:
    first_line = f.readline().strip()
    return [int(x) for x in first_line.split(' ')]

# This single line trivializes the problem completely
@lru_cache(maxsize=None)
def count_stones(stone, blinks):
  if blinks == 0:
    return 1

  if stone == 0:
    return count_stones(1, blinks - 1)

  str_stone = str(stone)
  if len(str_stone) % 2 == 0:
    left = int(str_stone[:len(str_stone) // 2])
    right = int(str_stone[len(str_stone) // 2:])
    return count_stones(left, blinks - 1) + count_stones(right, blinks - 1)
  
  return count_stones(stone * 2024, blinks - 1)

if __name__ == '__main__':
  stones = read_input('input.txt')

  blinks = 25
  stone_count = 0
  for stone in stones:
    stone_count += count_stones(stone, blinks)

  print("(Part one) stone count: " + str(stone_count))

  blinks = 75
  stone_count = 0
  for stone in stones:
    stone_count += count_stones(stone, blinks)

  print("(Part two) stone count: " + str(stone_count))