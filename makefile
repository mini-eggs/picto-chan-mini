default: dev

deploy: 
	rm gin-bin || true && \
	cat main.go | sed 's/package main/package handler/' | tee main.go >> /dev/null && \
	now && now alias chan.evanjon.es 

dev: 
	cat main.go | sed 's/package handler/package main/' | tee main.go >> /dev/null && \
	gin --all

