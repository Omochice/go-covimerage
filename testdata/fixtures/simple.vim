" Comment line
let g:count = 0
for i in range(5)
  let g:count += 1
endfor

function! SimpleFunc()
  echo "hello"
  return 42
endfunction
