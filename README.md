# Game of life

This is a project I did from following a tutorial. The link to that repo can me found here:
https://github.com/KyleBanks/conways-gol

I have made some changes mainly to the structure because I want to expand on the OpenGL parts a lot in another project.

## Install

You can download and build directly from the source: 

```sh
$ go get github.com/PetrusJPrinsloo/gameoflife
```

## Configure

To change how the application runs just edit the default.json file in the root of the project.

```json
{
    "width": 1000,
    "height": 1000,
    "rows": 80,
    "columns": 80,
    "threshold": 0.15,
    "fps": 10
}
```

* `"width": 1000` Width of the window.
* `"height": 1000` Height of the window.
* `"rows": 80` Rows of cells.
* `"columns": 80` Columns of cells.
* `"threshold": 0.15` Chance of cell being alive on start, so default is 15%.
* `"fps": 10` Frames Per Second, each frame is a new generation.
