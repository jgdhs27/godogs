import copy
from typing import List
from dog_breeds import *


class DoesNotFitError(Exception):
    pass


class FilledGrid:
    def __str__(self):
        s = ""
        for row in self.rows:
            s += "".join(row) + "\n"
        return s

    def __init__(self, input_string: str):
        self.rows = []
        strings = input_string.split("\n")[1:-1]
        for s in strings:
            self.rows.append(list(s))

    def set(self, x: int, y: int, char: str):
        self.rows[y][x] = char

    def is_square_occupied(self, x: int, y: int) -> bool:
        return self.rows[y][x] != " "

    def rotate(self):
        new_rows = []
        for i in range(len(self.rows[0])):
            new_rows.append([row[-i - 1] for row in self.rows])

        new_dog = FilledGrid("")
        new_dog.rows = new_rows
        return new_dog

    def size(self) -> int:
        return sum(sum(1 for val in row if val != " ") for row in self.rows)


class Board(FilledGrid):
    def insert(self, other, x: int, y: int, char_to_insert: str):
        new_state = copy.deepcopy(self)
        for y_offset, row in enumerate(other.rows):
            for x_offset, value in enumerate(row):
                if other.is_square_occupied(x_offset, y_offset):
                    if self.is_square_occupied(x + x_offset, y + y_offset):
                        raise DoesNotFitError()
                    new_state.set(x + x_offset, y + y_offset, char_to_insert)
        return new_state


class Dog:
    def __init__(self, input_string: str):
        self.orientations = []
        orientation = FilledGrid(input_string)
        for _ in range(4):
            self.orientations.append(orientation)
            orientation = orientation.rotate()

    def size(self) -> int:
        return sum(
            sum(1 for val in row if val != " ") for row in self.orientations[0].rows
        )


# ■
board = Board(
    """
■     ■■
       ■
        
        
        
        
■       
■■     ■
"""
)

dogs = [
    Dog(EIKA),
    Dog(WANNI),
    Dog(VIVI),
    Dog(BORKO),
    Dog(RUST),
    Dog(OAKLEY),
    Dog(MARBLE),
    Dog(MIKKO),
    Dog(BORKO),
    Dog(MUCKI),
]

dogs.sort(key=lambda dog: -dog.size())

board_size_x = len(board.rows[0])
board_size_y = len(board.rows)

dogs_total = sum([dog.size() for dog in dogs])
board_total = board_size_x * board_size_y - board.size()
print("Total dog size:", dogs_total)
print("Total board size:", board_total)
assert dogs_total == board_total


def solve(dogs: List[Dog], board: Board):
    if len(dogs) == 0:
        return board

    if len(dogs) < 3:
        print(board)

    dog = dogs[0]
    other_dogs = dogs[1:]
    for dog_orientation in dog.orientations:
        for x in range(board_size_x):
            for y in range(board_size_y):
                try:
                    new_board = board.insert(dog_orientation, x, y, str(len(dogs) - 1))
                    return solve(other_dogs, new_board)
                except (IndexError, DoesNotFitError):
                    continue
    raise DoesNotFitError()


solution = solve(dogs, board)
print(solution)
