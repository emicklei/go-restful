# https://github.com/matm/gocov-html
# go get github.com/axw/gocov/gocov
# go get github.com/matm/gocov-html

gocov test github.com/emicklei/go-restful | gocov-html > restful.html && open restful.html