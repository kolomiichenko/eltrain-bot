// i18n["ua"]["mainMenuText"]

package main

var (
	i18n = map[string]map[string]string{
		"ua": map[string]string{
			"mainMenuText":     "🚆 Для отримання розкладу електричок (приміських поїздів) заповніть всі поля та натисніть кнопку <b>🔎 Знайти</b>",
			"mainMenuLanguage": "🌐 Мова: ",
			"mainMenuFrom":     "🛫 Звідки: ",
			"mainMenuTo":       "🛬 Куди: ",
			"mainMenuReverse":  "↔️ Зворотній напрямок",
			"mainMenuWhen":     "📆 Тип відображення: ",
			"mainMenuSearch":   "🔎 Знайти",

			"selectLanguage": "🌐 Оберіть мову",
			"ua":             "🇺🇦 Українська",
			"en":             "🇺🇸 English",

			"typeFrom":   "🛫 <b>Введіть станцію відправлення</b>",
			"typeTo":     "🛬 <b>Введіть станцію прибуття</b>",
			"selectFrom": "🛫 <b>Виберіть станцію відправлення</b>",
			"selectTo":   "🛬 <b>Виберіть станцію прибуття</b>",

			"selectWhen": "📆 Виберіть тип відображення",
			"whenMenu0":  "Всі дні",
			"whenMenu1":  "Сьогодні",

			"searchError":  "🚫 Заповніть всі поля в /menu",
			"searchTotal":  "ℹ️ На запит знайдено %v результатів(та)\n",
			"searchFooter": "\n\n ⬅️ Назад до /menu",
		},
		"en": map[string]string{
			"mainMenuText":     "🚆 To obtain the timetable of trains (commuter trains) fill all the fields and push the <b>🔎 Find </b> button",
			"mainMenuLanguage": "🌐 Language: ",
			"mainMenuFrom":     "🛫 Departure: ",
			"mainMenuTo":       "🛬 Arrive: ",
			"mainMenuReverse":  "↔️ Reverse derection",
			"mainMenuWhen":     "📆 Display mode: ",
			"mainMenuSearch":   "🔎 Search",

			"typeFrom":   "🛫 <b>Type the departure station</b>",
			"typeTo":     "🛬 <b>Type the arrival station</b>",
			"selectFrom": "🛫 <b>Select the departure station</b>",
			"selectTo":   "🛬 <b>Select the arrival station</b>",

			"selectLanguage": "🌐 Select language",
			"ua":             "🇺🇦 Українська",
			"en":             "🇺🇸 English",

			"selectWhen": "📆 Select the type of display",
			"whenMenu0":  "All days",
			"whenMenu1":  "Today",

			"searchError":  "🚫 Fill all data. Try again /menu",
			"searchTotal":  "ℹ️ %v results found\n",
			"searchFooter": "\n\n ⬅️ Go back to /menu",
		},
	}
)
