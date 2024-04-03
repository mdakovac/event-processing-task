package description_service

import (
	"fmt"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/currency/currency_service"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type DescriptionService struct {
	currencyService currency_service.CurrencyServiceType
}

type DescriptionServiceType interface {
	AssignDescription(event *casino.Event) *casino.Event
}

func (service *DescriptionService) AssignDescription(event *casino.Event) *casino.Event {
	var formattedPlayerInfo = formatPlayerInfo(event.PlayerID, event.Player.Email)
	var formattedDate = formatDate(event.CreatedAt)

	var description string
	if event.Type == "game_start" || event.Type == "game_stop" {
		description = fmt.Sprintf("%s %s playing a game \"%s\" on %s",
			formattedPlayerInfo,
			getGamePlayingStatus(event.Type),
			casino.Games[event.GameID].Title,
			formattedDate,
		)
	} else {
		var formattedAmount = formatAmount(
			service.currencyService.ConvertAmountToFloat(event.Amount, event.Currency),
			event.Currency,
			service.currencyService.ConvertAmountToFloat(event.AmountEUR, "EUR"),
		)
		if event.Type == "deposit" {
			description = fmt.Sprintf("%s made a deposit of %s on %s", formattedPlayerInfo, formattedAmount, formattedDate)

		} else {
			description = fmt.Sprintf("%s placed a bet of %s on a game \"%s\" on %s",
				formattedPlayerInfo,
				formattedAmount,
				casino.Games[event.GameID].Title,
				formattedDate,
			)
		}
	}

	event.Description = description
	return event
}

func formatPlayerInfo(playerId int, playerEmail string) string {
	var playerInfoString = fmt.Sprintf("Player #%d", playerId)
	if playerEmail != "" {
		playerInfoString += fmt.Sprintf(" (%s)", playerEmail)
	}

	return playerInfoString
}

func getGamePlayingStatus(eventType string) string {
	if eventType == "game_start" {
		return "started"
	}
	return "stopped"
}

func formatDate(date time.Time) string {
	return fmt.Sprintf("%s %s, %d at %d:%d UTC",
		date.Month(),
		getOrdinalNumber(date.Day()),
		date.Year(),
		date.Hour(),
		date.Minute(),
	)
}

func getOrdinalNumber(n int) string {
	if n >= 11 && n <= 13 {
		return fmt.Sprintf("%dth", n)
	}

	switch n % 10 {
	case 1:
		return fmt.Sprintf("%dst", n)
	case 2:
		return fmt.Sprintf("%dnd", n)
	case 3:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
}

func formatAmount(amount float64, currency string, amountEur float64) string {
	var formattedAmount string
	if currency == "BTC" {
		formattedAmount = fmt.Sprintf("%f", amount)
	} else {
		formattedAmount = fmt.Sprintf("%.2f", amount)
	}

	formattedAmountWithCurrency := fmt.Sprintf("%s %s", formattedAmount, currency)
	if currency != "EUR" {
		formattedAmountWithCurrency += fmt.Sprintf(" (%.2f EUR)", amountEur)
	}

	return formattedAmountWithCurrency
}

func NewDescriptionService(currencyService currency_service.CurrencyServiceType) DescriptionServiceType {
	return &DescriptionService{
		currencyService,
	}
}
