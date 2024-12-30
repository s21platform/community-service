package postgres

import (
	"context"
)

func (r *Repository) AddPeerLogins(ctx context.Context, peerLogins []string) error {
	// в таблице participant есть foreign key на трайб и кампус что в них заносить? скорее всего нельзя будет добавить поля без добавления этих значений
	//или нужно сначала взять логины потом подтянуть остальную информацию с платформы и после уже добавлять в таблицу participant?
	return nil
}
