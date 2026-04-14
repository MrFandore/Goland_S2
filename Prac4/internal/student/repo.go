package student

import "errors"

var ErrStudentNotFound = errors.New("student not found")

type Repo struct {
	data map[int64]Student
}

func NewRepo() *Repo {
	return &Repo{
		data: map[int64]Student{
			1: {
				ID:       1,
				FullName: "Булискерия Будимир Бакуриевич",
				Group:    "ПИМО-01-25",
				Email:    "budimir_buliskeria@mail.ru",
			},
			2: {
				ID:       2,
				FullName: "Петрова Мария Сергеевна",
				Group:    "ИВБО-02-25",
				Email:    "petrova@example.com",
			},
			3: {
				ID:       3,
				FullName: "Сидоров Алексей Андреевич",
				Group:    "ИВБО-03-25",
				Email:    "sidorov@example.com",
			},
		},
	}
}

func (r *Repo) GetByID(id int64) (Student, error) {
	st, ok := r.data[id]
	if !ok {
		return Student{}, ErrStudentNotFound
	}
	return st, nil
}
