package aggregation_models

type TopPlayer struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

type CountByPlayerId map[int]int
type CountByTimestamp map[string]int

type Aggregation struct {
	EventsTotal                  int       `json:"events_total"`
	EventsPerMinute              float64   `json:"events_per_minute"`
	EventsPerSecondMovingAverage float64   `json:"events_per_second_moving_average"`
	TopPlayerBets                TopPlayer `json:"top_player_bets"`
	TopPlayerWins                TopPlayer `json:"top_player_wins"`
	TopPlayerDeposits            TopPlayer `json:"top_player_deposits"`
}
