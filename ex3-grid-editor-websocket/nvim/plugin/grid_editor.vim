if exists('g:loaded_grid_editor')
  finish
endif
let g:loaded_grid_editor = 1
lua require('grid_editor').setup()

