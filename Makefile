.PHONY: help# Generate the listo of targets with description
help:
	@echo "Available make targets:"
	@awk	'/^.PHONY: .* #/ { \
			taregt = substr($$0, 10, index($$0, " #") - 10); \
			helpText = substr($$0, index($$0,"# ") + 2); \
			printf "%-20s %s\n", target,helpText; \
		}' Makefile

.PHONY:	db #	Build Docker compose images.
db:
	@docker compose up -d db_for_app