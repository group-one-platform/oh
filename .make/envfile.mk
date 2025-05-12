ifneq ($(wildcard .env),)
  $(info $(YELLOW)including .env file$(RESET))
  include .env
endif

ifneq ($(wildcard .env.local),)
  $(info $(YELLOW)including .env.local file$(RESET))
  include .env.local
endif
