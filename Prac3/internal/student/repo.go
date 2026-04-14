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
				FullName: "Кузнецов Михаил Столярович",
				Group:    "ПИМО-01-25",
				Email:    "mikhail.kuznetsov_87@inbox.ru",
			},
			2: {
				ID:       2,
				FullName: "Сидорова Зоя Николаевна",
				Group:    "ИВБО-02-25",
				Email:    "liliya_art_2024@mail.ru",
			},
			3: {
				ID:       3,
				FullName: "Иванов Иван Иванович",
				Group:    "ИВБО-03-25",
				Email:    "random.cat.lover@yandex.ru",
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
