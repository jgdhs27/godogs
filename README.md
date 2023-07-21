Solver for the game [Dogs Organized Neatly](https://store.steampowered.com/app/1597730/Dogs_Organized_Neatly/) using a depth-first search.

# Benchmarks

All measurements were done on level 62, which is also the level checked in.

Description of versions:
1. Python: original python implementation
1. Go v1 (78f764cc): Direct port of python to Go
1. Go v2 (3717f006): Removed all board copying. Instead, modify a single
board in place. 

## Time to solve

The underlying search didn't change, so each version searches 36,620 nodes.

| Python | Go v1 | Go v2 |
|---|---|---|
| 294700ms | 4780ms | 86ms |

## Time to view all nodes

1,806,909 nodes

| Python | Go v1 | Go v2 |
|---|---|---|
| Long | 235200ms | 3760ms |
