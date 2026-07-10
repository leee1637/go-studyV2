package service_registration

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"
	"study/internal/core/domain"
	email_service "study/internal/features/email/service"
	"study/internal/features/registration/repository"
	user_repository "study/internal/features/users/repository/postgres"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationService struct {
	regRepo      *repository.RegistrationRepository
	userRepo     *user_repository.UserRepository
	emailService *email_service.EmailService
}

func NewRegistrationService(regRepo *repository.RegistrationRepository, userRepo *user_repository.UserRepository, email *email_service.EmailService) *RegistrationService {
	return &RegistrationService{regRepo: regRepo, userRepo: userRepo, emailService: email}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func mapHeaders(header []string) map[string]int {
	headers := make(map[string]int)
	for i, h := range header {
		key := strings.ToLower(strings.TrimSpace(h))
		headers[key] = i
	}
	return headers
}

func getFieldByIndex(row []string, headers map[string]int, name string) string {
	idx, ok := headers[name]
	if !ok || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func (s *RegistrationService) ProcessCSV(ctx context.Context, file []byte) (*domain.UploadResult, error) {
	const MaxSize = 5 * 1024 * 1024

	if len(file) > MaxSize {
		return nil, fmt.Errorf("файл слишком большой: %d байт (макс. %d)", len(file), MaxSize)
	}

	reader := csv.NewReader(bytes.NewReader(file))
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV пустой или нет строк данных")
	}

	headers := mapHeaders(records[0])

	requiredHeaders := []string{"фио", "роль", "email"}
	for _, name := range requiredHeaders {
		if _, ok := headers[name]; !ok {
			return nil, fmt.Errorf("отсутствует обязательная колонка: '%s'", name)
		}
	}

	emailSet := make(map[string]bool)
	var validRequests []domain.RegistrationRequest
	var warnings []domain.WarningRow

	for i, record := range records[1:] {
		rowNum := i + 2
		req, warning := s.validateRow(record, rowNum, headers, emailSet)
		if warning != nil {
			warnings = append(warnings, *warning)
			continue
		}
		validRequests = append(validRequests, *req)
	}

	if len(validRequests) > 0 {
		emails := make([]string, len(validRequests))
		for i, req := range validRequests {
			emails[i] = req.Email
		}

		existing, err := s.regRepo.EmailExistsInPending(ctx, emails)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки email в БД: %w", err)
		}

		if len(existing) > 0 {
			existingSet := make(map[string]bool)
			for _, e := range existing {
				existingSet[e] = true
			}

			var filtered []domain.RegistrationRequest
			for _, req := range validRequests {
				if existingSet[req.Email] {
					warnings = append(warnings, domain.WarningRow{
						Row:    0,
						Reason: fmt.Sprintf("Email %s уже имеет активную заявку", req.Email),
					})
				} else {
					filtered = append(filtered, req)
				}
			}
			validRequests = filtered
		}
	}

	if len(validRequests) > 0 {
		err := s.regRepo.SaveBatch(ctx, validRequests)
		if err != nil {
			return nil, fmt.Errorf("ошибка сохранения заявок: %w", err)
		}

		for _, req := range validRequests {
			err := s.emailService.SendVerificationLink(req.Email, req.FIO, req.Token.String())
			if err != nil {
				fmt.Printf("ошибка отправки email для %s: %v\n", req.Email, err)
			}
		}
	}

	result := &domain.UploadResult{
		Total:     len(records) - 1,
		Processed: len(validRequests),
		Skipped:   len(warnings),
		Warnings:  warnings,
	}

	return result, nil
}

func (s *RegistrationService) validateRow(row []string, rowNum int, headers map[string]int, seenEmails map[string]bool) (*domain.RegistrationRequest, *domain.WarningRow) {
	if len(row) < len(headers) {
		return nil, &domain.WarningRow{Row: rowNum, Reason: "Строка короче заголовка"}
	}

	fio := getFieldByIndex(row, headers, "фио")
	role := strings.ToUpper(getFieldByIndex(row, headers, "роль"))
	email := strings.ToLower(getFieldByIndex(row, headers, "email"))
	group := getFieldByIndex(row, headers, "группа")
	phone := getFieldByIndex(row, headers, "телефон")

	if len(fio) < 3 {
		return nil, &domain.WarningRow{Row: rowNum, Reason: "ФИО должно содержать минимум 3 символа"}
	}

	if role != "STUDENT" && role != "TEACHER" {
		return nil, &domain.WarningRow{Row: rowNum, Reason: fmt.Sprintf("Неверная роль: '%s' (допустимы STUDENT/TEACHER)", role)}
	}

	if email == "" {
		return nil, &domain.WarningRow{Row: rowNum, Reason: "Email не может быть пустым"}
	}

	if !isValidEmail(email) {
		return nil, &domain.WarningRow{Row: rowNum, Reason: fmt.Sprintf("Невалидный email: '%s'", email)}
	}

	if seenEmails[email] {
		return nil, &domain.WarningRow{
			Row:    rowNum,
			Reason: fmt.Sprintf("Дубликат email в файле: %s", email),
		}
	}
	seenEmails[email] = true

	if role == "STUDENT" && group == "" {
		return nil, &domain.WarningRow{Row: rowNum, Reason: "Для студента группа обязательна"}
	}

	if role == "TEACHER" {
		group = ""
	}

	var phoneNumber *string
	if phone != "" {
		phoneNumber = &phone
	}

	req := &domain.RegistrationRequest{
		ID:          uuid.New(),
		FIO:         fio,
		Role:        domain.Role(role),
		PhoneNumber: phoneNumber,
		Email:       email,
		GroupName:   group,
		Token:       uuid.New(),
		Status:      "pending",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	return req, nil
}

func (s *RegistrationService) GetByToken(ctx context.Context, token uuid.UUID) (domain.RegistrationRequest, error) {
	return s.regRepo.GetByToken(ctx, token)
}

func (s *RegistrationService) DeleteByToken(ctx context.Context, token uuid.UUID) error {
	return s.regRepo.DeleteByToken(ctx, token)
}

func (s *RegistrationService) CompleteRegistration(ctx context.Context, token uuid.UUID, password string) error {
	req, err := s.regRepo.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("Заявка не найдена: %w", err)
	}

	if req.Status != "pending" {
		return fmt.Errorf("Заявка уже обработана (статус: %s)", req.Status)
	}

	if time.Now().After(req.ExpiresAt) {
		s.regRepo.DeleteByToken(ctx, token)
		return fmt.Errorf("Срок действия ссылки истёк")
	}

	if len(password) < 8 {
		return fmt.Errorf("Пароль не может быть меньше 8 символов")
	}

	if len(password) > 100 {
		return fmt.Errorf("Пароль не может быть больше 100 символов")
	}

	_, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return fmt.Errorf("Пользователь с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Ошибка хеширования пароля: %w", err)
	}

	tx, err := s.userRepo.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	userID, err := s.userRepo.SaveUser(ctx, tx, user)
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	switch req.Role {
	case domain.RoleStudent:
		student := &domain.Student{
			ID:          userID,
			FIO:         req.FIO,
			GroupName:   req.GroupName,
			PhoneNumber: req.PhoneNumber,
		}
		err = s.userRepo.SaveUserStudent(ctx, tx, student)
		if err != nil {
			return fmt.Errorf("ошибка создания профиля студента: %w", err)
		}

	case domain.RoleTeacher:
		teacher := &domain.Teacher{
			ID:          userID,
			FIO:         req.FIO,
			PhoneNumber: "",
		}
		err = s.userRepo.SaveUserTeacher(ctx, tx, teacher)
		if err != nil {
			return fmt.Errorf("ошибка создания профиля преподавателя: %w", err)
		}
	}

	err = s.regRepo.MarkCompleted(ctx, token)
	if err != nil {
		return fmt.Errorf("ошибка обновления заявки: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil
}
