package dialogs

import (
	"github.com/labstack/echo"
	json "github.com/pquerna/ffjson/ffjson"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type (
	// Questions является вебхук-каналом входящих запросов от пользователя.
	Questions <-chan Question

	// Answers является вебхук-каналом исходящих ответов к пользователям.
	Answers chan Answer
)

func init() {
	log.SetLevel(log.InfoLevel)
}

// New создаёт простой роутер для прослушивания входящих данных по вебхуку и
// возвращает два канала: для чтения запросов и отправки ответов соответственно.
// и функцию обработчик для передачи в роут
func New() (Questions, Answers, func(c echo.Context) error) {

	var err error
	questions := make(chan Question)
	answers := make(chan Answer)

	handleFunc := func(c echo.Context) error {
		log.Debugln("Тело входящего запроса:")
		log.Debugln(c.Request().Body)

		log.Debugln("Декодируем запрос...")
		var question Question
		if err := c.Bind(&question); err != nil {
			return err
		}
		log.Debugln("Отправляем запрос в канал...", questions)
		questions <- question

		var answer Answer
		for answer = range answers {
			a := answer.Session
			q := question.Session
			if !strings.EqualFold(a.SessionID, q.SessionID) ||
				!strings.EqualFold(a.UserID, q.UserID) ||
				a.MessageID != q.MessageID {
				log.Debugln("Это не тот ответ...")
				continue
			}

			log.Debugln("Обнаружен подходящий запрос! Отвечаем...")
			break
		}

		log.Debugln("Дождались нужный ответ! Отправляем его...")
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusOK)

		log.Debugln("Кодируем ответ...")
		if err = json.NewEncoder(c.Response()).Encode(answer); err != nil {
			log.Debugln("Ошибка:", err.Error())
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		return err
	}
	return questions, answers, handleFunc
}
