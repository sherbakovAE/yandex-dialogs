package dialogs

import (
	"database/sql"
	"regexp"
	"strings"
)

//func GetFunctionName(i interface{}) string {
//	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
//}

type HandFunc func(question Question, answer *Answer, p *Pipeline) (finish bool, err error)
type Filter struct {
	State   State
	Command []string
	Pattern string
}
type PipeItem struct {
	Fn     HandFunc
	Filter Filter
}
type Pipeline struct {
	PipeOfFunc []PipeItem
	Storage    Storage
	Questions  Questions
	Answers    Answers
	DB         *sql.DB // база данных (при необходимости)
}

// регистрация функция в списке
func (p *Pipeline) Register(handFunc HandFunc, filter Filter) {
	p.PipeOfFunc = append(p.PipeOfFunc, PipeItem{Fn: handFunc, Filter: filter})
}

// запуск всех зарегистрированных функций, если текущее состояние и введённые команды пройдут их фильтры
func (p *Pipeline) Start() error {

	for question := range p.Questions {

		answer := NewAnswer(question, "")
		for _, item := range p.PipeOfFunc {
			// проверки фильтров
			// проверка на состояние
			if p.Storage.GetState(question.Session.UserID) != item.Filter.State && item.Filter.State != -1 {
				continue
			}
			//log.Println(question.Request.Command,GetFunctionName(item.Fn),item.Filter.State)

			// проверка на паттерн регулярного выражения
			if item.Filter.Pattern != "" {
				re := regexp.MustCompile(item.Filter.Pattern)
				if !re.MatchString(question.Request.Command) {
					//log.PrintLn("Результат выражения", re.MatchString(question.Request.Command), item.Filter.Pattern, question.Request.Command)
					continue
				}
			}

			// если в ответе пользователя есть хоть одно слово из списка , задача выполняется
			existsCommnd := false
			for i := 0; i < len(item.Filter.Command); i++ {
				for _, word := range strings.Split(question.Request.Command, " ") {
					if word == item.Filter.Command[i] {
						existsCommnd = true
					}
				}
			}

			// если нет фильтра по регулярному выражению , но есть фильтр по команде, но её нет - пропускаем
			if item.Filter.Pattern == "" && len(item.Filter.Command) != 0 && !existsCommnd {
				continue
			}
			// выполнение функции
			finish, err := item.Fn(question, &answer, p)
			// при ошибке закрываем сессию и удаляем из хранилища
			if err != nil {
				answer.Response.Text = "Ошибка сервера :" + err.Error()
				p.Storage.Delete(question.Session.UserID)
				answer.Response.EndSession = true
				break
			}
			if finish {
				break
			}
		}
		//log.Println(answer.Response.Text)
		p.Answers <- answer
	}

	return nil
}
