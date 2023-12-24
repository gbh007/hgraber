package agent

import (
	"app/internal/domain/agent"
	"context"
)

func (uc *UseCase) Books(ctx context.Context) []agent.BookToHandle {
	books, err := uc.agentAPI.UnprocessedBooks(ctx, booksLimit)
	if err != nil {
		uc.logger.Error(ctx, err)

		return nil
	}

	return books
}

func (uc *UseCase) BookHandle(ctx context.Context, book agent.BookToHandle) {
	err := uc.bookHandle(ctx, book)
	if err != nil {
		uc.logger.Error(ctx, err)
	}
}
