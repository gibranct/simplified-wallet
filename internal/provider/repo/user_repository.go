package repository

import (
	"context"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"log"
	"strings"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/provider/db/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db   *sqlx.DB
	otel telemetry.Telemetry
}

var allUserColumns = []string{
	"id",
	"name",
	"email",
	"password",
	"balance",
	"cpf",
	"cnpj",
	"user_type",
	"active",
	"created_at",
	"updated_at",
}

func (ur UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := ur.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exists, nil
}

func (ur UserRepository) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE cpf = $1)"
	err := ur.db.GetContext(ctx, &exists, query, cpf)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exists, nil
}

func (ur UserRepository) ExistsByCNPJ(ctx context.Context, cnpj string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE cnpj = $1)"
	err := ur.db.GetContext(ctx, &exists, query, cnpj)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exists, nil
}

func (ur UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var user model.UserModel
	query := "SELECT " + strings.Join(allUserColumns, ", ") + " FROM users WHERE id = $1"
	err := ur.db.GetContext(
		ctx,
		&user,
		query, userID,
	)
	if err != nil {
		return nil, err
	}
	return user.ToEntity()
}

func (ur UserRepository) Save(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users 
	(id, name, email, password, balance, cpf, cnpj, user_type, active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING id`
	var userID uuid.UUID
	userModel := model.NewUserModelFrom(user)
	err := ur.db.GetContext(
		ctx,
		&userID,
		query,
		userModel.ID,
		userModel.Name,
		userModel.Email,
		userModel.Password,
		userModel.Balance,
		userModel.CPF,
		userModel.CNPJ,
		userModel.UserType,
		userModel.Active,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ur UserRepository) UpdateBalance(ctx context.Context, senderID, receiverID string, updateFn func(sender, receiver *entity.User) (*entity.Transaction, error)) error {
	return runInTx(ctx, ur.db, func(tx *sqlx.Tx) error {
		query1 := "SELECT " + strings.Join(allUserColumns, ", ") + " FROM users WHERE id = $1"
		var sender model.UserModel
		err := tx.GetContext(ctx, &sender, query1, senderID)
		if err != nil {
			log.Println(err)
			return errs.ErrSenderNotFound
		}

		var receiver model.UserModel
		err = tx.GetContext(ctx, &receiver, query1, receiverID)
		if err != nil {
			log.Println(err)
			return errs.ErrReceiverNotFound
		}

		senderEntity, err := sender.ToEntity()
		if err != nil {
			return err
		}
		receiverEntity, err := receiver.ToEntity()
		if err != nil {
			return err
		}
		transaction, err := updateFn(senderEntity, receiverEntity)

		if err != nil {
			return err
		}

		query2 := "UPDATE users SET balance = $1, updated_at = NOW() WHERE id = $2"
		_, err = tx.ExecContext(ctx, query2, senderEntity.Balance(), senderID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query2, receiverEntity.Balance(), receiverID)
		if err != nil {
			return err
		}

		insertQuery := "INSERT INTO transactions (id, sender_id, receiver_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)"
		_, err = tx.ExecContext(ctx, insertQuery, transaction.ID(), transaction.SenderID(), transaction.ReceiverID(), transaction.Amount(), transaction.CreatedAt())

		return err

	})
}

func NewUserRepository(db *sqlx.DB, otel telemetry.Telemetry) UserRepository {
	return UserRepository{db: db, otel: otel}
}
