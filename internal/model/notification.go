package model

// Notification is what the notifier needs to send one email
type Notification struct {
	ID          int64       `db:"id"`
	OrderID     int64       `db:"order_id"`
	OrderStatus OrderStatus `db:"order_status"`
	UserEmail   string      `db:"user_email"`
	Attempts    int         `db:"attempts"`
}
