.PHONY: backend-dev extension-dev build-backend build-extension

# Backend
backend-dev:
	cd backend && go run main.go

build-backend:
	cd backend && go build -o server main.go

# Extension
extension-install:
	cd extension && npm install

extension-dev:
	cd extension && npm run dev

build-extension:
	cd extension && npm run build
