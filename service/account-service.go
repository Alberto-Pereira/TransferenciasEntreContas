package service

import (
	"desafio-tecnico/model"
	"desafio-tecnico/repository"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func ReadAccounts() ([]model.Account, error) {

	accounts, err := repository.ReadAccounts()

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func ReadAccountBalance(accountId int64, accountSecret string) (int64, error) {

	accountBalance, err := repository.ReadAccountBalance(accountId, accountSecret)

	if err != nil {
		return 0, err
	}

	return accountBalance, nil
}

func CreateAccount(account model.Account) (int, error) {

	err := validateAccount(account)

	if err != nil {
		return 0, fmt.Errorf("Error validating account! %s", err.Error())
	}

	hashSecret, err := hashAccountSecret(account.Secret)

	if err != nil {
		return 0, err
	}

	id, err := repository.CreateAccount(model.Account{
		Name: account.Name, CPF: account.CPF, Secret: hashSecret,
		Balance: account.Balance, Created_at: account.Created_at})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func validateAccount(account model.Account) error {

	// ! Format - First letter uppercase, at least 2 characters and accepts space
	validName := regexp.MustCompile(`^([A-Z]{1}[a-z]+\s?)+$`)
	// ! Format - Accepts the following 000.000.000-00 pattern
	validCPF := regexp.MustCompile("^[0-9]{3}[.][0-9]{3}[.][0-9]{3}[-][0-9]{2}$")

	if !validName.MatchString(account.Name) {
		return errors.New("Invalid account name!")
	}

	if !validCPF.MatchString(account.CPF) {
		return errors.New("Invalid account cpf!")
	}

	if len(account.Secret) == 0 || account.Secret == " " {
		return errors.New("Invalid account secret!")
	}

	if account.Balance < 0 {
		return errors.New("Invalid account balance!")
	}

	if account.Created_at < 0 {
		return errors.New("Invalid account created time!")
	}

	return nil
}

func hashAccountSecret(accountSecret string) (string, error) {

	hashSecret, err := bcrypt.GenerateFromPassword([]byte(accountSecret), 14)

	if err != nil {
		return "", errors.New("Error while try to hash account secret!")
	}

	err = checkAccountSecretHash(accountSecret, string(hashSecret))

	if err != nil {
		return "", errors.New("Error while try to check if account secret and account secret hash matches!")
	}

	return string(hashSecret), nil
}

func checkAccountSecretHash(accountSecret string, hashSecret string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashSecret), []byte(accountSecret))

	return err
}
