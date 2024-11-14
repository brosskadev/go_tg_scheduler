package scraper

import (
	"context"
	"edushedule2/config"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func ScrapeSchedule(faculty, form, course, group, date string) (string, error) {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	var schedule string

	err := chromedp.Run(ctx,
		chromedp.Navigate(config.GetUrl()),
		logStep("Перешли на страницу"),

		chromedp.SetValue(`select[name="ddlFac"]`, faculty, chromedp.ByQuery),
		logStep("Установлено значение для факультета"),

		chromedp.SetValue(`select[name="ddlDep"]`, form, chromedp.ByQuery),
		logStep("Установлено значение для формы"),

		chromedp.SetValue(`select[name="ddlCourse"]`, course, chromedp.ByQuery),
		logStep("Установлено значение для курса"),

		chromedp.SetValue(`select[name="ddlGroup"]`, group, chromedp.ByQuery),
		logStep("Установлено значение для группы"),

		chromedp.SetValue(`select[name="ddlWeek"]`, date, chromedp.ByQuery),
		logStep("Выбрана дата "+date),

		chromedp.Click(`a.chosen-single.button`, chromedp.ByQuery),
		logStep("Нажата кнопка 'Показать'"),

		logStep("Расписание стало видимым"),

		chromedp.WaitVisible("#TT", chromedp.ByID),
		chromedp.Text(`#TT`, &schedule, chromedp.ByID),
		logStep("Получено расписание"),
	)

	if err != nil {
		return "", err
	}
	return schedule, nil

}

func logStep(message string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Println(message)
		return nil
	})
}
